// Package auth provides authentication and authorization for the Ambari MCP Server
package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// Permission represents a specific authorization capability
type Permission string

const (
	// Cluster permissions
	ClusterView  Permission = "cluster:view"
	ClusterAdmin Permission = "cluster:admin"

	// Service permissions
	ServiceView    Permission = "service:view"
	ServiceOperate Permission = "service:operate"
	ServiceRestart Permission = "service:restart"
	ServiceAdmin   Permission = "service:admin"

	// Host permissions
	HostView   Permission = "host:view"
	HostManage Permission = "host:manage"

	// Alert permissions
	AlertView   Permission = "alert:view"
	AlertManage Permission = "alert:manage"
	AlertAdmin  Permission = "alert:admin"

	// Config permissions
	ConfigView   Permission = "config:view"
	ConfigModify Permission = "config:modify"
)

// PermissionGroups maps group names to their permissions
var PermissionGroups = map[string][]Permission{
	"ADMIN": {
		ClusterView, ClusterAdmin, ServiceView, ServiceOperate, ServiceRestart, ServiceAdmin,
		HostView, HostManage, AlertView, AlertManage, AlertAdmin, ConfigView, ConfigModify,
	},
	"OPERATOR": {
		ClusterView, ServiceView, ServiceOperate, ServiceRestart,
		HostView, AlertView, AlertManage, ConfigView,
	},
	"VIEWER": {
		ClusterView, ServiceView, HostView, AlertView, ConfigView,
	},
}

// AuthContext holds authentication and authorization information
type AuthContext struct {
	Username    string       `json:"username"`
	Groups      []string     `json:"groups"`
	Permissions []Permission `json:"permissions"`
	IsValidated bool         `json:"is_validated"`
	Source      string       `json:"source"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// HasPermission checks if the user has a specific permission
func (a *AuthContext) HasPermission(perm Permission) bool {
	for _, p := range a.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if the user has all required permissions
func (a *AuthContext) HasAllPermissions(perms ...Permission) bool {
	for _, perm := range perms {
		if !a.HasPermission(perm) {
			return false
		}
	}
	return true
}

// AuthProvider is the Strategy interface for different authentication methods
type AuthProvider interface {
	Authenticate(ctx context.Context, headers map[string]string) (*AuthContext, error)
	Name() string
}

// contextKey is used for storing auth context in request context
type contextKey string

const authContextKey contextKey = "auth_context"

// WithAuthContext stores the auth context in the request context
func WithAuthContext(ctx context.Context, authCtx *AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, authCtx)
}

// GetAuthContext retrieves the auth context from the request context
func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	authCtx, ok := ctx.Value(authContextKey).(*AuthContext)
	return authCtx, ok
}

// LDAPProvider implements AuthProvider for LDAP authentication via headers
type LDAPProvider struct {
	headerPrefix       string
	groupMappings      map[string][]string
	defaultPermissions []Permission
	logger             *logrus.Logger
}

// NewLDAPProvider creates a new LDAP authentication provider
func NewLDAPProvider(headerPrefix string, groupMappings map[string][]string, defaultPerms []string, logger *logrus.Logger) *LDAPProvider {
	// Convert string permissions to Permission type
	perms := make([]Permission, len(defaultPerms))
	for i, p := range defaultPerms {
		perms[i] = Permission(p)
	}
	
	return &LDAPProvider{
		headerPrefix:       headerPrefix,
		groupMappings:      groupMappings,
		defaultPermissions: perms,
		logger:             logger,
	}
}

func (p *LDAPProvider) Name() string {
	return "LDAP"
}

func (p *LDAPProvider) Authenticate(ctx context.Context, headers map[string]string) (*AuthContext, error) {
	// Extract username from headers
	username := headers[p.headerPrefix+"name"]
	if username == "" {
		username = headers[p.headerPrefix+"username"]
	}
	if username == "" {
		return nil, fmt.Errorf("username not found in headers")
	}

	// Extract groups from headers
	groupsHeader := headers[p.headerPrefix+"groups"]
	var groups []string
	if groupsHeader != "" {
		groups = strings.Split(groupsHeader, ",")
		for i, group := range groups {
			groups[i] = strings.TrimSpace(group)
		}
	}

	// Map groups to permissions
	permSet := make(map[Permission]bool)
	for _, group := range groups {
		if mappedPerms, exists := p.groupMappings[group]; exists {
			for _, perm := range mappedPerms {
				permSet[Permission(perm)] = true
			}
		}
	}

	// Add default permissions if no group mappings found
	if len(permSet) == 0 {
		for _, perm := range p.defaultPermissions {
			permSet[perm] = true
		}
	}

	// Convert permission set to slice
	permissions := make([]Permission, 0, len(permSet))
	for perm := range permSet {
		permissions = append(permissions, perm)
	}

	return &AuthContext{
		Username:    username,
		Groups:      groups,
		Permissions: permissions,
		IsValidated: true,
		Source:      "LDAP",
		Headers:     headers,
	}, nil
}

// Middleware provides HTTP middleware for authentication
type Middleware struct {
	provider AuthProvider
	enabled  bool
	logger   *logrus.Logger
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(provider AuthProvider, enabled bool, logger *logrus.Logger) *Middleware {
	return &Middleware{
		provider: provider,
		enabled:  enabled,
		logger:   logger,
	}
}

// Handler wraps an HTTP handler with authentication middleware
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if m.enabled {
			// Extract headers for authentication
			headers := make(map[string]string)
			for name, values := range r.Header {
				if len(values) > 0 {
					headers[strings.ToLower(name)] = values[0]
				}
			}

			// Authenticate the request
			authCtx, err := m.provider.Authenticate(ctx, headers)
			if err != nil {
				m.logger.WithError(err).Warn("Authentication failed")
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
				return
			}

			// Add auth context to request context
			ctx = WithAuthContext(ctx, authCtx)
			r = r.WithContext(ctx)

			m.logger.WithFields(logrus.Fields{
				"user":   authCtx.Username,
				"groups": authCtx.Groups,
				"source": authCtx.Source,
			}).Debug("Request authenticated")
		} else {
			// Create default auth context for disabled auth
			defaultCtx := &AuthContext{
				Username:    "default-user",
				Groups:      []string{"ambari-admins"},
				Permissions: PermissionGroups["ADMIN"],
				IsValidated: false,
				Source:      "disabled",
			}
			ctx = WithAuthContext(ctx, defaultCtx)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}