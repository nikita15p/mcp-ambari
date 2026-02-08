package readonly

import (
	"context"
	"fmt"

	"github.com/niita15p/mcp-ambari/internal/auth"
	"github.com/niita15p/mcp-ambari/internal/client"
	ops "github.com/niita15p/mcp-ambari/internal/operations"
	"github.com/sirupsen/logrus"
)

// ---- GetHost ----
type GetHost struct{ ops.ReadOnlyBase }

func NewGetHost(c client.AmbariClient, l *logrus.Logger) *GetHost {
	return &GetHost{ops.ReadOnlyBase{OpName: "ambari_hosts_gethost", OpDescription: "Returns information about a single host", OpCategory: "hosts", Permissions: []auth.Permission{auth.HostView}, Client: c, Logger: l}}
}
func (o *GetHost) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"hostName": map[string]interface{}{"type": "string", "description": "The name of the host"}, "fields": map[string]interface{}{"type": "string", "description": "Filter fields", "default": "Hosts/*"}}, Required: []string{"hostName"}}}
}
func (o *GetHost) Validate(args map[string]interface{}) error {
	if _, ok := args["hostName"].(string); !ok {
		return fmt.Errorf("hostName is required")
	}
	return nil
}
func (o *GetHost) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	p := map[string]string{"fields": "Hosts/*"}
	if f, ok := args["fields"].(string); ok {
		p["fields"] = f
	}
	return o.Client.Get(ctx, fmt.Sprintf("/hosts/%s", args["hostName"].(string)), p)
}

// ---- GetAlertTargets ----
type GetAlertTargets struct{ ops.ReadOnlyBase }

func NewGetAlertTargets(c client.AmbariClient, l *logrus.Logger) *GetAlertTargets {
	return &GetAlertTargets{ops.ReadOnlyBase{OpName: "ambari_alerts_gettargets", OpDescription: "Returns all alert targets", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l}}
}
func (o *GetAlertTargets) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"fields": map[string]interface{}{"type": "string", "description": "Filter fields", "default": "AlertTarget/*"}}, Required: []string{}}}
}
func (o *GetAlertTargets) Validate(args map[string]interface{}) error { return nil }
func (o *GetAlertTargets) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	p := map[string]string{"fields": "AlertTarget/*"}
	if f, ok := args["fields"].(string); ok {
		p["fields"] = f
	}
	return o.Client.Get(ctx, "/alert_targets", p)
}

// ---- GetAlertSummary ----
type GetAlertSummary struct{ ops.ReadOnlyBase }

func NewGetAlertSummary(c client.AmbariClient, l *logrus.Logger) *GetAlertSummary {
	return &GetAlertSummary{ops.ReadOnlyBase{OpName: "ambari_alerts_getalertsummary", OpDescription: "Get alert summary in grouped format for a cluster", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l}}
}
func (o *GetAlertSummary) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "maintenanceFilter": map[string]interface{}{"type": "boolean", "description": "Filter out maintenance alerts"}}, Required: []string{"clusterName"}}}
}
func (o *GetAlertSummary) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	return nil
}
func (o *GetAlertSummary) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	p := map[string]string{"format": "groupedSummary"}
	if mf, ok := args["maintenanceFilter"].(bool); ok && mf {
		p["Alert/maintenance_state.in"] = "OFF"
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/alerts", args["clusterName"].(string)), p)
}

// ---- GetAlertDetails ----
type GetAlertDetails struct{ ops.ReadOnlyBase }

func NewGetAlertDetails(c client.AmbariClient, l *logrus.Logger) *GetAlertDetails {
	return &GetAlertDetails{ops.ReadOnlyBase{OpName: "ambari_alerts_getalertdetails", OpDescription: "Get details for a specific alert definition", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l}}
}
func (o *GetAlertDetails) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "alertId": map[string]interface{}{"type": "string", "description": "Alert definition ID"}}, Required: []string{"clusterName", "alertId"}}}
}
func (o *GetAlertDetails) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	if _, ok := args["alertId"].(string); !ok {
		return fmt.Errorf("alertId required")
	}
	return nil
}
func (o *GetAlertDetails) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	p := map[string]string{"fields": "*", "Alert/definition_id": args["alertId"].(string)}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/alerts", args["clusterName"].(string)), p)
}

// ---- GetAlertDefinitions ----
type GetAlertDefinitions struct{ ops.ReadOnlyBase }

func NewGetAlertDefinitions(c client.AmbariClient, l *logrus.Logger) *GetAlertDefinitions {
	return &GetAlertDefinitions{ops.ReadOnlyBase{OpName: "ambari_alerts_getalertdefinitions", OpDescription: "Get all alert definitions for a cluster", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l}}
}
func (o *GetAlertDefinitions) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "fields": map[string]interface{}{"type": "string", "description": "Filter fields", "default": "*"}}, Required: []string{"clusterName"}}}
}
func (o *GetAlertDefinitions) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetAlertDefinitions) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	p := map[string]string{"fields": "*"}
	if f, ok := args["fields"].(string); ok {
		p["fields"] = f
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/alert_definitions", args["clusterName"].(string)), p)
}

// ---- GetAlertGroups ----
type GetAlertGroups struct{ ops.ReadOnlyBase }

func NewGetAlertGroups(c client.AmbariClient, l *logrus.Logger) *GetAlertGroups {
	return &GetAlertGroups{ops.ReadOnlyBase{OpName: "ambari_alerts_getalertgroups", OpDescription: "Get all alert groups for a cluster", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l}}
}
func (o *GetAlertGroups) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}}, Required: []string{"clusterName"}}}
}
func (o *GetAlertGroups) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetAlertGroups) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/alert_groups", args["clusterName"].(string)), map[string]string{"fields": "*"})
}

// ---- GetNotifications ----
type GetNotifications struct{ ops.ReadOnlyBase }

func NewGetNotifications(c client.AmbariClient, l *logrus.Logger) *GetNotifications {
	return &GetNotifications{ops.ReadOnlyBase{OpName: "ambari_alerts_getnotifications", OpDescription: "Get all alert notification targets", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l}}
}
func (o *GetNotifications) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}}, Required: []string{"clusterName"}}}
}
func (o *GetNotifications) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetNotifications) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return o.Client.Get(ctx, "/alert_targets", map[string]string{"fields": "*"})
}

// ---- GetServicesWithStaleConfigs ----
type GetServicesWithStaleConfigs struct{ ops.ReadOnlyBase }

func NewGetServicesWithStaleConfigs(c client.AmbariClient, l *logrus.Logger) *GetServicesWithStaleConfigs {
	return &GetServicesWithStaleConfigs{ops.ReadOnlyBase{OpName: "ambari_services_getserviceswithstaleconfigs", OpDescription: "Get services with stale configurations requiring restart", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l}}
}
func (o *GetServicesWithStaleConfigs) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "serviceName": map[string]interface{}{"type": "string", "description": "Filter by service (optional)"}}, Required: []string{"clusterName"}}}
}
func (o *GetServicesWithStaleConfigs) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetServicesWithStaleConfigs) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster := args["clusterName"].(string)
	p := map[string]string{"fields": "ServiceInfo/service_name,ServiceInfo/state,ServiceInfo/maintenance_state,components/ServiceComponentInfo/component_name,components/host_components/HostRoles/state,components/host_components/HostRoles/stale_configs,components/host_components/HostRoles/host_name"}
	if svc, ok := args["serviceName"].(string); ok && svc != "" {
		return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, svc), p)
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/services", cluster), p)
}

// ---- GetHostComponentsWithStaleConfigs ----
type GetHostComponentsWithStaleConfigs struct{ ops.ReadOnlyBase }

func NewGetHostComponentsWithStaleConfigs(c client.AmbariClient, l *logrus.Logger) *GetHostComponentsWithStaleConfigs {
	return &GetHostComponentsWithStaleConfigs{ops.ReadOnlyBase{OpName: "ambari_services_gethostcomponentswithstaleconfigs", OpDescription: "Get host components needing restart due to stale configurations", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l}}
}
func (o *GetHostComponentsWithStaleConfigs) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "hostName": map[string]interface{}{"type": "string", "description": "Filter by host"}, "serviceName": map[string]interface{}{"type": "string", "description": "Filter by service"}, "componentName": map[string]interface{}{"type": "string", "description": "Filter by component"}}, Required: []string{"clusterName"}}}
}
func (o *GetHostComponentsWithStaleConfigs) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetHostComponentsWithStaleConfigs) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	p := map[string]string{"fields": "HostRoles/component_name,HostRoles/host_name,HostRoles/service_name,HostRoles/state,HostRoles/stale_configs,HostRoles/maintenance_state", "HostRoles/stale_configs": "true"}
	if h, ok := args["hostName"].(string); ok {
		p["HostRoles/host_name"] = h
	}
	if s, ok := args["serviceName"].(string); ok {
		p["HostRoles/service_name"] = s
	}
	if c, ok := args["componentName"].(string); ok {
		p["HostRoles/component_name"] = c
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/host_components", args["clusterName"].(string)), p)
}

// ---- GetRollingRestartStatus ----
type GetRollingRestartStatus struct{ ops.ReadOnlyBase }

func NewGetRollingRestartStatus(c client.AmbariClient, l *logrus.Logger) *GetRollingRestartStatus {
	return &GetRollingRestartStatus{ops.ReadOnlyBase{OpName: "ambari_services_getrollingrestartstatus", OpDescription: "Get status of rolling restart operations", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l}}
}
func (o *GetRollingRestartStatus) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "serviceName": map[string]interface{}{"type": "string", "description": "Filter by service"}, "requestId": map[string]interface{}{"type": "string", "description": "Filter by request ID"}}, Required: []string{"clusterName"}}}
}
func (o *GetRollingRestartStatus) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetRollingRestartStatus) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster := args["clusterName"].(string)
	p := map[string]string{"fields": "Requests/id,Requests/request_context,Requests/request_status,Requests/progress_percent,Requests/start_time,Requests/end_time,tasks/Tasks/command_name,tasks/Tasks/status,tasks/Tasks/host_name,tasks/Tasks/role"}
	if s, ok := args["serviceName"].(string); ok {
		p["tasks/Tasks/role.in"] = s
	}
	if rid, ok := args["requestId"].(string); ok && rid != "" {
		return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/requests/%s", cluster, rid), p)
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/requests", cluster), p)
}

// ---- IsServiceCheckSupported ----
type IsServiceCheckSupported struct{ ops.ReadOnlyBase }

func NewIsServiceCheckSupported(c client.AmbariClient, l *logrus.Logger) *IsServiceCheckSupported {
	return &IsServiceCheckSupported{ops.ReadOnlyBase{OpName: "ambari_services_isservicechecksupported", OpDescription: "Check if service check is supported for a service in the stack", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l}}
}
func (o *IsServiceCheckSupported) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "serviceName": map[string]interface{}{"type": "string", "description": "Service name"}, "stackName": map[string]interface{}{"type": "string", "description": "Stack name (e.g., HDP, VDP)"}, "stackVersion": map[string]interface{}{"type": "string", "description": "Stack version (e.g., 3.1)"}}, Required: []string{"clusterName", "serviceName", "stackName", "stackVersion"}}}
}
func (o *IsServiceCheckSupported) Validate(args map[string]interface{}) error {
	for _, k := range []string{"clusterName", "serviceName", "stackName", "stackVersion"} {
		if _, ok := args[k].(string); !ok {
			return fmt.Errorf("%s required", k)
		}
	}
	return nil
}
func (o *IsServiceCheckSupported) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return o.Client.Get(ctx, fmt.Sprintf("/stacks/%s/versions/%s/services/%s", args["stackName"].(string), args["stackVersion"].(string), args["serviceName"].(string)), map[string]string{"fields": "StackServices/service_check_supported"})
}

// ---- GetServiceCheckStatus ----
type GetServiceCheckStatus struct{ ops.ReadOnlyBase }

func NewGetServiceCheckStatus(c client.AmbariClient, l *logrus.Logger) *GetServiceCheckStatus {
	return &GetServiceCheckStatus{ops.ReadOnlyBase{OpName: "ambari_services_getservicecheckstatus", OpDescription: "Get status of recent service check operations", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l}}
}
func (o *GetServiceCheckStatus) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"}, "serviceName": map[string]interface{}{"type": "string", "description": "Filter by service"}, "requestId": map[string]interface{}{"type": "string", "description": "Filter by request ID"}}, Required: []string{"clusterName"}}}
}
func (o *GetServiceCheckStatus) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName required")
	}
	return nil
}
func (o *GetServiceCheckStatus) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster := args["clusterName"].(string)
	p := map[string]string{"fields": "Requests/id,Requests/request_context,Requests/request_status,Requests/progress_percent,tasks/Tasks/command_name,tasks/Tasks/status,tasks/Tasks/host_name,tasks/Tasks/role", "Requests/request_context.matches": ".*Service Check.*"}
	if s, ok := args["serviceName"].(string); ok {
		p["tasks/Tasks/role.in"] = s
	}
	if rid, ok := args["requestId"].(string); ok && rid != "" {
		return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/requests/%s", cluster, rid), p)
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/requests", cluster), p)
}
