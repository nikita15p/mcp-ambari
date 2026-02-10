// Package readonly contains all read-only (GET) Ambari operations
package readonly

import (
	"context"
	"fmt"

	"mcp-ambari/internal/auth"
	"mcp-ambari/internal/client"
	ops "mcp-ambari/internal/operations"
	"github.com/sirupsen/logrus"
)

// ---------- GetClusters ----------

type GetClusters struct {
	ops.ReadOnlyBase
}

func NewGetClusters(c client.AmbariClient, l *logrus.Logger) *GetClusters {
	return &GetClusters{ops.ReadOnlyBase{
		OpName: "ambari_clusters_getclusters", OpDescription: "Returns all clusters",
		OpCategory: "clusters", Permissions: []auth.Permission{auth.ClusterView}, Client: c, Logger: l,
	}}
}

func (o *GetClusters) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"fields":    map[string]interface{}{"type": "string", "description": "Filter fields", "default": "Clusters/*"},
			"page_size": map[string]interface{}{"type": "integer", "description": "Page size", "default": 10},
		}, Required: []string{}},
	}
}

func (o *GetClusters) Validate(args map[string]interface{}) error { return nil }

func (o *GetClusters) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := map[string]string{}
	if f, ok := args["fields"].(string); ok {
		params["fields"] = f
	} else {
		params["fields"] = "Clusters/*"
	}
	if ps, ok := args["page_size"].(float64); ok {
		params["page_size"] = fmt.Sprintf("%d", int(ps))
	}
	return o.Client.Get(ctx, "/clusters", params)
}

// ---------- GetCluster ----------

type GetCluster struct {
	ops.ReadOnlyBase
}

func NewGetCluster(c client.AmbariClient, l *logrus.Logger) *GetCluster {
	return &GetCluster{ops.ReadOnlyBase{
		OpName: "ambari_clusters_getcluster", OpDescription: "Returns information about a specific cluster",
		OpCategory: "clusters", Permissions: []auth.Permission{auth.ClusterView}, Client: c, Logger: l,
	}}
}

func (o *GetCluster) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"fields":      map[string]interface{}{"type": "string", "description": "Filter fields", "default": "Clusters/*"},
		}, Required: []string{"clusterName"}},
	}
}

func (o *GetCluster) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok || args["clusterName"] == "" {
		return fmt.Errorf("clusterName is required")
	}
	return nil
}

func (o *GetCluster) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster := args["clusterName"].(string)
	params := map[string]string{"fields": "Clusters/*"}
	if f, ok := args["fields"].(string); ok {
		params["fields"] = f
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s", cluster), params)
}

// ---------- GetServices ----------

type GetServices struct {
	ops.ReadOnlyBase
}

func NewGetServices(c client.AmbariClient, l *logrus.Logger) *GetServices {
	return &GetServices{ops.ReadOnlyBase{
		OpName: "ambari_services_getservices", OpDescription: "Get all services for a cluster",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l,
	}}
}

func (o *GetServices) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"fields":      map[string]interface{}{"type": "string", "description": "Filter fields"},
		}, Required: []string{"clusterName"}},
	}
}

func (o *GetServices) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok || args["clusterName"] == "" {
		return fmt.Errorf("clusterName is required")
	}
	return nil
}

func (o *GetServices) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster := args["clusterName"].(string)
	params := map[string]string{"fields": "ServiceInfo/service_name,ServiceInfo/cluster_name"}
	if f, ok := args["fields"].(string); ok {
		params["fields"] = f
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/services", cluster), params)
}

// ---------- GetService ----------

type GetService struct {
	ops.ReadOnlyBase
}

func NewGetService(c client.AmbariClient, l *logrus.Logger) *GetService {
	return &GetService{ops.ReadOnlyBase{
		OpName: "ambari_services_getservice", OpDescription: "Get details of a service",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l,
	}}
}

func (o *GetService) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
			"fields":      map[string]interface{}{"type": "string", "description": "Filter fields", "default": "ServiceInfo/*"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *GetService) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *GetService) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	params := map[string]string{"fields": "ServiceInfo/*"}
	if f, ok := args["fields"].(string); ok {
		params["fields"] = f
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, service), params)
}

// ---------- GetHosts ----------

type GetHosts struct {
	ops.ReadOnlyBase
}

func NewGetHosts(c client.AmbariClient, l *logrus.Logger) *GetHosts {
	return &GetHosts{ops.ReadOnlyBase{
		OpName: "ambari_hosts_gethosts", OpDescription: "Returns all hosts",
		OpCategory: "hosts", Permissions: []auth.Permission{auth.HostView}, Client: c, Logger: l,
	}}
}

func (o *GetHosts) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"fields": map[string]interface{}{"type": "string", "description": "Filter fields", "default": "Hosts/*"},
		}, Required: []string{}},
	}
}

func (o *GetHosts) Validate(args map[string]interface{}) error { return nil }

func (o *GetHosts) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := map[string]string{"fields": "Hosts/*"}
	if f, ok := args["fields"].(string); ok {
		params["fields"] = f
	}
	return o.Client.Get(ctx, "/hosts", params)
}

// ---------- GetAlerts ----------

type GetAlerts struct {
	ops.ReadOnlyBase
}

func NewGetAlerts(c client.AmbariClient, l *logrus.Logger) *GetAlerts {
	return &GetAlerts{ops.ReadOnlyBase{
		OpName: "ambari_alerts_getalerts", OpDescription: "Get all alerts for a cluster",
		OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertView}, Client: c, Logger: l,
	}}
}

func (o *GetAlerts) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"state":       map[string]interface{}{"type": "string", "description": "Filter by state (CRITICAL,WARNING,OK,UNKNOWN)"},
		}, Required: []string{"clusterName"}},
	}
}

func (o *GetAlerts) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	return nil
}

func (o *GetAlerts) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster := args["clusterName"].(string)
	params := map[string]string{"fields": "*"}
	if s, ok := args["state"].(string); ok {
		params["Alert/state"] = s
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/alerts", cluster), params)
}

// ---------- GetServiceState ----------

type GetServiceState struct {
	ops.ReadOnlyBase
}

func NewGetServiceState(c client.AmbariClient, l *logrus.Logger) *GetServiceState {
	return &GetServiceState{ops.ReadOnlyBase{
		OpName: "ambari_services_getservicestate", OpDescription: "Get detailed state of a service",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceView}, Client: c, Logger: l,
	}}
}

func (o *GetServiceState) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *GetServiceState) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *GetServiceState) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	params := map[string]string{
		"fields": "ServiceInfo/*,components/ServiceComponentInfo/*,components/host_components/HostRoles/state,components/host_components/HostRoles/stale_configs",
	}
	return o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, service), params)
}
