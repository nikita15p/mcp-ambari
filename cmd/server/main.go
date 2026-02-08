// Package main is the entry point for the Ambari MCP Server.
// It uses the official MCP Go SDK for protocol compliance and clean architecture.
//
// Architecture highlights:
//   - Read-only vs Actionable operations are cleanly separated
//   - Template Method pattern for operation execution lifecycle
//   - Strategy pattern for pluggable authentication (LDAP/mTLS)
//   - Registry/Factory pattern for operation management
//   - Official MCP Go SDK for protocol compliance
package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"

	"github.com/nikita15p/mcp-ambari/internal/auth"
	"github.com/nikita15p/mcp-ambari/internal/client"
	ops "github.com/nikita15p/mcp-ambari/internal/operations"
	"github.com/nikita15p/mcp-ambari/internal/operations/actionable"
	"github.com/nikita15p/mcp-ambari/internal/operations/readonly"
	"github.com/nikita15p/mcp-ambari/internal/resources"
	"github.com/nikita15p/mcp-ambari/internal/transport"
)

func main() {
	// --- Logger ---
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	level, err := logrus.ParseLevel(envOr("LOG_LEVEL", "info"))
	if err == nil {
		logger.SetLevel(level)
	}

	logger.Info("Starting Tusker Ambari MCP Server (Go)")

	// --- Ambari Client ---
	timeout, _ := time.ParseDuration(envOr("AMBARI_TIMEOUT", "30s"))
	ambariClient := client.NewAmbariClient(client.Config{
		BaseURL:  envOr("AMBARI_BASE_URL", "http://localhost:8080/api/v1"),
		Username: envOr("AMBARI_USERNAME", "admin"),
		Password: envOr("AMBARI_PASSWORD", "admin"),
		Timeout:  timeout,
		Retries:  3,
	}, logger)

	// --- Operation Registry (Registry/Factory pattern) ---
	registry := ops.NewRegistry(logger)

	// Register READ-ONLY operations (safe, GET-only, lower permissions) — 24 tools
	readOnlyOps := []ops.Operation{
		// Clusters
		readonly.NewGetClusters(ambariClient, logger),
		readonly.NewGetCluster(ambariClient, logger),
		// Services
		readonly.NewGetServices(ambariClient, logger),
		readonly.NewGetService(ambariClient, logger),
		readonly.NewGetServiceState(ambariClient, logger),
		readonly.NewGetServicesWithStaleConfigs(ambariClient, logger),
		readonly.NewGetHostComponentsWithStaleConfigs(ambariClient, logger),
		readonly.NewGetRollingRestartStatus(ambariClient, logger),
		readonly.NewIsServiceCheckSupported(ambariClient, logger),
		readonly.NewGetServiceCheckStatus(ambariClient, logger),
		// Hosts
		readonly.NewGetHosts(ambariClient, logger),
		readonly.NewGetHost(ambariClient, logger),
		// Alerts
		readonly.NewGetAlerts(ambariClient, logger),
		readonly.NewGetAlertSummary(ambariClient, logger),
		readonly.NewGetAlertDetails(ambariClient, logger),
		readonly.NewGetAlertDefinitions(ambariClient, logger),
		readonly.NewGetAlertGroups(ambariClient, logger),
		readonly.NewGetAlertTargets(ambariClient, logger),
		readonly.NewGetNotifications(ambariClient, logger),
		// Users and Groups
		readonly.NewGetUsers(ambariClient, logger),
		readonly.NewGetUser(ambariClient, logger),
		readonly.NewGetGroups(ambariClient, logger),
		readonly.NewGetGroup(ambariClient, logger),
		readonly.NewGetUserPrivileges(ambariClient, logger),
	}
	for _, op := range readOnlyOps {
		if err := registry.Register(op); err != nil {
			logger.WithError(err).Fatal("Failed to register read-only operation")
		}
	}

	// Register ACTIONABLE operations (state-changing, higher permissions) — 27 tools
	// Can be disabled via ENABLE_ACTIONABLE_TOOLS=false environment variable
	enableActionable := strings.ToLower(envOr("ENABLE_ACTIONABLE_TOOLS", "true")) == "true"

	if enableActionable {
		actionableOps := []ops.Operation{
			// Cluster management
			actionable.NewCreateCluster(ambariClient, logger),
			// Service lifecycle
			actionable.NewStartService(ambariClient, logger),
			actionable.NewStopService(ambariClient, logger),
			actionable.NewRestartService(ambariClient, logger),
			actionable.NewRestartComponents(ambariClient, logger),
			actionable.NewEnableMaintenanceMode(ambariClient, logger),
			actionable.NewDisableMaintenanceMode(ambariClient, logger),
			actionable.NewRunServiceCheck(ambariClient, logger),
			// Alert definitions
			actionable.NewUpdateAlertDefinition(ambariClient, logger),
			// Alert groups
			actionable.NewCreateAlertGroup(ambariClient, logger),
			actionable.NewUpdateAlertGroup(ambariClient, logger),
			actionable.NewDeleteAlertGroup(ambariClient, logger),
			actionable.NewDuplicateAlertGroup(ambariClient, logger),
			actionable.NewAddDefinitionToGroup(ambariClient, logger),
			actionable.NewRemoveDefinitionFromGroup(ambariClient, logger),
			// Alert notifications
			actionable.NewCreateNotification(ambariClient, logger),
			actionable.NewUpdateNotification(ambariClient, logger),
			actionable.NewDeleteNotification(ambariClient, logger),
			actionable.NewAddNotificationToGroup(ambariClient, logger),
			actionable.NewRemoveNotificationFromGroup(ambariClient, logger),
			// Alert settings
			actionable.NewSaveAlertSettings(ambariClient, logger),
			// User and Group management
			actionable.NewCreateUser(ambariClient, logger),
			actionable.NewUpdateUser(ambariClient, logger),
			actionable.NewDeleteUser(ambariClient, logger),
			actionable.NewCreateUserGroup(ambariClient, logger),
			actionable.NewDeleteUserGroup(ambariClient, logger),
			actionable.NewAddUserToGroup(ambariClient, logger),
			actionable.NewRemoveUserFromGroup(ambariClient, logger),
		}
		for _, op := range actionableOps {
			if err := registry.Register(op); err != nil {
				logger.WithError(err).Fatal("Failed to register actionable operation")
			}
		}
	} else {
		logger.Info("Actionable tools disabled via ENABLE_ACTIONABLE_TOOLS=false")
	}

	total, ro, act := registry.Count()
	logger.WithFields(logrus.Fields{
		"total": total, "readonly": ro, "actionable": act,
	}).Info("Operations registered")

	// --- Operation Executor (Template Method pattern) ---
	executor := ops.NewExecutor(ambariClient, logger)

	// --- MCP Server using Go SDK ---
	implementation := &mcp.Implementation{
		Name:    "mcp-ambari",
		Version: "1.0.0",
	}
	mcpServer := mcp.NewServer(implementation, nil)

	// Register each operation as an MCP tool via the SDK
	for _, op := range registry.All() {
		registerMCPTool(mcpServer, op, executor, logger)
	}

	// --- MCP Resources (all read-only, accessed by URI) ---
	resRegistry := resources.NewRegistry(ambariClient, logger)
	for _, resDef := range resRegistry.Definitions() {
		registerMCPResource(mcpServer, resDef, resRegistry, logger)
	}

	logger.WithFields(logrus.Fields{
		"tools": total, "resources": resRegistry.Count(),
	}).Info("MCP server fully initialized")

	// --- Parse CLI flags ---
	var flagTransport string
	var flagHost string
	var flagPort string
	flag.StringVar(&flagTransport, "transport", envOr("MCP_TRANSPORT", "stdio"), "Transport mode (stdio, http, https, https-mtls)")
	flag.StringVar(&flagHost, "host", envOr("MCP_HOST", "127.0.0.1"), "Host address for HTTP transport")
	flag.StringVar(&flagPort, "port", envOr("MCP_PORT", "8090"), "Port for HTTP transport")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info("Shutdown signal received")
		cancel()
	}()

	// --- Authentication Middleware ---
	// Create default LDAP provider for development (disabled)
	ldapProvider := auth.NewLDAPProvider("x-remote-", defaultGroupMappings(), []string{"cluster:admin", "service:admin"}, logger)
	authMW := auth.NewMiddleware(ldapProvider, false, logger) // Disabled auth for development

	// --- Transport Configuration ---
	transportCfg := transport.Config{
		Mode:       transport.Mode(strings.ToLower(flagTransport)),
		Host:       flagHost,
		Port:       flagPort,
		SSLCert:    envOr("TLS_CERT_FILE", ""),
		SSLKey:     envOr("TLS_KEY_FILE", ""),
		SSLCACerts: envOr("TLS_CA_FILE", ""),
	}

	// Create transport using factory
	transportImpl, err := transport.Factory(transportCfg, authMW, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create transport")
	}

	logger.WithFields(logrus.Fields{
		"mode":        transportCfg.Mode,
		"description": transportImpl.Description(),
	}).Info("Transport configured")

	// Start the transport
	mcpServerWrapper := &transport.MCPServer{Server: mcpServer}
	if err := transportImpl.Start(ctx, mcpServerWrapper); err != nil {
		logger.WithError(err).Fatal("Transport failed")
	}
}

// registerMCPTool bridges our Operation interface to the SDK's mcp.Server using the proper API
func registerMCPTool(server *mcp.Server, op ops.Operation, executor *ops.Executor, logger *logrus.Logger) {
	def := op.Definition()

	// Create MCP tool definition
	tool := &mcp.Tool{
		Name:        def.Name,
		Description: def.Description,
	}

	// Create the tool handler function that matches the SDK's expected signature
	handler := func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, map[string]interface{}, error) {
		// Create default auth context for stdio
		authCtx := &auth.AuthContext{
			Username: "stdio-user", Groups: []string{"ambari-admins"},
			Permissions: auth.PermissionGroups["ADMIN"],
			IsValidated: true, Source: "stdio",
		}

		// Execute the operation through our executor
		result, err := executor.Run(ctx, op, input, authCtx)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"tool": op.Name(), "type": op.Type(), "error": err,
			}).Error("Operation failed")
			return nil, nil, err
		}

		return nil, map[string]interface{}{"result": result}, nil
	}

	// Register the tool with the SDK
	mcp.AddTool(server, tool, handler)

	logger.WithFields(logrus.Fields{
		"tool": def.Name, "type": op.Type(), "category": op.Category(),
	}).Debug("MCP tool registered")
}

// registerMCPResource bridges our resource registry to the SDK's mcp.Server using the proper API
func registerMCPResource(server *mcp.Server, resDef resources.ResourceDefinition, resReg *resources.Registry, logger *logrus.Logger) {
	// Create MCP resource definition
	resource := &mcp.Resource{
		URI:         resDef.URI,
		Name:        resDef.Name,
		Description: resDef.Description,
		MIMEType:    resDef.MimeType,
	}

	// Create resource handler
	handler := mcp.ResourceHandler(func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		// Use our resource registry to resolve the resource
		result, err := resReg.Read(ctx, req.Params.URI)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"uri":   req.Params.URI,
				"error": err,
			}).Error("Resource read failed")
			return nil, err
		}

		// Convert data to JSON string for MCP
		var content string
		if result.Data != nil {
			if str, ok := result.Data.(string); ok {
				content = str
			} else {
				// Convert to JSON
				jsonBytes, jsonErr := json.Marshal(result.Data)
				if jsonErr != nil {
					return nil, jsonErr
				}
				content = string(jsonBytes)
			}
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      req.Params.URI,
					MIMEType: resDef.MimeType,
					Text:     content,
				},
			},
		}, nil
	})

	// Register the resource with the SDK
	server.AddResource(resource, handler)

	logger.WithFields(logrus.Fields{
		"uri":  resDef.URI,
		"name": resDef.Name,
	}).Debug("MCP resource registered")
}

// --- Helpers ---

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func defaultGroupMappings() map[string][]string {
	return map[string][]string{
		"ambari-admins":    {"cluster:admin", "service:admin", "alert:admin", "config:modify", "host:manage"},
		"hadoop-operators": {"cluster:operate", "service:operate", "service:restart", "alert:manage"},
		"data-engineers":   {"cluster:view", "service:view", "service:operate", "config:view"},
		"bigdata-viewers":  {"cluster:view", "service:view", "alert:view", "config:view", "host:view"},
	}
}
