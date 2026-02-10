// Package actionable provides actionable operations for user and group management
package actionable

import (
	"context"
	"fmt"

	"github.com/nikita15p/mcp-ambari/internal/auth"
	"github.com/nikita15p/mcp-ambari/internal/client"
	ops "github.com/nikita15p/mcp-ambari/internal/operations"
	"github.com/sirupsen/logrus"
)

// ---- CreateUser ----

type CreateUser struct{ ops.ActionableBase }

func NewCreateUser(c client.AmbariClient, l *logrus.Logger) *CreateUser {
	return &CreateUser{ops.ActionableBase{
		OpName:        "ambari_users_createuser",
		OpDescription: "Create a new Ambari user",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     false,
		Client:        c, Logger: l,
	}}
}

func (o *CreateUser) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username":      m("string", "Username for the new user"),
				"password":      m("string", "Password for the new user"),
				"displayName":   m("string", "Display name (optional)"),
				"localUsername": m("string", "Local username (optional)"),
			},
			Required: []string{"username", "password"},
		},
	}
}

func (o *CreateUser) Validate(a map[string]interface{}) error {
	return req(a, "username", "password")
}

func (o *CreateUser) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	body := map[string]interface{}{
		"Users": map[string]interface{}{
			"user_name": a["username"].(string),
			"password":  a["password"].(string),
		},
	}

	if displayName, ok := a["displayName"].(string); ok {
		body["Users"].(map[string]interface{})["display_name"] = displayName
	}

	if localUsername, ok := a["localUsername"].(string); ok {
		body["Users"].(map[string]interface{})["local_username"] = localUsername
	}

	return o.Client.Post(ctx, "/users", nil, body)
}

// ---- UpdateUser ----

type UpdateUser struct{ ops.ActionableBase }

func NewUpdateUser(c client.AmbariClient, l *logrus.Logger) *UpdateUser {
	return &UpdateUser{ops.ActionableBase{
		OpName:        "ambari_users_updateuser",
		OpDescription: "Update an existing Ambari user",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     false,
		Client:        c, Logger: l,
	}}
}

func (o *UpdateUser) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username":    m("string", "Username to update"),
				"password":    m("string", "New password (optional)"),
				"displayName": m("string", "New display name (optional)"),
				"active":      m("boolean", "User active status (optional)"),
			},
			Required: []string{"username"},
		},
	}
}

func (o *UpdateUser) Validate(a map[string]interface{}) error {
	return req(a, "username")
}

func (o *UpdateUser) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	username := a["username"].(string)
	body := map[string]interface{}{
		"Users": map[string]interface{}{},
	}

	if password, ok := a["password"].(string); ok {
		body["Users"].(map[string]interface{})["password"] = password
	}

	if displayName, ok := a["displayName"].(string); ok {
		body["Users"].(map[string]interface{})["display_name"] = displayName
	}

	if active, ok := a["active"].(bool); ok {
		body["Users"].(map[string]interface{})["active"] = active
	}

	return o.Client.Put(ctx, fmt.Sprintf("/users/%s", username), nil, body)
}

// ---- DeleteUser ----

type DeleteUser struct{ ops.ActionableBase }

func NewDeleteUser(c client.AmbariClient, l *logrus.Logger) *DeleteUser {
	return &DeleteUser{ops.ActionableBase{
		OpName:        "ambari_users_deleteuser",
		OpDescription: "Delete an Ambari user",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     true,
		Client:        c, Logger: l,
	}}
}

func (o *DeleteUser) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username": m("string", "Username to delete"),
			},
			Required: []string{"username"},
		},
	}
}

func (o *DeleteUser) Validate(a map[string]interface{}) error {
	return req(a, "username")
}

func (o *DeleteUser) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	username := a["username"].(string)
	return o.Client.Delete(ctx, fmt.Sprintf("/users/%s", username), nil)
}

// ---- CreateUserGroup ----

type CreateUserGroup struct{ ops.ActionableBase }

func NewCreateUserGroup(c client.AmbariClient, l *logrus.Logger) *CreateUserGroup {
	return &CreateUserGroup{ops.ActionableBase{
		OpName:        "ambari_users_creategroup",
		OpDescription: "Create a new Ambari group",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     false,
		Client:        c, Logger: l,
	}}
}

func (o *CreateUserGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"groupName": m("string", "Group name"),
			},
			Required: []string{"groupName"},
		},
	}
}

func (o *CreateUserGroup) Validate(a map[string]interface{}) error {
	return req(a, "groupName")
}

func (o *CreateUserGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	body := map[string]interface{}{
		"Groups": map[string]interface{}{
			"group_name": a["groupName"].(string),
		},
	}

	return o.Client.Post(ctx, "/groups", nil, body)
}

// ---- DeleteUserGroup ----

type DeleteUserGroup struct{ ops.ActionableBase }

func NewDeleteUserGroup(c client.AmbariClient, l *logrus.Logger) *DeleteUserGroup {
	return &DeleteUserGroup{ops.ActionableBase{
		OpName:        "ambari_users_deletegroup",
		OpDescription: "Delete an Ambari group",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     true,
		Client:        c, Logger: l,
	}}
}

func (o *DeleteUserGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"groupName": m("string", "Group name to delete"),
			},
			Required: []string{"groupName"},
		},
	}
}

func (o *DeleteUserGroup) Validate(a map[string]interface{}) error {
	return req(a, "groupName")
}

func (o *DeleteUserGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	groupName := a["groupName"].(string)
	return o.Client.Delete(ctx, fmt.Sprintf("/groups/%s", groupName), nil)
}

// ---- AddUserToGroup ----

type AddUserToGroup struct{ ops.ActionableBase }

func NewAddUserToGroup(c client.AmbariClient, l *logrus.Logger) *AddUserToGroup {
	return &AddUserToGroup{ops.ActionableBase{
		OpName:        "ambari_users_addusertogroup",
		OpDescription: "Add a user to a group",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     false,
		Client:        c, Logger: l,
	}}
}

func (o *AddUserToGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username":  m("string", "Username to add"),
				"groupName": m("string", "Group name to add user to"),
			},
			Required: []string{"username", "groupName"},
		},
	}
}

func (o *AddUserToGroup) Validate(a map[string]interface{}) error {
	return req(a, "username", "groupName")
}

func (o *AddUserToGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	username := a["username"].(string)
	groupName := a["groupName"].(string)

	return o.Client.Post(ctx, fmt.Sprintf("/groups/%s/members/%s", groupName, username), nil, nil)
}

// ---- RemoveUserFromGroup ----

type RemoveUserFromGroup struct{ ops.ActionableBase }

func NewRemoveUserFromGroup(c client.AmbariClient, l *logrus.Logger) *RemoveUserFromGroup {
	return &RemoveUserFromGroup{ops.ActionableBase{
		OpName:        "ambari_users_removeuserfromgroup",
		OpDescription: "Remove a user from a group",
		OpCategory:    "users",
		Permissions:   []auth.Permission{auth.ClusterAdmin},
		Dangerous:     false,
		Client:        c, Logger: l,
	}}
}

func (o *RemoveUserFromGroup) Definition() ops.ToolDefinition {
	return ops.ToolDefinition{
		Name:        o.OpName,
		Description: o.OpDescription,
		InputSchema: ops.ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"username":  m("string", "Username to remove"),
				"groupName": m("string", "Group name to remove user from"),
			},
			Required: []string{"username", "groupName"},
		},
	}
}

func (o *RemoveUserFromGroup) Validate(a map[string]interface{}) error {
	return req(a, "username", "groupName")
}

func (o *RemoveUserFromGroup) Execute(ctx context.Context, a map[string]interface{}) (interface{}, error) {
	username := a["username"].(string)
	groupName := a["groupName"].(string)

	return o.Client.Delete(ctx, fmt.Sprintf("/groups/%s/members/%s", groupName, username), nil)
}
