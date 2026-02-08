/* START GENAI */
// Package readonly provides read-only operations for user and group management
package readonly

import (
	"context"
	"fmt"

	"github.com/nikita15p/mcp-ambari/internal/auth"
	"github.com/nikita15p/mcp-ambari/internal/client"
	ops "github.com/nikita15p/mcp-ambari/internal/operations"
	"github.com/sirupsen/logrus"
)

// ---- GetUsers ----

type GetUsers struct{ ops.ReadOnlyBase }

func NewGetUsers(c client.AmbariClient, l *logrus.Logger) *GetUsers {
	return &GetUsers{ops.ReadOnlyBase{
		OpName:        "ambari_users_getusers",
		OpDescription: "Get all users in Ambari",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterView},
		Client:        c, Logger: l,
	}}
}

func (o *GetUsers) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"fields": map[string]interface{}{
					"type":        "string",
					"description": "Comma-separated fields to return (optional)",
					"default":     "Users/*",
				},
			},
			Required: []string{},
		},
	}
}

func (o *GetUsers) Validate(args map[string]interface{}) error {
	return nil // No required parameters
}

func (o *GetUsers) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := map[string]string{
		"fields": "Users/*",
	}

	if fields, ok := args["fields"].(string); ok && fields != "" {
		params["fields"] = fields
	}

	return o.Client.Get(ctx, "/users", params)
}

// ---- GetUser ----

type GetUser struct{ ops.ReadOnlyBase }

func NewGetUser(c client.AmbariClient, l *logrus.Logger) *GetUser {
	return &GetUser{ops.ReadOnlyBase{
		OpName:        "ambari_users_getuser",
		OpDescription: "Get details of a specific user",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterView},
		Client:        c, Logger: l,
	}}
}

func (o *GetUser) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username": map[string]interface{}{
					"type":        "string",
					"description": "Username to get details for",
				},
			},
			Required: []string{"username"},
		},
	}
}

func (o *GetUser) Validate(args map[string]interface{}) error {
	if _, ok := args["username"].(string); !ok {
		return fmt.Errorf("username required")
	}
	return nil
}

func (o *GetUser) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	username := args["username"].(string)
	return o.Client.Get(ctx, fmt.Sprintf("/users/%s", username), map[string]string{
		"fields": "Users/*",
	})
}

// ---- GetGroups ----

type GetGroups struct{ ops.ReadOnlyBase }

func NewGetGroups(c client.AmbariClient, l *logrus.Logger) *GetGroups {
	return &GetGroups{ops.ReadOnlyBase{
		OpName:        "ambari_users_getgroups",
		OpDescription: "Get all groups in Ambari",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterView},
		Client:        c, Logger: l,
	}}
}

func (o *GetGroups) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"fields": map[string]interface{}{
					"type":        "string",
					"description": "Comma-separated fields to return (optional)",
					"default":     "Groups/*",
				},
			},
			Required: []string{},
		},
	}
}

func (o *GetGroups) Validate(args map[string]interface{}) error {
	return nil // No required parameters
}

func (o *GetGroups) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := map[string]string{
		"fields": "Groups/*",
	}

	if fields, ok := args["fields"].(string); ok && fields != "" {
		params["fields"] = fields
	}

	return o.Client.Get(ctx, "/groups", params)
}

// ---- GetGroup ----

type GetGroup struct{ ops.ReadOnlyBase }

func NewGetGroup(c client.AmbariClient, l *logrus.Logger) *GetGroup {
	return &GetGroup{ops.ReadOnlyBase{
		OpName:        "ambari_users_getgroup",
		OpDescription: "Get details of a specific group",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterView},
		Client:        c, Logger: l,
	}}
}

func (o *GetGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"groupName": map[string]interface{}{
					"type":        "string",
					"description": "Group name to get details for",
				},
			},
			Required: []string{"groupName"},
		},
	}
}

func (o *GetGroup) Validate(args map[string]interface{}) error {
	if _, ok := args["groupName"].(string); !ok {
		return fmt.Errorf("groupName required")
	}
	return nil
}

func (o *GetGroup) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	groupName := args["groupName"].(string)
	return o.Client.Get(ctx, fmt.Sprintf("/groups/%s", groupName), map[string]string{
		"fields": "Groups/*",
	})
}

// ---- GetUserPrivileges ----

type GetUserPrivileges struct{ ops.ReadOnlyBase }

func NewGetUserPrivileges(c client.AmbariClient, l *logrus.Logger) *GetUserPrivileges {
	return &GetUserPrivileges{ops.ReadOnlyBase{
		OpName:        "ambari_users_getuserprivileges",
		OpDescription: "Get privileges assigned to a specific user",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterView},
		Client:        c, Logger: l,
	}}
}

func (o *GetUserPrivileges) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username": map[string]interface{}{
					"type":        "string",
					"description": "Username to get privileges for",
				},
			},
			Required: []string{"username"},
		},
	}
}

func (o *GetUserPrivileges) Validate(args map[string]interface{}) error {
	if _, ok := args["username"].(string); !ok {
		return fmt.Errorf("username required")
	}
	return nil
}

func (o *GetUserPrivileges) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	username := args["username"].(string)
	return o.Client.Get(ctx, fmt.Sprintf("/users/%s/privileges", username), map[string]string{
		"fields": "PrivilegeInfo/*",
	})
}

/* END GENAI */
