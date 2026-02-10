/* START GENAI */
// Package certs provides certificate generation utilities for mTLS authentication
package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// CAConfig holds configuration for CA certificate generation
type CAConfig struct {
	Organization string
	Country      string
	ValidDays    int
	KeySize      int
}

// CertConfig holds configuration for server/client certificate generation
type CertConfig struct {
	CommonName   string
	Organization string
	Country      string
	ValidDays    int
	KeySize      int
	DNSNames     []string
	IPAddresses  []net.IP
	IsServer     bool // true for server cert, false for client cert
}

// CAResult holds the generated CA certificate and private key
type CAResult struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
	CertPEM     []byte
	KeyPEM      []byte
}

// CertResult holds the generated certificate and private key
type CertResult struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
	CertPEM     []byte
	KeyPEM      []byte
}

// GenerateCA creates a new Certificate Authority
func GenerateCA(config CAConfig) (*CAResult, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, config.KeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CA private key: %v", err)
	}

	// Create CA certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{config.Organization},
			Country:      []string{config.Country},
			CommonName:   "Ambari MCP Server CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(config.ValidDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	// Self-sign the CA certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA certificate: %v", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return &CAResult{
		Certificate: cert,
		PrivateKey:  privateKey,
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}

// GenerateCertificate creates a server or client certificate signed by the CA
func GenerateCertificate(config CertConfig, ca *CAResult) (*CertResult, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, config.KeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{config.Organization},
			Country:      []string{config.Country},
			CommonName:   config.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(config.ValidDays) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	// Set extended key usage based on certificate type
	if config.IsServer {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		template.DNSNames = config.DNSNames
		template.IPAddresses = config.IPAddresses
	} else {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	// Sign the certificate with CA
	certDER, err := x509.CreateCertificate(rand.Reader, &template, ca.Certificate, &privateKey.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return &CertResult{
		Certificate: cert,
		PrivateKey:  privateKey,
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}

// SaveCAToFiles writes CA certificate and key to files
func SaveCAToFiles(ca *CAResult, certPath, keyPath string) error {
	// Create directory if needed
	dir := filepath.Dir(certPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write certificate
	if err := os.WriteFile(certPath, ca.CertPEM, 0644); err != nil {
		return fmt.Errorf("failed to write CA certificate: %v", err)
	}

	// Write private key with restricted permissions
	if err := os.WriteFile(keyPath, ca.KeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write CA private key: %v", err)
	}

	return nil
}

// SaveCertToFiles writes certificate and key to files
func SaveCertToFiles(cert *CertResult, certPath, keyPath string) error {
	// Create directory if needed
	dir := filepath.Dir(certPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write certificate
	if err := os.WriteFile(certPath, cert.CertPEM, 0644); err != nil {
		return fmt.Errorf("failed to write certificate: %v", err)
	}

	// Write private key with restricted permissions
	if err := os.WriteFile(keyPath, cert.KeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	return nil
}

// LoadCA loads a CA certificate and private key from files
func LoadCA(certPath, keyPath string) (*CAResult, error) {
	// Read certificate
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode CA certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	// Read private key
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA private key: %v", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil || keyBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode CA private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA private key: %v", err)
	}

	return &CAResult{
		Certificate: cert,
		PrivateKey:  privateKey,
		CertPEM:     certPEM,
		KeyPEM:      keyPEM,
	}, nil
}
/* END GENAI */
