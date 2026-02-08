# Ambari MCP Server Setup Guide

## 1. Configure Cluster Connection

The Ambari MCP Server connects to your Ambari cluster via environment variables. Configure these before running the server:

### Required Environment Variables

```bash
# Copy the example configuration
cp .env.example .env

# Edit the configuration
export AMBARI_BASE_URL=http://your-ambari-server:8080/api/v1
export AMBARI_USERNAME=admin
export AMBARI_PASSWORD=your-password
```

### Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `AMBARI_BASE_URL` | Ambari REST API endpoint | `http://localhost:8080/api/v1` |
| `AMBARI_USERNAME` | Ambari admin username | `admin` |
| `AMBARI_PASSWORD` | Ambari admin password | `admin` |
| `AMBARI_TIMEOUT` | Request timeout | `30s` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` |

## 2. Add Server to MCP Clients

### For Claude Desktop

Add to your Claude Desktop config file (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "ambari-server": {
      "command": "/path/to/mcp-ambari/server",
      "args": ["-transport", "stdio"],
      "env": {
        "AMBARI_BASE_URL": "http://your-ambari-server:8080/api/v1",
        "AMBARI_USERNAME": "admin", 
        "AMBARI_PASSWORD": "your-password"
      }
    }
  }
}
```

### For Cline (VS Code Extension)

Add to Cline's MCP settings (`~/.cline/mcp_servers.json`):

```json
{
  "mcpServers": {
    "ambari-server": {
      "command": "/path/to/mcp-ambari/server",
      "args": ["-transport", "stdio"],
      "env": {
        "AMBARI_BASE_URL": "http://your-ambari-server:8080/api/v1",
        "AMBARI_USERNAME": "admin",
        "AMBARI_PASSWORD": "your-password"
      }
    }
  }
}
```

## 3. Available Tools & Resources

Once configured, the MCP server provides:

### Tools (40 total)
- **Read-only (19)**: Get clusters, services, hosts, alerts, configurations
- **Actionable (21)**: Start/stop services, restart components, manage alerts

### Resources (12 total)
- `ambari://clusters` - List all clusters
- `ambari://cluster/{clusterName}` - Cluster details
- `ambari://cluster/{clusterName}/services` - Cluster services
- `ambari://cluster/{clusterName}/hosts` - Cluster hosts
- `ambari://cluster/{clusterName}/alerts` - Cluster alerts
- And more...

## 4. Testing the Connection

### Manual Test
```bash
# Build the server
go build ./cmd/server

# Test with your configuration
AMBARI_BASE_URL=http://your-ambari:8080/api/v1 \
AMBARI_USERNAME=admin \
AMBARI_PASSWORD=your-password \
./server -transport stdio
```

### Verify in MCP Client
1. Restart your MCP client (Claude Desktop/Cline)
2. Look for "ambari-server" in the available tools
3. Try running: "List all clusters" or "Get cluster status"

## 5. Security Configuration

### Authentication (Optional)
For production use with authentication headers:

```bash
export AUTH_ENABLED=true
export LDAP_HEADER_PREFIX=x-user-
```

### HTTPS/mTLS (Optional)
For secure transport modes:

```bash
export MCP_TRANSPORT=mtls
export SSL_CERTFILE=certs/server-cert.pem
export SSL_KEYFILE=certs/server-key.pem
export SSL_CA_CERTS=certs/ca.pem
```

## 6. Troubleshooting

### Common Issues

1. **"Connection refused"**: Check `AMBARI_BASE_URL` is correct and accessible
2. **"Authentication failed"**: Verify `AMBARI_USERNAME` and `AMBARI_PASSWORD`
3. **"Server not found in MCP client"**: Check the client config file path and JSON syntax
4. **"Permission denied"**: Ensure the server binary has execute permissions (`chmod +x server`)

### Debug Logging
```bash
export LOG_LEVEL=debug
./server -transport stdio
```

### Test Ambari API Manually
```bash
curl -u admin:password http://your-ambari:8080/api/v1/clusters
```

## 7. Available Operations

The server exposes 40 Ambari operations as MCP tools:

**Cluster Operations:**
- `ambari_clusters_getclusters` - List all clusters
- `ambari_clusters_getcluster` - Get cluster details  
- `ambari_clusters_createcluster` - Create new cluster

**Service Operations:**
- `ambari_services_getservices` - List cluster services
- `ambari_services_startservice` - Start a service
- `ambari_services_stopservice` - Stop a service
- `ambari_services_restartservice` - Restart a service

**Alert Operations:**
- `ambari_alerts_getalerts` - Get cluster alerts
- `ambari_alerts_getalertsummary` - Get alert summary
- `ambari_alerts_createalertgroup` - Create alert group

And many more! Each tool includes parameter validation and permission checking.