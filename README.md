# Ambari MCP Server

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A high-performance Model Context Protocol (MCP) server for Apache Ambari implemented in Go. This project enables AI agents and assistants to seamlessly interact with Ambari clusters through standardized MCP tools and resources.

## Overview

The Ambari MCP Server acts as a bridge between AI agents and Apache Ambari clusters, providing comprehensive access for automated cluster management, service operations, monitoring, and troubleshooting through the Model Context Protocol. **The server is designed to be consumed by AI agents** rather than human users directly.

**Agent-Centric Architecture:**
```
AI Agents → Ambari MCP Server → Apache Ambari REST API → Hadoop Cluster
```

**Key Benefits:**
- 🚀 **High Performance**: Built in Go with connection pooling and retry logic  
- 🔒 **Enterprise Security**: LDAP authentication with role-based permissions [WIP]
- 📊 **Comprehensive Coverage**: 52+ tools covering all major Ambari operations  including user/group management
- 🔧 **Production Ready**: Robust error handling , TLS/mTLS support, and graceful shutdown

## Architecture

The server implements several design patterns for maintainability and extensibility:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AI Agents     │    │  Transport Layer │    │  Auth Provider  │
│                 │    │                  │    │                 │  
│ • Claude         │◄──►│ • Stdio (MCP)    │◄──►│ • LDAP Headers  │
│ • Cline         │    │ • HTTP/HTTPS     │    │ • Permission    │
│ • Custom Agents │    │ • mTLS           │    │   Groups        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
            ┌─────────────────────────────────────────────┐
            │         Operation Registry                  │
            │                                             │
            │  ┌─────────────────┐ ┌─────────────────────┐│
            │  │   Read-Only     │ │    Actionable       ││
            │  │   Operations    │ │    Operations       ││
            │  │   (24 tools)    │ │    (28 tools)       ││
            │  │                 │ │                     ││
            │  │ • Get clusters  │ │ • Start services    ││
            │  │ • List services │ │ • Restart components││
            │  │ • View alerts   │ │ • Create clusters   ││
            │  │ • Check status  │ │ • Manage alerts     ││
            │  └─────────────────┘ └─────────────────────┘│
            └─────────────────────────────────────────────┘
                                │
                                ▼
            ┌─────────────────────────────────────────────┐
            │         Template Method Executor            │
            │                                             │
            │  Authorization → Validation → Execution     │
            └─────────────────────────────────────────────┘
                                │
                                ▼
            ┌─────────────────────────────────────────────┐
            │          Ambari REST Client                 │
            │                                             │
            │  • Connection pooling  • Retry logic        │
            │  • Request timeout     • Error handling     │  
            └─────────────────────────────────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │    Ambari Cluster     │
                    │                       │
                    │  <server name>        │
                    │                       │
                    │  • Services          │
                    │  • Hosts             │
                    │  • Configurations    │
                    │  • Alerts            │
                    └───────────────────────┘
```

### Design Patterns Implemented

- **Strategy Pattern**: Pluggable authentication providers and transport modes
- **Template Method**: Standardized operation execution lifecycle  
- **Factory Pattern**: Dynamic operation and transport creation
- **Registry Pattern**: Centralized operation management
- **Repository Pattern**: Ambari client with connection abstraction

### SOLID Principles

- **Single Responsibility**: Each package has one clear purpose
- **Open/Closed**: Extensible through interfaces without modification  
- **Liskov Substitution**: All operations implement Operation interface
- **Interface Segregation**: Minimal, focused interfaces
- **Dependency Inversion**: Dependencies injected via constructors

## Features

### 🛠️ **51 MCP Tools Available**


- **51 Tools**: Operations for managing and monitoring Ambari clusters
- **12 Resources**: URI-based access to cluster data
- **8 Prompts**: Guided workflows for common tasks

---

## 🛠️ Tools (51 Total)

### Read-Only Tools (24) - Always Available

These tools are safe, GET-only operations with lower permission requirements:

#### Cluster Operations (2)
| Tool Name | Description |
|-----------|-------------|
| `ambari_clusters_getclusters` | List all Ambari clusters |
| `ambari_clusters_getcluster` | Get detailed information about a specific cluster |

#### Service Operations (8)
| Tool Name | Description |
|-----------|-------------|
| `ambari_services_getservices` | List all services in a cluster |
| `ambari_services_getservice` | Get detailed service information |
| `ambari_services_getservicestate` | Get comprehensive service state details |
| `ambari_services_getserviceswithstaleconfigs` | Find services needing restart |
| `ambari_services_gethostcomponentswithstaleconfigs` | Find components needing restart |
| `ambari_services_getrollingrestartstatus` | Monitor rolling restart operations |
| `ambari_services_isservicechecksupported` | Check if service supports health checks |
| `ambari_services_getservicecheckstatus` | Get status of service health checks |

#### Host Operations (2)
| Tool Name | Description |
|-----------|-------------|
| `ambari_hosts_gethosts` | List all hosts in the cluster |
| `ambari_hosts_gethost` | Get detailed information about a specific host |

#### Alert Operations (7)
| Tool Name | Description |
|-----------|-------------|
| `ambari_alerts_getalerts` | Get all cluster alerts |
| `ambari_alerts_getalertsummary` | Get summarized alert information |
| `ambari_alerts_getalertdetails` | Get detailed alert information |
| `ambari_alerts_getalertdefinitions` | List all alert definitions |
| `ambari_alerts_getalertgroups` | List all alert groups |
| `ambari_alerts_gettargets` | List all alert targets |
| `ambari_alerts_getnotifications` | List all alert notification targets |

#### User & Group Operations (5)
| Tool Name | Description |
|-----------|-------------|
| `ambari_users_getusers` | Get all Ambari users |
| `ambari_users_getuser` | Get specific user details |
| `ambari_users_getgroups` | Get all Ambari groups |
| `ambari_users_getgroup` | Get specific group details |
| `ambari_users_getuserprivileges` | Get privileges assigned to a user |

### Actionable Tools (27) - Require ENABLE_ACTIONABLE_TOOLS=true

These tools perform state-changing operations and require explicit enablement:

#### Cluster Management (1)
| Tool Name | Description |
|-----------|-------------|
| `ambari_clusters_createcluster` | Create a new Ambari cluster |

#### Service Lifecycle (7)
| Tool Name | Description |
|-----------|-------------|
| `ambari_services_startservice` | Start a service |
| `ambari_services_stopservice` | Stop a service |
| `ambari_services_restartservice` | Restart a service |
| `ambari_services_restartcomponents` | Restart specific components |
| `ambari_services_enablemaintenancemode` | Enable maintenance mode for a service |
| `ambari_services_disablemaintenancemode` | Disable maintenance mode for a service |
| `ambari_services_runservicecheck` | Run health checks for a service |

#### Alert Definition Management (1)
| Tool Name | Description |
|-----------|-------------|
| `ambari_alerts_updatealertdefinition` | Update an alert definition |

#### Alert Group Management (5)
| Tool Name | Description |
|-----------|-------------|
| `ambari_alerts_createalertgroup` | Create a new alert group |
| `ambari_alerts_updatealertgroup` | Update an existing alert group |
| `ambari_alerts_deletealertgroup` | Delete an alert group |
| `ambari_alerts_duplicatealertgroup` | Duplicate an alert group |
| `ambari_alerts_adddefinitiontogroup` | Add an alert definition to a group |
| `ambari_alerts_removedefinitionfromgroup` | Remove an alert definition from a group |

#### Alert Notification Management (5)
| Tool Name | Description |
|-----------|-------------|
| `ambari_alerts_createnotification` | Create a new notification target |
| `ambari_alerts_updatenotification` | Update a notification target |
| `ambari_alerts_deletenotification` | Delete a notification target |
| `ambari_alerts_addnotificationtogroup` | Add notification to alert group |
| `ambari_alerts_removenotificationfromgroup` | Remove notification from alert group |
| `ambari_alerts_savealertsettings` | Save alert settings |

#### User & Group Management (8)
| Tool Name | Description |
|-----------|-------------|
| `ambari_users_createuser` | Create a new Ambari user |
| `ambari_users_updateuser` | Update an existing user |
| `ambari_users_deleteuser` | Delete an Ambari user |
| `ambari_users_creategroup` | Create a new Ambari group |
| `ambari_users_deletegroup` | Delete an Ambari group |
| `ambari_users_addusertogroup` | Add a user to a group |
| `ambari_users_removeuserfromgroup` | Remove a user from a group |

---

## 📊 Resources (12 Total)

Resources provide read-only access to cluster data via URI patterns:

| Resource URI | Description |
|--------------|-------------|
| `ambari://clusters` | List of all Ambari clusters with basic information |
| `ambari://cluster/{clusterName}` | Detailed information about a specific cluster |
| `ambari://cluster/{clusterName}/services` | All services running in a cluster with their status |
| `ambari://cluster/{clusterName}/hosts` | All hosts in a cluster with status and components |
| `ambari://cluster/{clusterName}/alerts` | Current alerts for a cluster grouped by severity |
| `ambari://cluster/{clusterName}/alerts/summary` | Summarized alert information for quick health overview |
| `ambari://cluster/{clusterName}/services/stale-configs` | Services needing restart due to configuration changes |
| `ambari://cluster/{clusterName}/service/{serviceName}` | Detailed information about a specific service |
| `ambari://cluster/{clusterName}/service/{serviceName}/components` | All components of a service with host assignments |
| `ambari://host/{hostName}` | Detailed information about a specific host |
| `ambari://cluster/{clusterName}/requests/recent` | Recent operations and their status |
| `ambari://cluster/{clusterName}/configurations` | Current configuration types for all services |

### Example Resource Access

```bash
# Get cluster details
URI: ambari://cluster/sagarautomation

# Get service components
URI: ambari://cluster/sagarautomation/service/HDFS/components

# Get alert summary
URI: ambari://cluster/sagarautomation/alerts/summary
```

---

## 📝 Prompts (8 Total)

Prompts are reusable templates that guide AI agents through multi-step workflows:

### 1. cluster_health_check
**Description**: Comprehensive health check for an Ambari cluster including services, hosts, and alerts

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster to check

**Workflow**: Checks cluster overview, service status, alerts, host health, and stale configurations, then provides overall health assessment and recommendations.

---

### 2. service_troubleshooting
**Description**: Troubleshoot issues with a specific Ambari service

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster
- `serviceName` - Name of the service to troubleshoot (e.g., HDFS, YARN)

**Workflow**: Analyzes service state, alerts, components, and configurations to identify issues and provide troubleshooting recommendations.

---

### 3. alert_investigation
**Description**: Investigate and analyze alerts in an Ambari cluster

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster

**Optional Arguments**:
- `severity` - Alert severity to investigate (CRITICAL, WARNING, OK)

**Workflow**: Reviews alert summary, details, affected services/hosts, and provides root cause analysis with remediation actions.

---

### 4. performance_analysis
**Description**: Analyze performance and resource usage of cluster services

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster

**Optional Arguments**:
- `serviceName` - Specific service to analyze

**Workflow**: Evaluates service status, host resource usage, bottlenecks, and provides performance optimization recommendations.

---

### 5. configuration_review
**Description**: Review cluster configurations and identify potential issues

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster

**Workflow**: Identifies stale configurations, components needing restart, and recommends restart strategies with impact assessment.

---

### 6. user_permissions_audit
**Description**: Audit Ambari users, groups, and permissions

**Required Arguments**: None

**Workflow**: Reviews all users, groups, privileges, and provides security audit report with access control recommendations.

---

### 7. upgrade_readiness_check
**Description**: Check if cluster is ready for upgrade or maintenance

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster

**Workflow**: Verifies all services are running, no critical alerts, no stale configs, all hosts healthy, then provides readiness assessment and upgrade blockers.

---

### 8. service_dependency_analysis
**Description**: Analyze service dependencies and start/stop order

**Required Arguments**:
- `clusterName` - Name of the Ambari cluster
- `serviceName` - Service to analyze dependencies for

**Workflow**: Identifies service dependencies, recommended start/stop order, impact analysis, and safe restart procedures.

---

### 📝 **8 MCP Prompts Available**

Reusable prompt templates for common Ambari workflows:

- `cluster_health_check` - Comprehensive cluster health analysis
- `service_troubleshooting` - Diagnose and fix service issues
- `alert_investigation` - Investigate and analyze cluster alerts
- `performance_analysis` - Analyze service and cluster performance
- `configuration_review` - Review and validate cluster configurations
- `user_permissions_audit` - Audit users, groups, and permissions
- `upgrade_readiness_check` - Verify cluster is ready for maintenance
- `service_dependency_analysis` - Understand service dependencies

### 📊 **12 MCP Resources Available**

Direct access to cluster data via URI patterns:

- `ambari://clusters` - List all clusters
- `ambari://cluster/{clusterName}` - Cluster details
- `ambari://cluster/{clusterName}/services` - Cluster services
- `ambari://cluster/{clusterName}/hosts` - Cluster hosts  
- `ambari://cluster/{clusterName}/alerts` - Cluster alerts
- `ambari://cluster/{clusterName}/alerts/summary` - Alert summary
- `ambari://cluster/{clusterName}/services/stale-configs` - Stale configurations
- `ambari://cluster/{clusterName}/service/{serviceName}` - Service details
- `ambari://cluster/{clusterName}/service/{serviceName}/components` - Service components
- `ambari://host/{hostName}` - Host details
- `ambari://cluster/{clusterName}/requests/recent` - Recent operations
- `ambari://cluster/{clusterName}/configurations` - Configuration types

## Installation

### Prerequisites

- **Go 1.23+** (with Go 1.24 toolchain)
- **Access to Apache Ambari cluster**

### Build from Source

```bash
# Clone the repository
git clone https://mcp-ambari.git
cd mcp-ambari

# Install dependencies
go mod download

# Build the server
go build ./cmd/server

# The binary 'server' is now ready to use
```

### Binary Installation

```bash
# Build and install to $GOPATH/bin
go install ./cmd/server

# Or build locally
make build
```

## Configuration

### Environment Variables

```bash
# Copy the example configuration
cp .env.example .env

# Required: Ambari connection
export AMBARI_BASE_URL=http://your-ambari-server:8080/api/v1
export AMBARI_USERNAME=admin  
export AMBARI_PASSWORD=your-password

# Optional: Timeouts and logging
export AMBARI_TIMEOUT=30s
export LOG_LEVEL=info

# Optional: Authentication (for HTTP transport)
export AUTH_ENABLED=false
export LDAP_HEADER_PREFIX=x-user-
export DEFAULT_PERMISSIONS=cluster:view,service:view

# Optional: Transport mode
export MCP_TRANSPORT=stdio  # Options: stdio, http, ssl, mtls
```

### Configuration Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `AMBARI_BASE_URL` | Ambari REST API endpoint | `http://localhost:8080/api/v1` | ✅ |
| `AMBARI_USERNAME` | Ambari username | `admin` | ✅ |
| `AMBARI_PASSWORD` | Ambari password | `admin` | ✅ |
| `AMBARI_TIMEOUT` | Request timeout | `30s` | ❌ |
| `LOG_LEVEL` | Logging level | `info` | ❌ |
| `MCP_TRANSPORT` | Transport mode | `stdio` | ❌ |
| `AUTH_ENABLED` | Enable authentication | `false` | ❌ |

## Usage

### With MCP Clients

#### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

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

#### Cline (VS Code Extension)

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

### Manual Testing

```bash
# Run the server directly
./server -transport stdio

# With custom configuration
AMBARI_BASE_URL=http://your-ambari:8080/api/v1 \
AMBARI_USERNAME=admin \
AMBARI_PASSWORD=your-password \
LOG_LEVEL=debug \
./server -transport stdio
```

### Command Line Options

```bash
./server [OPTIONS]

Options:
  -transport string
        Transport mode: stdio, http, ssl, mtls (default "stdio")
  -host string
        Server host for HTTP modes (default "0.0.0.0")  
  -port string
        Server port for HTTP modes (default "9001")
  -ssl-certfile string
        SSL certificate file (default "certs/server-cert.pem")
  -ssl-keyfile string
        SSL private key file (default "certs/server-key.pem")
  -ssl-ca-certs string
        CA certs for mTLS client verification (default "certs/ca.pem")
```

## Project Structure

```
mcp-ambari/
├── cmd/
│   └── server/                 # Main application entry point
│       └── main.go            
├── internal/                   # Private application code
│   ├── auth/                  # Authentication & authorization
│   │   └── auth.go            # LDAP provider, permissions, middleware
│   ├── client/                # Ambari REST client
│   │   └── ambari.go          # HTTP client with pooling & retries
│   ├── operations/            # Business logic layer
│   │   ├── base.go            # Base interfaces & executor
│   │   ├── registry.go        # Operation registry & factory
│   │   ├── actionable/        # State-changing operations
│   │   │   ├── alerts.go      # Alert management operations
│   │   │   └── services.go    # Service lifecycle operations
│   │   └── readonly/          # Safe, read-only operations
│   │       ├── alerts.go      # Alert querying operations
│   │       └── clusters.go    # Cluster & service queries
│   │   │   └── users.go       # User and group queries
│   ├── resources/             # MCP resources (data endpoints)
│   │   └── resources.go       # 12 cluster data resources
│   └── transport/             # Transport layer abstraction
│       └── transport.go       # Stdio/HTTP/TLS transport modes
├── .env.example              # Configuration template
├── MCP_SETUP.md             # Detailed setup guide
├── go.mod                   # Go module definition
└── README.md               # This file
```

### Architecture Highlights

- **Separation of Concerns**: Read-only vs actionable operations
- **Template Method Pattern**: Standardized execution lifecycle
- **Strategy Pattern**: Pluggable auth providers and transports  
- **Registry Pattern**: Dynamic operation management
- **Dependency Injection**: Clean, testable architecture

## Example: Cluster Information

Based on your connected cluster `sagarautomation`:

```json
{
  "cluster_name": "sagarautomation",
  "cluster_id": 2,
  "version": "VDP-3.4", 
  "total_hosts": 3,
  "security_type": "KERBEROS",
  "provisioning_state": "INSTALLED",
  "health_report": {
    "Host/host_status/HEALTHY": 2,
    "Host/host_status/UNHEALTHY": 1,
    "Host/stale_config": 1
  }
}
```

**Installed Services**: HDFS, YARN, HIVE, HBASE, SPARK3, RANGER, RANGER_KMS, AMBARI_METRICS, MAPREDUCE2, ZOOKEEPER, KERBEROS, TEZ

## Authentication & Security

### Permission System

The server implements a comprehensive permission system:

```go
// Available Permissions
ClusterView, ClusterAdmin, ServiceView, ServiceOperate, 
ServiceRestart, ServiceAdmin, HostView, HostManage,
AlertView, AlertManage, AlertAdmin, ConfigView, ConfigModify

// Permission Groups  
"ADMIN":    All permissions
"OPERATOR": View, operate, and restart permissions
"VIEWER":   Read-only permissions only
```

### LDAP Integration

```bash
export AUTH_ENABLED=true
export LDAP_HEADER_PREFIX=x-user-
```

Headers expected:
- `x-user-name` or `x-user-username`: Username
- `x-user-groups`: Comma-separated group list

### Group Mappings

```go
"ambari-admins":    Full admin access
"hadoop-operators": Operational permissions
"data-engineers":   View and operate services  
"bigdata-viewers":  Read-only access
```

## Transport Modes

### Stdio (Default)
For MCP clients like Claude Desktop and Cline:
```bash
./server -transport stdio
```

### HTTP
For web applications and streamableHttp clients:
```bash  
./server -transport http -host 127.0.0.1 -port 8094
```

### HTTPS/TLS
For secure deployments with TLS encryption:
```bash
# Set TLS certificate environment variables
export TLS_CERT_FILE=/path/to/server.crt
export TLS_KEY_FILE=/path/to/server.key

./server -transport ssl -host 127.0.0.1 -port 8443
```

### HTTPS/mTLS
For enterprise deployments with mutual TLS authentication:
```bash
# Set TLS certificate and CA environment variables
export TLS_CERT_FILE=/path/to/server.crt
export TLS_KEY_FILE=/path/to/server.key
export TLS_CA_FILE=/path/to/ca.crt

./server -transport mtls -host 127.0.0.1 -port 8443
```

### Actionable Tool Control
**Actionable tools are disabled by default for security.** Only read-only operations are available unless explicitly enabled:
```bash
# Default: Only readonly tools (24 tools - safe operations only)
./server -transport http -port 8094

# Explicitly enable all tools (51 tools - includes user management, service control, etc.)
export ENABLE_ACTIONABLE_TOOLS=true
./server -transport http -port 8094
```

## Error Handling & Reliability

- **Retry Logic**: Automatic retry with exponential backoff
- **Connection Pooling**: Efficient HTTP connection reuse
- **Graceful Shutdown**: Clean resource cleanup on termination
- **Comprehensive Logging**: Structured JSON logging with correlation IDs
- **Input Validation**: Parameter validation before execution
- **Permission Checks**: Authorization validation for all operations

## Development

### Prerequisites

- Go 1.23+ with Go 1.24 toolchain
- Access to Apache Ambari cluster
- Optional: Docker for containerized deployment

### Local Development

```bash
# Install dependencies
go mod download

# Run with development settings
export AMBARI_BASE_URL=http://localhost:8080/api/v1
export AMBARI_USERNAME=admin
export AMBARI_PASSWORD=admin
export LOG_LEVEL=debug

# Build and run
go run ./cmd/server -transport stdio
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Adding New Operations

1. **Create the operation struct** implementing the `Operation` interface
2. **Add to the registry** in `main.go`
3. **Implement required methods**: `Name()`, `Description()`, `Definition()`, `Validate()`, `Execute()`
4. **Add proper permissions** and error handling

Example:
```go
type GetNewData struct {
    ops.ReadOnlyBase
}

func (o *GetNewData) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    return o.Client.Get(ctx, "/new-endpoint", params)
}
```

## Deployment

### Binary Deployment

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" ./cmd/server

# Deploy binary
scp server user@server:/usr/local/bin/mcp-ambari
```

### Docker Deployment (Future)

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

## Monitoring & Observability

### Structured Logging

All operations log with structured JSON format:

```json
{
  "level": "info",
  "msg": "Operation completed",
  "tool": "ambari_services_getservices",
  "type": "readonly", 
  "execution_ms": 245,
  "timestamp": "2026-02-08T13:39:44Z"
}
```

### Performance Metrics

- **Operation execution times** tracked
- **Error rates** by operation type
- **Authentication success/failure** rates
- **Ambari API response times**

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|--------|----------|
| Connection refused | Ambari server not accessible | Check `AMBARI_BASE_URL` |
| Authentication failed | Invalid credentials | Verify `AMBARI_USERNAME`/`AMBARI_PASSWORD` |
| Permission denied | Insufficient Ambari permissions | Use admin account or grant permissions |
| Timeout errors | Network latency | Increase `AMBARI_TIMEOUT` |
| MCP client not connecting | Configuration issues | Check client config syntax |

### Debug Mode

```bash
export LOG_LEVEL=debug
./server -transport stdio
```

### Health Check

Test Ambari connectivity:
```bash
curl -u admin:password http://your-ambari:8080/api/v1/clusters
```

## Performance

- **Concurrent Operations**: Multiple operations can run simultaneously  
- **Connection Pooling**: Efficient HTTP connection reuse
- **Memory Efficient**: Streaming JSON parsing for large responses
- **Fast Startup**: Sub-second initialization time
- **Low Latency**: Direct REST API access without additional layers

## Contributing

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`) 
5. **Open** a Pull Request

### Coding Standards

- Follow Go conventions and `gofmt` formatting
- Add tests for new functionality  
- Update documentation for API changes
- Use meaningful commit messages
- Ensure SOLID principles are maintained

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Model Context Protocol](https://modelcontextprotocol.io/) - Official MCP specification
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official Go SDK for MCP
- [Apache Ambari](https://ambari.apache.org/) - Apache Ambari project

## Support

For support and questions:
1. Check the [MCP_SETUP.md](MCP_SETUP.md) guide
2. Review the troubleshooting section above
3. Open an issue in the repository

---

**Built with ❤️ and Go for the Apache Ambari community**
