/* START GENAI */
// Package prompts defines MCP Prompts for the Ambari MCP Server.
// Prompts are reusable templates that guide AI agents through common
// Ambari operations and troubleshooting workflows.
package prompts

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
)

// PromptDefinition describes a single MCP prompt template
type PromptDefinition struct {
	Name        string
	Description string
	Arguments   []PromptArgument
	Template    string
}

// PromptArgument defines a required argument for a prompt
type PromptArgument struct {
	Name        string
	Description string
	Required    bool
}

// Registry holds all prompt definitions
type Registry struct {
	prompts []PromptDefinition
	logger  *logrus.Logger
}

// NewRegistry creates a prompt registry with all Ambari prompts
func NewRegistry(logger *logrus.Logger) *Registry {
	r := &Registry{
		prompts: []PromptDefinition{},
		logger:  logger,
	}
	r.registerAll()
	return r
}

// Definitions returns all prompt definitions
func (r *Registry) Definitions() []PromptDefinition {
	return r.prompts
}

// Get retrieves a prompt by name
func (r *Registry) Get(name string) (*PromptDefinition, error) {
	for _, p := range r.prompts {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("prompt not found: %s", name)
}

// Count returns number of registered prompts
func (r *Registry) Count() int {
	return len(r.prompts)
}

// GetPrompt renders a prompt with provided arguments
func (r *Registry) GetPrompt(name string, args map[string]string) (string, error) {
	prompt, err := r.Get(name)
	if err != nil {
		return "", err
	}

	// Validate required arguments
	for _, arg := range prompt.Arguments {
		if arg.Required {
			if _, ok := args[arg.Name]; !ok {
				return "", fmt.Errorf("missing required argument: %s", arg.Name)
			}
		}
	}

	// Simple template rendering (replace {argName} with values)
	result := prompt.Template
	for key, value := range args {
		placeholder := "{" + key + "}"
		result = replaceAll(result, placeholder, value)
	}

	return result, nil
}

func (r *Registry) registerAll() {
	// 1. Cluster Health Check
	r.add(PromptDefinition{
		Name:        "cluster_health_check",
		Description: "Comprehensive health check for an Ambari cluster including services, hosts, and alerts",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster to check", Required: true},
		},
		Template: `Perform a comprehensive health check for the Ambari cluster "{clusterName}":

1. Get cluster overview using ambari_clusters_getcluster
2. Check all services status using ambari_services_getservices
3. Review critical and warning alerts using ambari_alerts_getalertsummary
4. Check host health status using ambari_hosts_gethosts
5. Identify services with stale configurations using ambari_services_getserviceswithstaleconfigs

Analyze the results and provide:
- Overall cluster health status (Healthy/Degraded/Critical)
- List of services that are not running
- Critical alerts that need immediate attention
- Hosts with issues
- Services needing restart due to configuration changes
- Recommended actions to improve cluster health

Cluster: {clusterName}`,
	})

	// 2. Service Troubleshooting
	r.add(PromptDefinition{
		Name:        "service_troubleshooting",
		Description: "Troubleshoot issues with a specific Ambari service",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster", Required: true},
			{Name: "serviceName", Description: "Name of the service to troubleshoot (e.g., HDFS, YARN)", Required: true},
		},
		Template: `Troubleshoot the service "{serviceName}" in cluster "{clusterName}":

1. Get service current state using ambari_services_getservicestate
2. Check service-specific alerts using ambari_alerts_getalerts
3. Verify all components are running using ambari_services_getservice
4. Check for stale configurations using ambari_services_gethostcomponentswithstaleconfigs
5. Review recent operations related to this service

Analyze and provide:
- Current service status and health
- Components that are down or in maintenance mode
- Active alerts related to this service
- Recent configuration changes
- Components with stale configurations
- Recommended troubleshooting steps
- Commands to resolve common issues

Service: {serviceName}
Cluster: {clusterName}`,
	})

	// 3. Alert Investigation
	r.add(PromptDefinition{
		Name:        "alert_investigation",
		Description: "Investigate and analyze alerts in an Ambari cluster",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster", Required: true},
			{Name: "severity", Description: "Alert severity to investigate (CRITICAL, WARNING, OK)", Required: false},
		},
		Template: `Investigate alerts in cluster "{clusterName}"{severity_filter}:

1. Get alert summary using ambari_alerts_getalertsummary
2. List all alerts using ambari_alerts_getalerts
3. For each alert, get detailed information using ambari_alerts_getalertdetails
4. Check affected services and hosts
5. Review alert definitions using ambari_alerts_getalertdefinitions

Provide analysis including:
- Total alert count by severity
- Top 5 most critical alerts
- Services with the most alerts
- Hosts with the most alerts
- Patterns in alert occurrences
- Root cause analysis for critical alerts
- Recommended remediation actions
- Alerts that may require alert definition updates

Cluster: {clusterName}`,
	})

	// 4. Performance Analysis
	r.add(PromptDefinition{
		Name:        "performance_analysis",
		Description: "Analyze performance and resource usage of cluster services",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster", Required: true},
			{Name: "serviceName", Description: "Specific service to analyze (optional)", Required: false},
		},
		Template: `Analyze performance of cluster "{clusterName}"{service_filter}:

1. Get all services status using ambari_services_getservices
2. Check host resource usage using ambari_hosts_gethosts
3. Review service-specific metrics
4. Check for bottlenecks or resource constraints
5. Identify services with stale configs that may affect performance

Provide performance analysis:
- Services running vs stopped
- Host resource utilization (CPU, memory, disk)
- Services with potential performance issues
- Hosts under heavy load
- Configuration changes needed for optimization
- Services that need restart for performance improvements
- Capacity planning recommendations

Cluster: {clusterName}`,
	})

	// 5. Configuration Review
	r.add(PromptDefinition{
		Name:        "configuration_review",
		Description: "Review cluster configurations and identify potential issues",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster", Required: true},
		},
		Template: `Review configurations for cluster "{clusterName}":

1. Get cluster details using ambari_clusters_getcluster
2. Check services with stale configurations using ambari_services_getserviceswithstaleconfigs
3. Review components needing restart using ambari_services_gethostcomponentswithstaleconfigs
4. Check if rolling restart is in progress using ambari_services_getrollingrestartstatus

Provide configuration analysis:
- Total services with stale configurations
- Components requiring restart (grouped by service)
- Priority order for restart operations
- Impact assessment of pending configuration changes
- Recommended restart strategy (rolling vs full)
- Services that can be restarted independently
- Services with dependencies that need coordinated restart

Cluster: {clusterName}`,
	})

	// 6. User and Permissions Audit
	r.add(PromptDefinition{
		Name:        "user_permissions_audit",
		Description: "Audit Ambari users, groups, and permissions",
		Arguments: []PromptArgument{},
		Template: `Audit Ambari users and permissions:

1. List all users using ambari_users_getusers
2. List all groups using ambari_users_getgroups
3. For each user, check privileges using ambari_users_getuserprivileges
4. Review group memberships

Provide security audit report:
- Total number of users and groups
- Users with administrative privileges
- Users with no assigned groups
- Groups and their members
- Users with inactive or disabled accounts
- Privilege distribution across users
- Recommendations for access control improvements
- Users that may have excessive permissions

This helps ensure proper access control and security compliance.`,
	})

	// 7. Cluster Upgrade Readiness
	r.add(PromptDefinition{
		Name:        "upgrade_readiness_check",
		Description: "Check if cluster is ready for upgrade or maintenance",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster", Required: true},
		},
		Template: `Check upgrade readiness for cluster "{clusterName}":

1. Get cluster current version using ambari_clusters_getcluster
2. Check all services are running using ambari_services_getservices
3. Verify no critical alerts using ambari_alerts_getalertsummary
4. Ensure no stale configurations using ambari_services_getserviceswithstaleconfigs
5. Check all hosts are healthy using ambari_hosts_gethosts

Provide readiness assessment:
- Current cluster version and stack info
- Services not in STARTED state
- Critical or warning alerts present
- Services with pending configuration changes
- Unhealthy or unreachable hosts
- Recent failed operations
- Upgrade blockers (critical issues)
- Pre-upgrade recommendations
- Estimated downtime impact

Cluster: {clusterName}`,
	})

	// 8. Service Dependency Analysis
	r.add(PromptDefinition{
		Name:        "service_dependency_analysis",
		Description: "Analyze service dependencies and start/stop order",
		Arguments: []PromptArgument{
			{Name: "clusterName", Description: "Name of the Ambari cluster", Required: true},
			{Name: "serviceName", Description: "Service to analyze dependencies for", Required: true},
		},
		Template: `Analyze service dependencies for "{serviceName}" in cluster "{clusterName}":

1. Get service details using ambari_services_getservice
2. Check all services status using ambari_services_getservices
3. Identify dependent services

Provide dependency analysis:
- Services that {serviceName} depends on (must start first)
- Services that depend on {serviceName} (affected by {serviceName} outage)
- Recommended start order for related services
- Impact of stopping {serviceName}
- Safe restart procedure
- Related services that should be monitored during operation

This helps plan maintenance windows and understand service relationships.

Service: {serviceName}
Cluster: {clusterName}`,
	})

	r.logger.WithField("count", len(r.prompts)).Info("MCP prompts registered")
}

func (r *Registry) add(prompt PromptDefinition) {
	r.prompts = append(r.prompts, prompt)
}

// ToMCPPrompt converts our PromptDefinition to MCP SDK Prompt
func (p *PromptDefinition) ToMCPPrompt() *mcp.Prompt {
	arguments := []*mcp.PromptArgument{}
	for _, arg := range p.Arguments {
		arguments = append(arguments, &mcp.PromptArgument{
			Name:        arg.Name,
			Description: arg.Description,
			Required:    arg.Required,
		})
	}

	return &mcp.Prompt{
		Name:        p.Name,
		Description: p.Description,
		Arguments:   arguments,
	}
}

// Simple string replace function
func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx == -1 {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

/* END GENAI */
