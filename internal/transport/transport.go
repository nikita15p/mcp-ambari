// Package transport provides pluggable transport implementations using the
// Strategy pattern. Each transport mode (stdio, HTTP, SSL, mTLS) implements
// the Transport interface. A Factory creates the correct transport based on
// configuration.
package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/niita15p/mcp-ambari/internal/auth"
	"github.com/sirupsen/logrus"
)

// Mode defines the transport mode
type Mode string

const (
	ModeStdio Mode = "stdio"
	ModeHTTP  Mode = "http"
	ModeSSL   Mode = "ssl"
	ModeMTLS  Mode = "mtls"
)

// Config holds transport configuration
type Config struct {
	Mode       Mode   `json:"mode"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	SSLCert    string `json:"ssl_cert"`
	SSLKey     string `json:"ssl_key"`
	SSLCACerts string `json:"ssl_ca_certs"`
}

// MCPServer wraps the actual mcp.Server for transport use
type MCPServer struct {
	*mcp.Server
}

// Listen implements a simplified interface for stdio transport
func (m *MCPServer) Listen(ctx context.Context) error {
	// For stdio, use the SDK's stdio transport
	return m.Server.Run(ctx, &mcp.StdioTransport{})
}

// Transport is the Strategy interface for MCP server transport
type Transport interface {
	// Name returns the transport mode name
	Name() Mode

	// Start launches the transport and blocks until shutdown
	Start(ctx context.Context, mcpServer *MCPServer) error

	// Description returns human-readable transport info for logging
	Description() string
}

// Factory creates the appropriate transport based on configuration (Factory pattern)
func Factory(cfg Config, authMW *auth.Middleware, logger *logrus.Logger) (Transport, error) {
	switch cfg.Mode {
	case ModeStdio:
		return &StdioTransport{logger: logger}, nil
	case ModeHTTP:
		return &HTTPTransport{cfg: cfg, authMW: authMW, logger: logger}, nil
	case ModeSSL:
		return &SSLTransport{cfg: cfg, authMW: authMW, logger: logger}, nil
	case ModeMTLS:
		return &MTLSTransport{cfg: cfg, authMW: authMW, logger: logger}, nil
	default:
		return nil, fmt.Errorf("unsupported transport mode: %s (supported: stdio, http, ssl, mtls)", cfg.Mode)
	}
}

// ---------- HTTP Transport ----------

// HTTPTransport implements Transport for HTTP/streamableHttp mode
type HTTPTransport struct {
	cfg    Config
	authMW *auth.Middleware
	logger *logrus.Logger
}

func (t *HTTPTransport) Name() Mode { return ModeHTTP }
func (t *HTTPTransport) Description() string {
	return fmt.Sprintf("HTTP transport on http://%s:%s — for streamableHttp clients", t.cfg.Host, t.cfg.Port)
}

func (t *HTTPTransport) Start(ctx context.Context, mcpServer *MCPServer) error {
	addr := fmt.Sprintf("%s:%s", t.cfg.Host, t.cfg.Port)

	t.logger.WithFields(logrus.Fields{
		"host": t.cfg.Host,
		"port": t.cfg.Port,
		"mode": "http",
		"addr": addr,
	}).Info("Starting HTTP transport for MCP")
	
	// Create streamable HTTP handler for MCP-over-HTTP (using MCP Go SDK)
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return mcpServer.Server
	}, nil)
	
	// Apply auth middleware if provided
	var httpHandler http.Handler = handler
	if t.authMW != nil {
		httpHandler = t.authMW.Handler(handler)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: httpHandler,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		t.logger.Info("Shutting down HTTP server")
		server.Close()
	}()

	t.logger.WithField("addr", addr).Info("HTTP transport server starting")
	return server.ListenAndServe()
}

// ---------- SSL Transport ----------

// SSLTransport implements Transport for HTTPS/TLS mode
type SSLTransport struct {
	cfg    Config
	authMW *auth.Middleware
	logger *logrus.Logger
}

func (t *SSLTransport) Name() Mode { return ModeSSL }
func (t *SSLTransport) Description() string {
	return fmt.Sprintf("HTTPS transport on https://%s:%s — TLS encrypted streamableHttp", t.cfg.Host, t.cfg.Port)
}

func (t *SSLTransport) Start(ctx context.Context, mcpServer *MCPServer) error {
	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair(t.cfg.SSLCert, t.cfg.SSLKey)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Create streamable HTTP handler for MCP-over-HTTPS
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return mcpServer.Server
	}, nil)

	// Apply auth middleware if provided
	var httpHandler http.Handler = handler
	if t.authMW != nil {
		httpHandler = t.authMW.Handler(handler)
	}

	addr := fmt.Sprintf("%s:%s", t.cfg.Host, t.cfg.Port)
	server := &http.Server{
		Addr:      addr,
		Handler:   httpHandler,
		TLSConfig: tlsConfig,
	}

	t.logger.WithFields(logrus.Fields{
		"host": t.cfg.Host,
		"port": t.cfg.Port,
		"mode": "https",
		"cert": t.cfg.SSLCert,
		"addr": addr,
	}).Info("Starting HTTPS transport for MCP")

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		t.logger.Info("Shutting down HTTPS server")
		server.Close()
	}()

	return server.ListenAndServeTLS("", "") // Certificates are in TLS config
}

// ---------- mTLS Transport ----------

// MTLSTransport implements Transport for HTTPS with mutual TLS mode
type MTLSTransport struct {
	cfg    Config
	authMW *auth.Middleware
	logger *logrus.Logger
}

func (t *MTLSTransport) Name() Mode { return ModeMTLS }
func (t *MTLSTransport) Description() string {
	return fmt.Sprintf("HTTPS mTLS transport on https://%s:%s — mutual TLS authentication", t.cfg.Host, t.cfg.Port)
}

func (t *MTLSTransport) Start(ctx context.Context, mcpServer *MCPServer) error {
	// Load server TLS certificate
	cert, err := tls.LoadX509KeyPair(t.cfg.SSLCert, t.cfg.SSLKey)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %v", err)
	}

	// Load CA certificate for client verification
	caCert, err := ioutil.ReadFile(t.cfg.SSLCACerts)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}

	// Create streamable HTTP handler for MCP-over-HTTPS-mTLS
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return mcpServer.Server
	}, nil)

	// Apply auth middleware if provided
	var httpHandler http.Handler = handler
	if t.authMW != nil {
		httpHandler = t.authMW.Handler(handler)
	}

	addr := fmt.Sprintf("%s:%s", t.cfg.Host, t.cfg.Port)
	server := &http.Server{
		Addr:      addr,
		Handler:   httpHandler,
		TLSConfig: tlsConfig,
	}

	t.logger.WithFields(logrus.Fields{
		"host": t.cfg.Host,
		"port": t.cfg.Port,
		"mode": "https-mtls",
		"cert": t.cfg.SSLCert,
		"ca":   t.cfg.SSLCACerts,
		"addr": addr,
	}).Info("Starting HTTPS mTLS transport for MCP")

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		t.logger.Info("Shutting down HTTPS mTLS server")
		server.Close()
	}()

	return server.ListenAndServeTLS("", "") // Certificates are in TLS config
}

// ---------- Stdio Transport ----------

// StdioTransport implements Transport for stdio (default MCP transport)
type StdioTransport struct {
	logger *logrus.Logger
}

func (t *StdioTransport) Name() Mode { return ModeStdio }
func (t *StdioTransport) Description() string {
	return "Stdio transport (stdin/stdout) — for local MCP clients like Claude Desktop"
}

func (t *StdioTransport) Start(ctx context.Context, mcpServer *MCPServer) error {
	t.logger.Info("Starting MCP server on stdio transport")
	return mcpServer.Listen(ctx)
}
