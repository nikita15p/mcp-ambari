/* START GENAI */
// Package certs provides certificate management utilities for mTLS authentication
package certs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// CertManager handles certificate operations for mTLS
type CertManager struct {
	logger     *logrus.Logger
	certsDir   string
	caCertPath string
	caKeyPath  string
}

// NewCertManager creates a new certificate manager
func NewCertManager(certsDir string, logger *logrus.Logger) *CertManager {
	return &CertManager{
		logger:     logger,
		certsDir:   certsDir,
		caCertPath: filepath.Join(certsDir, "ca", "ca-cert.pem"),
		caKeyPath:  filepath.Join(certsDir, "ca", "ca-key.pem"),
	}
}

// SignClientCert signs a client certificate with the CA
func (cm *CertManager) SignClientCert(config CertConfig) (*CertResult, error) {
	cm.logger.WithField("common_name", config.CommonName).Info("Signing client certificate with CA")
	
	// Load CA certificate and key
	ca, err := LoadCA(cm.caCertPath, cm.caKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA: %v", err)
	}
	
	// Ensure this is a client certificate
	config.IsServer = false
	
	// Generate and sign the certificate
	cert, err := GenerateCertificate(config, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client certificate: %v", err)
	}
	
	cm.logger.WithFields(logrus.Fields{
		"common_name": config.CommonName,
		"valid_days":  config.ValidDays,
	}).Info("Successfully signed client certificate")
	
	return cert, nil
}

// GenerateClientCert creates and saves a CA-signed client certificate
func (cm *CertManager) GenerateClientCert(commonName, outputDir string, validDays, keySize int) error {
	cm.logger.WithField("common_name", commonName).Info("Generating CA-signed client certificate")
	
	// Load CA certificate and key
	ca, err := LoadCA(cm.caCertPath, cm.caKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load CA: %v", err)
	}
	
	// Configure client certificate
	config := CertConfig{
		CommonName:   commonName,
		Organization: "Ambari MCP Client",
		Country:      "US",
		ValidDays:    validDays,
		KeySize:      keySize,
		IsServer:     false,
	}
	
	// Generate certificate
	cert, err := GenerateCertificate(config, ca)
	if err != nil {
		return fmt.Errorf("failed to generate client certificate: %v", err)
	}
	
	// Save certificate to files
	certPath := filepath.Join(outputDir, fmt.Sprintf("%s-cert.pem", commonName))
	keyPath := filepath.Join(outputDir, fmt.Sprintf("%s-key.pem", commonName))
	
	if err := SaveCertToFiles(cert, certPath, keyPath); err != nil {
		return fmt.Errorf("failed to save client certificate: %v", err)
	}
	
	cm.logger.WithFields(logrus.Fields{
		"common_name": commonName,
		"cert_path":   certPath,
		"key_path":    keyPath,
	}).Info("Successfully generated and saved client certificate")
	
	return nil
}

// GetCAPaths returns the paths to CA certificate and key
func (cm *CertManager) GetCAPaths() (certPath, keyPath string) {
	return cm.caCertPath, cm.caKeyPath
}

// CAExists checks if CA certificate and key exist
func (cm *CertManager) CAExists() bool {
	_, certErr := os.Stat(cm.caCertPath)
	_, keyErr := os.Stat(cm.caKeyPath)
	return certErr == nil && keyErr == nil
}
/* END GENAI */
