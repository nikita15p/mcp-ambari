# Certificate Generation Guide

This guide explains how to generate and manage certificates for mTLS (mutual TLS) authentication in the Ambari MCP Server.

## Overview

The Ambari MCP Server includes a comprehensive certificate generation utility that creates a complete PKI (Public Key Infrastructure) for mTLS authentication:

1. **Generate a Certificate Authority (CA)** - Creates a self-signed CA that can sign other certificates
2. **Generate Server Certificates** - Creates server certificates signed by the CA for the MCP server
3. **Generate Client Certificates** - Creates client certificates signed by the CA for MCP clients

**Important:** All certificates must be generated in advance using this utility. The server loads pre-generated CA certificates at startup - there is no runtime CA generation.

All certificates are properly configured for their specific purpose (server auth or client auth) and are cryptographically signed by the CA to establish a trust chain.

## Quick Start

### Generate Complete Certificate Infrastructure

To generate all certificates (CA, server, and client) in one command:

```bash
./generate-certs --type all --dir certs
```

This will create:
- **CA certificate and key** in `certs/ca/`
- **Server certificate and key** in `certs/server/`
- **Default client certificate and key** in `certs/client/`

### Generate Individual Certificate Types

#### 1. Generate CA Only

```bash
./generate-certs --type ca --dir certs
```

Creates:
- `certs/ca/ca-cert.pem` - CA certificate (public)
- `certs/ca/ca-key.pem` - CA private key (mode 600)

#### 2. Generate Server Certificate

```bash
./generate-certs --type server --dir certs
```

Creates a server certificate signed by the CA:
- `certs/server/server-cert.pem` - Server certificate
- `certs/server/server-key.pem` - Server private key (mode 600)

**Note:** CA must already exist when generating server certificates.

#### 3. Generate Client Certificate

```bash
./generate-certs --type client --client-name mcp-client --dir certs
```

Creates a client certificate signed by the CA:
- `certs/client/mcp-client-cert.pem` - Client certificate
- `certs/client/mcp-client-key.pem` - Client private key (mode 600)

**Note:** CA must already exist when generating client certificates.

## Command-Line Options

### Common Options

| Flag | Default | Description |
|------|---------|-------------|
| `--type` | `all` | Certificate type: `ca`, `server`, `client`, or `all` |
| `--dir` | `certs` | Base directory for certificates |
| `--keysize` | `2048` | RSA key size (2048 or 4096 recommended) |
| `--days` | `365` | Certificate validity period in days |
| `--org` | `Ambari MCP Server` | Organization name |
| `--country` | `US` | Country code (2 letters) |

### Server-Specific Options

| Flag | Default | Description |
|------|---------|-------------|
| `--server-name` | `ambari-mcp-server` | Server common name (CN) |
| `--server-dns` | `localhost,ambari-mcp-server,...` | Comma-separated DNS names for SAN |

### Client-Specific Options

| Flag | Default | Description |
|------|---------|-------------|
| `--client-name` | `mcp-client` | Client common name (CN) |

## Usage Examples

### Example 1: Production Setup with Custom Organization

```bash
# Generate all certificates for production
./generate-certs \
  --type all \
  --dir /etc/ambari-mcp/certs \
  --org "<org>" \
  --country "<country>" \
  --days 730 \
  --keysize 4096
```

### Example 2: Generate Multiple Client Certificates

```bash
# First, generate CA and server (if not already done)
./generate-certs --type ca --dir certs
./generate-certs --type server --dir certs

# Generate client certificates for different users/services
./generate-certs --type client --client-name alice --dir certs
./generate-certs --type client --client-name bob --dir certs
./generate-certs --type client --client-name service-account --dir certs
```

### Example 3: Custom Server DNS Names

```bash
# Generate server certificate with custom DNS names
./generate-certs \
  --type server \
  --server-name mcp.example.com \
  --server-dns "mcp.example.com,*.mcp.example.com,localhost" \
  --dir certs
```

## Certificate Structure

After running `--type all`, your certificate directory will look like:

```
certs/
├── ca/
│   ├── ca-cert.pem       # CA certificate (distribute to clients)
│   └── ca-key.pem        # CA private key (keep secure!)
├── server/
│   ├── server-cert.pem   # Server certificate
│   └── server-key.pem    # Server private key
└── client/
    ├── mcp-client-cert.pem    # Client certificate
    └── mcp-client-key.pem     # Client private key
```

## Server Configuration

After generating certificates, configure your server:

```bash
# Set environment variables
export TLS_CERT_FILE=certs/server/server-cert.pem
export TLS_KEY_FILE=certs/server/server-key.pem
export TLS_CA_FILE=certs/ca/ca-cert.pem
export MCP_TRANSPORT=mtls
export HOST=0.0.0.0
export PORT=9001

# Start server
./server
```

Or use a `.env` file:

```env
MCP_TRANSPORT=mtls
HOST=0.0.0.0
PORT=9001
TLS_CERT_FILE=certs/server/server-cert.pem
TLS_KEY_FILE=certs/server/server-key.pem
TLS_CA_FILE=certs/ca/ca-cert.pem
```

## Client Configuration

Clients need three files to connect:

1. **Client Certificate** (`client/mcp-client-cert.pem`) - Client's identity
2. **Client Private Key** (`client/mcp-client-key.pem`) - Client's private key
3. **CA Certificate** (`ca/ca-cert.pem`) - To verify server's certificate

### Example Client Connection (curl)

```bash
curl -v \
  --cert certs/client/mcp-client-cert.pem \
  --key certs/client/mcp-client-key.pem \
  --cacert certs/ca/ca-cert.pem \
  https://localhost:9001
```

### Example Client Connection (openssl)

```bash
openssl s_client \
  -connect localhost:9001 \
  -cert certs/client/mcp-client-cert.pem \
  -key certs/client/mcp-client-key.pem \
  -CAfile certs/ca/ca-cert.pem
```

## Certificate Verification

### Verify CA Certificate

```bash
# Check CA certificate details
openssl x509 -in certs/ca/ca-cert.pem -text -noout

# Verify CA is self-signed
openssl verify -CAfile certs/ca/ca-cert.pem certs/ca/ca-cert.pem
```

### Verify Server Certificate

```bash
# Check server certificate details
openssl x509 -in certs/server/server-cert.pem -text -noout

# Verify server cert is signed by CA
openssl verify -CAfile certs/ca/ca-cert.pem certs/server/server-cert.pem
```

### Verify Client Certificate

```bash
# Check client certificate details
openssl x509 -in certs/client/mcp-client-cert.pem -text -noout

# Verify client cert is signed by CA
openssl verify -CAfile certs/ca/ca-cert.pem certs/client/mcp-client-cert.pem
```

## CertManager API

The `internal/certs/manager.go` package provides programmatic access to certificate management:

### Generate Client Certificate Programmatically

```go
import (
    "github.com/sirupsen/logrus"
    "ambari-mcp-server/internal/certs"
)

// Create certificate manager
logger := logrus.New()
certMgr := certs.NewCertManager("certs", logger)

// Generate a new client certificate
err := certMgr.GenerateClientCert("new-client", "certs/client", 365, 2048)
if err != nil {
    log.Fatalf("Failed to generate client cert: %v", err)
}
```

### Sign Custom Client Certificate

```go
// Create custom certificate configuration
config := certs.CertConfig{
    CommonName:   "custom-client",
    Organization: "My Organization",
    Country:      "US",
    ValidDays:    365,
    KeySize:      2048,
    IsServer:     false,
}

// Sign with CA
cert, err := certMgr.SignClientCert(config)
if err != nil {
    log.Fatalf("Failed to sign client cert: %v", err)
}

// Save to files
err = certs.SaveCertToFiles(cert, "client-cert.pem", "client-key.pem")
```

## Security Best Practices

### CA Private Key Protection

1. **Restrict Access**: CA private key should have `600` permissions (owner read/write only)
2. **Secure Storage**: Store CA private key in a secure location (HSM, vault, encrypted filesystem)
3. **Limited Distribution**: Never distribute CA private key to clients or servers
4. **Backup**: Keep encrypted backups of CA private key in multiple secure locations

### Certificate Rotation

1. **Monitor Expiration**: Track certificate expiration dates
2. **Renew Early**: Renew certificates 30 days before expiration
3. **Gradual Rollout**: Deploy new certificates gradually in production
4. **Keep Old Certs**: Maintain old certificates during transition period

### Private Key Management

1. **Never Share**: Private keys should never leave their host system
2. **File Permissions**: Always use `600` (owner-only) permissions
3. **No Version Control**: Never commit private keys to git/version control
4. **Secure Transmission**: Use encrypted channels (SSH, TLS) when transferring certificates

## Troubleshooting

### "Failed to load CA" Error

**Problem**: When generating server/client certificates, you get "failed to load CA"

**Solution**: Generate CA first:
```bash
./generate-certs --type ca --dir certs
```

### Certificate Verification Fails

**Problem**: `openssl verify` fails with "unable to get local issuer certificate"

**Solution**: Ensure you're using the correct CA file:
```bash
openssl verify -CAfile certs/ca/ca-cert.pem certs/server/server-cert.pem
```

### mTLS Handshake Failures

**Problem**: Server rejects client connection with "bad certificate"

**Solutions**:
1. Verify client cert is signed by the same CA as server trusts
2. Check certificate hasn't expired
3. Ensure client is presenting both cert and key
4. Verify CA certificate is loaded correctly on server

### Permission Denied on Private Key

**Problem**: Cannot read private key file

**Solution**: Check file permissions:
```bash
# Fix permissions (owner read/write only)
chmod 600 certs/ca/ca-key.pem
chmod 600 certs/server/server-key.pem
chmod 600 certs/client/*-key.pem
```

## Building the Tool

To rebuild the certificate generation tool:

```bash
go build -o generate-certs ./cmd/generate-certs
```

## Related Documentation

- [MTLS_SETUP.md](./MTLS_SETUP.md) - Complete mTLS setup guide
- [MCP_SETUP.md](./MCP_SETUP.md) - MCP server configuration
- [README.md](./README.md) - General project documentation

## Summary

The certificate generation tool provides a complete solution for creating a PKI (Public Key Infrastructure) for mTLS authentication. By generating a CA and using it to sign both server and client certificates, you establish a trust chain that enables secure mutual authentication between the MCP server and its clients.
