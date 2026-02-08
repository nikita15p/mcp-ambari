package actionable

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/niita15p/mcp-ambari/internal/auth"
	"github.com/niita15p/mcp-ambari/internal/client"
	ops "github.com/niita15p/mcp-ambari/internal/operations"
	"github.com/sirupsen/logrus"
)

// ---- CreateCluster ----
type CreateCluster struct{ ops.ActionableBase }

func NewCreateCluster(c client.AmbariClient, l *logrus.Logger) *CreateCluster {
	return &CreateCluster{ops.ActionableBase{OpName: "ambari_clusters_createcluster", OpDescription: "Creates a cluster", OpCategory: "clusters", Permissions: []auth.Permission{auth.ClusterAdmin}, Dangerous: true, Client: c, Logger: l}}
}
func (o *CreateCluster) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "body": m("string", "JSON body for cluster creation")}, Required: []string{"clusterName", "body"}}}
}
func (o *CreateCluster) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "body")
}
func (o *CreateCluster) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	var body interface{}
	if s, ok := a["body"].(string); ok {
		json.Unmarshal([]byte(s), &body)
	} else {
		body = a["body"]
	}
	return o.Client.Post(ctx, fmt.Sprintf("/clusters/%s", a["clusterName"].(string)), nil, body)
}

// ---- UpdateAlertDefinition ----
type UpdateAlertDefinition struct{ ops.ActionableBase }

func NewUpdateAlertDefinition(c client.AmbariClient, l *logrus.Logger) *UpdateAlertDefinition {
	return &UpdateAlertDefinition{ops.ActionableBase{OpName: "ambari_alerts_updatealertdefinition", OpDescription: "Update an alert definition (enable/disable or modify)", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *UpdateAlertDefinition) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "definitionId": m("string", "Alert definition ID"), "enabled": m("boolean", "Enable or disable"), "data": m("string", "JSON of additional properties")}, Required: []string{"clusterName", "definitionId"}}}
}
func (o *UpdateAlertDefinition) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "definitionId")
}
func (o *UpdateAlertDefinition) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	body := map[string]interface{}{}
	if e, ok := a["enabled"].(bool); ok {
		body["AlertDefinition/enabled"] = e
	}
	if d, ok := a["data"].(string); ok {
		var extra map[string]interface{}
		json.Unmarshal([]byte(d), &extra)
		for k, v := range extra {
			body[k] = v
		}
	}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/alert_definitions/%s", a["clusterName"].(string), a["definitionId"].(string)), nil, body)
}

// ---- CreateAlertGroup ----
type CreateAlertGroup struct{ ops.ActionableBase }

func NewCreateAlertGroup(c client.AmbariClient, l *logrus.Logger) *CreateAlertGroup {
	return &CreateAlertGroup{ops.ActionableBase{OpName: "ambari_alerts_createalertgroup", OpDescription: "Create a new alert group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *CreateAlertGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "groupName": m("string", "Alert group name"), "definitions": m("string", "JSON array of definition IDs")}, Required: []string{"clusterName", "groupName"}}}
}
func (o *CreateAlertGroup) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "groupName")
}
func (o *CreateAlertGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	body := map[string]interface{}{"AlertGroup": map[string]interface{}{"name": a["groupName"].(string)}}
	if d, ok := a["definitions"].(string); ok {
		var defs interface{}
		json.Unmarshal([]byte(d), &defs)
		body["AlertGroup"].(map[string]interface{})["definitions"] = defs
	}
	return o.Client.Post(ctx, fmt.Sprintf("/clusters/%s/alert_groups", a["clusterName"].(string)), nil, body)
}

// ---- UpdateAlertGroup ----
type UpdateAlertGroup struct{ ops.ActionableBase }

func NewUpdateAlertGroup(c client.AmbariClient, l *logrus.Logger) *UpdateAlertGroup {
	return &UpdateAlertGroup{ops.ActionableBase{OpName: "ambari_alerts_updatealertgroup", OpDescription: "Update an existing alert group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *UpdateAlertGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "groupId": m("integer", "Alert group ID"), "groupName": m("string", "New group name"), "definitions": m("string", "JSON array of definition IDs")}, Required: []string{"clusterName", "groupId", "groupName"}}}
}
func (o *UpdateAlertGroup) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "groupName")
}
func (o *UpdateAlertGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	gid := fmt.Sprintf("%v", a["groupId"])
	body := map[string]interface{}{"AlertGroup": map[string]interface{}{"name": a["groupName"].(string)}}
	if d, ok := a["definitions"].(string); ok {
		var defs interface{}
		json.Unmarshal([]byte(d), &defs)
		body["AlertGroup"].(map[string]interface{})["definitions"] = defs
	}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%s", a["clusterName"].(string), gid), nil, body)
}

// ---- DeleteAlertGroup ----
type DeleteAlertGroup struct{ ops.ActionableBase }

func NewDeleteAlertGroup(c client.AmbariClient, l *logrus.Logger) *DeleteAlertGroup {
	return &DeleteAlertGroup{ops.ActionableBase{OpName: "ambari_alerts_deletealertgroup", OpDescription: "Delete an alert group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertAdmin}, Dangerous: true, Client: c, Logger: l}}
}
func (o *DeleteAlertGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "groupId": m("integer", "Alert group ID to delete")}, Required: []string{"clusterName", "groupId"}}}
}
func (o *DeleteAlertGroup) Validate(a map[string]interface{}) error { return req(a, "clusterName") }
func (o *DeleteAlertGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	return o.Client.Delete(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%v", a["clusterName"].(string), a["groupId"]), nil)
}

// ---- DuplicateAlertGroup ----
type DuplicateAlertGroup struct{ ops.ActionableBase }

func NewDuplicateAlertGroup(c client.AmbariClient, l *logrus.Logger) *DuplicateAlertGroup {
	return &DuplicateAlertGroup{ops.ActionableBase{OpName: "ambari_alerts_duplicatealertgroup", OpDescription: "Duplicate an alert group with a new name", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *DuplicateAlertGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "sourceGroupId": m("integer", "Source group ID"), "newGroupName": m("string", "New group name")}, Required: []string{"clusterName", "sourceGroupId", "newGroupName"}}}
}
func (o *DuplicateAlertGroup) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "newGroupName")
}
func (o *DuplicateAlertGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	cluster := a["clusterName"].(string)
	src, _ := o.Client.Get(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%v", cluster, a["sourceGroupId"]), map[string]string{"fields": "*"})
	body := map[string]interface{}{"AlertGroup": map[string]interface{}{"name": a["newGroupName"].(string)}}
	if src != nil {
		if ag, ok := src["AlertGroup"].(map[string]interface{}); ok {
			if defs, ok := ag["definitions"]; ok {
				body["AlertGroup"].(map[string]interface{})["definitions"] = defs
			}
		}
	}
	return o.Client.Post(ctx, fmt.Sprintf("/clusters/%s/alert_groups", cluster), nil, body)
}

// ---- AddDefinitionToGroup ----
type AddDefinitionToGroup struct{ ops.ActionableBase }

func NewAddDefinitionToGroup(c client.AmbariClient, l *logrus.Logger) *AddDefinitionToGroup {
	return &AddDefinitionToGroup{ops.ActionableBase{OpName: "ambari_alerts_adddefinitiontogroup", OpDescription: "Add an alert definition to a group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *AddDefinitionToGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster name"), "groupId": m("integer", "Group ID"), "definitionId": m("integer", "Definition ID to add")}, Required: []string{"clusterName", "groupId", "definitionId"}}}
}
func (o *AddDefinitionToGroup) Validate(a map[string]interface{}) error { return req(a, "clusterName") }
func (o *AddDefinitionToGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	return o.Client.Post(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%v/alert_definitions/%v", a["clusterName"].(string), a["groupId"], a["definitionId"]), nil, nil)
}

// ---- RemoveDefinitionFromGroup ----
type RemoveDefinitionFromGroup struct{ ops.ActionableBase }

func NewRemoveDefinitionFromGroup(c client.AmbariClient, l *logrus.Logger) *RemoveDefinitionFromGroup {
	return &RemoveDefinitionFromGroup{ops.ActionableBase{OpName: "ambari_alerts_removedefinitionfromgroup", OpDescription: "Remove an alert definition from a group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *RemoveDefinitionFromGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "groupId": m("integer", "Group ID"), "definitionId": m("integer", "Definition ID")}, Required: []string{"clusterName", "groupId", "definitionId"}}}
}
func (o *RemoveDefinitionFromGroup) Validate(a map[string]interface{}) error {
	return req(a, "clusterName")
}
func (o *RemoveDefinitionFromGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	return o.Client.Delete(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%v/alert_definitions/%v", a["clusterName"].(string), a["groupId"], a["definitionId"]), nil)
}

// ---- CreateNotification ----
type CreateNotification struct{ ops.ActionableBase }

func NewCreateNotification(c client.AmbariClient, l *logrus.Logger) *CreateNotification {
	return &CreateNotification{ops.ActionableBase{OpName: "ambari_alerts_createnotification", OpDescription: "Create a new alert notification target", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *CreateNotification) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "notificationData": m("string", "JSON notification target data")}, Required: []string{"clusterName", "notificationData"}}}
}
func (o *CreateNotification) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "notificationData")
}
func (o *CreateNotification) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	var body interface{}
	json.Unmarshal([]byte(a["notificationData"].(string)), &body)
	return o.Client.Post(ctx, "/alert_targets", nil, body)
}

// ---- UpdateNotification ----
type UpdateNotification struct{ ops.ActionableBase }

func NewUpdateNotification(c client.AmbariClient, l *logrus.Logger) *UpdateNotification {
	return &UpdateNotification{ops.ActionableBase{OpName: "ambari_alerts_updatenotification", OpDescription: "Update an alert notification target", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *UpdateNotification) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "targetId": m("integer", "Target ID"), "notificationData": m("string", "JSON updated data")}, Required: []string{"clusterName", "targetId", "notificationData"}}}
}
func (o *UpdateNotification) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "notificationData")
}
func (o *UpdateNotification) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	var body interface{}
	json.Unmarshal([]byte(a["notificationData"].(string)), &body)
	return o.Client.Put(ctx, fmt.Sprintf("/alert_targets/%v", a["targetId"]), nil, body)
}

// ---- DeleteNotification ----
type DeleteNotification struct{ ops.ActionableBase }

func NewDeleteNotification(c client.AmbariClient, l *logrus.Logger) *DeleteNotification {
	return &DeleteNotification{ops.ActionableBase{OpName: "ambari_alerts_deletenotification", OpDescription: "Delete an alert notification target", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertAdmin}, Dangerous: true, Client: c, Logger: l}}
}
func (o *DeleteNotification) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "targetId": m("integer", "Target ID to delete")}, Required: []string{"clusterName", "targetId"}}}
}
func (o *DeleteNotification) Validate(a map[string]interface{}) error { return req(a, "clusterName") }
func (o *DeleteNotification) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	return o.Client.Delete(ctx, fmt.Sprintf("/alert_targets/%v", a["targetId"]), nil)
}

// ---- AddNotificationToGroup ----
type AddNotificationToGroup struct{ ops.ActionableBase }

func NewAddNotificationToGroup(c client.AmbariClient, l *logrus.Logger) *AddNotificationToGroup {
	return &AddNotificationToGroup{ops.ActionableBase{OpName: "ambari_alerts_addnotificationtogroup", OpDescription: "Add notification target to alert group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *AddNotificationToGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "groupId": m("integer", "Group ID"), "targetId": m("integer", "Target ID")}, Required: []string{"clusterName", "groupId", "targetId"}}}
}
func (o *AddNotificationToGroup) Validate(a map[string]interface{}) error {
	return req(a, "clusterName")
}
func (o *AddNotificationToGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	return o.Client.Post(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%v/alert_targets/%v", a["clusterName"].(string), a["groupId"], a["targetId"]), nil, nil)
}

// ---- RemoveNotificationFromGroup ----
type RemoveNotificationFromGroup struct{ ops.ActionableBase }

func NewRemoveNotificationFromGroup(c client.AmbariClient, l *logrus.Logger) *RemoveNotificationFromGroup {
	return &RemoveNotificationFromGroup{ops.ActionableBase{OpName: "ambari_alerts_removenotificationfromgroup", OpDescription: "Remove notification target from alert group", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertManage}, Dangerous: false, Client: c, Logger: l}}
}
func (o *RemoveNotificationFromGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "groupId": m("integer", "Group ID"), "targetId": m("integer", "Target ID")}, Required: []string{"clusterName", "groupId", "targetId"}}}
}
func (o *RemoveNotificationFromGroup) Validate(a map[string]interface{}) error {
	return req(a, "clusterName")
}
func (o *RemoveNotificationFromGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	return o.Client.Delete(ctx, fmt.Sprintf("/clusters/%s/alert_groups/%v/alert_targets/%v", a["clusterName"].(string), a["groupId"], a["targetId"]), nil)
}

// ---- SaveAlertSettings ----
type SaveAlertSettings struct{ ops.ActionableBase }

func NewSaveAlertSettings(c client.AmbariClient, l *logrus.Logger) *SaveAlertSettings {
	return &SaveAlertSettings{ops.ActionableBase{OpName: "ambari_alerts_savealertsettings", OpDescription: "Save cluster-level alert settings", OpCategory: "alerts", Permissions: []auth.Permission{auth.AlertAdmin}, Dangerous: false, Client: c, Logger: l}}
}
func (o *SaveAlertSettings) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "alertRepeatTolerance": m("integer", "Alert repeat tolerance value")}, Required: []string{"clusterName", "alertRepeatTolerance"}}}
}
func (o *SaveAlertSettings) Validate(a map[string]interface{}) error { return req(a, "clusterName") }
func (o *SaveAlertSettings) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	cluster := a["clusterName"].(string)
	tolerance := fmt.Sprintf("%v", a["alertRepeatTolerance"])
	body := map[string]interface{}{"Clusters": map[string]interface{}{"desired_config": map[string]interface{}{"type": "cluster-env", "properties": map[string]interface{}{"alerts_repeat_tolerance": tolerance}}}}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s", cluster), nil, body)
}

// ---- RestartComponents ----
type RestartComponents struct{ ops.ActionableBase }

func NewRestartComponents(c client.AmbariClient, l *logrus.Logger) *RestartComponents {
	return &RestartComponents{ops.ActionableBase{OpName: "ambari_services_restartcomponents", OpDescription: "Restart specific components with stale configurations", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceRestart}, Dangerous: true, Client: c, Logger: l}}
}
func (o *RestartComponents) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "serviceName": m("string", "Service"), "componentName": m("string", "Component to restart"), "hostNames": m("string", "JSON array of host names"), "context": m("string", "Context message")}, Required: []string{"clusterName", "serviceName", "componentName"}}}
}
func (o *RestartComponents) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "serviceName", "componentName")
}
func (o *RestartComponents) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	cluster, svc, comp := a["clusterName"].(string), a["serviceName"].(string), a["componentName"].(string)
	ctxMsg := "Restart components via MCP"
	if c, ok := a["context"].(string); ok {
		ctxMsg = c
	}
	body := map[string]interface{}{"RequestInfo": map[string]interface{}{"context": ctxMsg, "command": "RESTART", "operation_level": map[string]interface{}{"level": "HOST_COMPONENT", "cluster_name": cluster, "service_name": svc, "hostcomponent_name": comp}}, "Body": map[string]interface{}{"HostRoles": map[string]interface{}{"state": "STARTED"}}}
	path := fmt.Sprintf("/clusters/%s/host_components?HostRoles/component_name=%s&HostRoles/service_name=%s", cluster, comp, svc)
	return o.Client.Put(ctx, path, nil, body)
}

// ---- DisableMaintenanceMode ----
type DisableMaintenanceMode struct{ ops.ActionableBase }

func NewDisableMaintenanceMode(c client.AmbariClient, l *logrus.Logger) *DisableMaintenanceMode {
	return &DisableMaintenanceMode{ops.ActionableBase{OpName: "ambari_services_disablemaintenancemode", OpDescription: "Disable maintenance mode for a service", OpCategory: "services", Permissions: []auth.Permission{auth.ServiceOperate}, Dangerous: false, Client: c, Logger: l}}
}
func (o *DisableMaintenanceMode) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{Name: o.OpName, Description: o.OpDescription, InputSchema: ops.ToolSchema{Type: "object", Properties: map[string]interface{}{"clusterName": m("string", "Cluster"), "serviceName": m("string", "Service"), "componentName": m("string", "Component (optional)"), "hostName": m("string", "Host (required if component)")}, Required: []string{"clusterName", "serviceName"}}}
}
func (o *DisableMaintenanceMode) Validate(a map[string]interface{}) error {
	return req(a, "clusterName", "serviceName")
}
func (o *DisableMaintenanceMode) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	cluster, svc := a["clusterName"].(string), a["serviceName"].(string)
	if comp, ok := a["componentName"].(string); ok && comp != "" {
		host := a["hostName"].(string)
		return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/hosts/%s/host_components/%s", cluster, host, comp), nil, map[string]interface{}{"HostRoles": map[string]interface{}{"maintenance_state": "OFF"}})
	}
	body := map[string]interface{}{"RequestInfo": map[string]interface{}{"context": "Disable Maintenance Mode via MCP"}, "Body": map[string]interface{}{"ServiceInfo": map[string]interface{}{"maintenance_state": "OFF"}}}
	return o.Client.Put(ctx, fmt.Sprintf("/clusters/%s/services/%s", cluster, svc), nil, body)
}

// ---- Helpers ----
func m(t, desc string) map[string]interface{} {
	return map[string]interface{}{"type": t, "description": desc}
}
func req(a map[string]interface{}, keys ...string) error {
	for _, k := range keys {
		if _, ok := a[k].(string); !ok {
			if _, ok2 := a[k]; !ok2 {
				return fmt.Errorf("%s is required", k)
			}
		}
	}
	return nil
}
