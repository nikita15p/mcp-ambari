// Package actionable contains all state-changing (POST/PUT/DELETE) Ambari operations
package actionable

import (
	"context"
	"fmt"

	"mcp-ambari/internal/auth"
	"mcp-ambari/internal/client"
	ops "mcp-ambari/internal/operations"
	"github.com/sirupsen/logrus"
)

// ---------- StartService ----------

type StartService struct {
	ops.ActionableBase
}

func NewStartService(c client.AmbariClient, l *logrus.Logger) *StartService {
	return &StartService{ops.ActionableBase{
		OpName: "ambari_services_startservice", OpDescription: "Start a specific service on the cluster",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceOperate}, Dangerous: false, Client: c, Logger: l,
	}}
}

func (o *StartService) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
			"context":     map[string]interface{}{"type": "string", "description": "Context message", "default": "Start service via MCP"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *StartService) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *StartService) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	ctxMsg := "Start service via MCP"
	if c, ok := args["context"].(string); ok {
		ctxMsg = c
	}
	body := map[string]interface{}{
		"RequestInfo": map[string]interface{}{
			"context": ctxMsg,
			"operation_level": map[string]interface{}{
				"level": "SERVICE", "cluster_name": cluster, "service_name": service,
			},
		},
		"Body": map[string]interface{}{
			"ServiceInfo": map[string]interface{}{"state": "STARTED"},
		},
	}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, service), nil, body)
}

// ---------- StopService ----------

type StopService struct {
	ops.ActionableBase
}

func NewStopService(c client.AmbariClient, l *logrus.Logger) *StopService {
	return &StopService{ops.ActionableBase{
		OpName: "ambari_services_stopservice", OpDescription: "Stop a specific service on the cluster",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceOperate}, Dangerous: true, Client: c, Logger: l,
	}}
}

func (o *StopService) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
			"context":     map[string]interface{}{"type": "string", "description": "Context message", "default": "Stop service via MCP"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *StopService) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *StopService) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	ctxMsg := "Stop service via MCP"
	if c, ok := args["context"].(string); ok {
		ctxMsg = c
	}
	body := map[string]interface{}{
		"RequestInfo": map[string]interface{}{
			"context": ctxMsg,
			"operation_level": map[string]interface{}{
				"level": "SERVICE", "cluster_name": cluster, "service_name": service,
			},
		},
		"Body": map[string]interface{}{
			"ServiceInfo": map[string]interface{}{"state": "INSTALLED"},
		},
	}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, service), nil, body)
}

// ---------- RestartService ----------

type RestartService struct {
	ops.ActionableBase
}

func NewRestartService(c client.AmbariClient, l *logrus.Logger) *RestartService {
	return &RestartService{ops.ActionableBase{
		OpName: "ambari_services_restartservice", OpDescription: "Restart a specific service",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceRestart}, Dangerous: true, Client: c, Logger: l,
	}}
}

func (o *RestartService) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
			"context":     map[string]interface{}{"type": "string", "description": "Context message"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *RestartService) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *RestartService) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	ctxMsg := "Restart service via MCP"
	if c, ok := args["context"].(string); ok {
		ctxMsg = c
	}
	body := map[string]interface{}{
		"RequestInfo": map[string]interface{}{
			"context": ctxMsg, "command": "RESTART",
			"operation_level": map[string]interface{}{
				"level": "SERVICE", "cluster_name": cluster, "service_name": service,
			},
		},
		"Body": map[string]interface{}{
			"ServiceInfo": map[string]interface{}{"state": "STARTED"},
		},
	}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, service), nil, body)
}

// ---------- EnableMaintenanceMode ----------

type EnableMaintenanceMode struct {
	ops.ActionableBase
}

func NewEnableMaintenanceMode(c client.AmbariClient, l *logrus.Logger) *EnableMaintenanceMode {
	return &EnableMaintenanceMode{ops.ActionableBase{
		OpName: "ambari_services_enablemaintenancemode", OpDescription: "Enable maintenance mode for a service",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceOperate}, Dangerous: false, Client: c, Logger: l,
	}}
}

func (o *EnableMaintenanceMode) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *EnableMaintenanceMode) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *EnableMaintenanceMode) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	body := map[string]interface{}{
		"RequestInfo": map[string]interface{}{"context": "Enable Maintenance Mode via MCP"},
		"Body":        map[string]interface{}{"ServiceInfo": map[string]interface{}{"maintenance_state": "ON"}},
	}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, service), nil, body)
}

// ---------- RunServiceCheck ----------

type RunServiceCheck struct {
	ops.ActionableBase
}

func NewRunServiceCheck(c client.AmbariClient, l *logrus.Logger) *RunServiceCheck {
	return &RunServiceCheck{ops.ActionableBase{
		OpName: "ambari_services_runservicecheck", OpDescription: "Run service check for a service",
		OpCategory: "services", Permissions: []auth.Permission{auth.ServiceOperate}, Dangerous: false, Client: c, Logger: l,
	}}
}

func (o *RunServiceCheck) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name: o.OpName, Description: o.OpDescription,
		InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{
			"clusterName": map[string]interface{}{"type": "string", "description": "Cluster name"},
			"serviceName": map[string]interface{}{"type": "string", "description": "Service name"},
		}, Required: []string{"clusterName", "serviceName"}},
	}
}

func (o *RunServiceCheck) Validate(args map[string]interface{}) error {
	if _, ok := args["clusterName"].(string); !ok {
		return fmt.Errorf("clusterName is required")
	}
	if _, ok := args["serviceName"].(string); !ok {
		return fmt.Errorf("serviceName is required")
	}
	return nil
}

func (o *RunServiceCheck) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cluster, service := args["clusterName"].(string), args["serviceName"].(string)
	body := map[string]interface{}{
		"RequestInfo": map[string]interface{}{
			"command": service + "_SERVICE_CHECK",
			"context": service + " Service Check",
			"operation_level": map[string]interface{}{
				"level": "CLUSTER", "cluster_name": cluster,
			},
		},
		"Requests/resource_filters": []map[string]interface{}{
			{"service_name": service},
		},
	}
	return o.Client.Post(ctx, fmt.Sprintf("/clusters/%s/requests", cluster), nil, body)
}
