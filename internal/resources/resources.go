// Package resources defines MCP Resources for the Ambari MCP Server.
// Resources are always READ-ONLY data endpoints accessed by URI pattern.
// They complement Tools by providing structured, browsable cluster data.
package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"mcp-ambari/internal/client"
	"github.com/sirupsen/logrus"
)

// ResourceDefinition describes a single MCP resource
type ResourceDefinition struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// ResourceResult wraps a resource read result
type ResourceResult struct {
	URI       string      `json:"uri"`
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Handler is a function that resolves a resource URI to data
type Handler func(ctx context.Context, params map[string]string) (*ResourceResult, error)

// Registry holds all resource definitions and their handlers
type Registry struct {
	definitions []ResourceDefinition
	handlers    map[string]Handler
	client      client.AmbariClient
	logger      *logrus.Logger
}

// NewRegistry creates a resource registry with all Ambari resources
func NewRegistry(c client.AmbariClient, logger *logrus.Logger) *Registry {
	r := &Registry{
		handlers: make(map[string]Handler),
		client:   c,
		logger:   logger,
	}
	r.registerAll()
	return r
}

// Definitions returns all resource definitions for MCP ListResources
func (r *Registry) Definitions() []ResourceDefinition {
	return r.definitions
}

// Read resolves a resource URI and returns the data
func (r *Registry) Read(ctx context.Context, uri string) (*ResourceResult, error) {
	resType, params, err := r.parseURI(uri)
	if err != nil {
		return nil, err
	}
	handler, ok := r.handlers[resType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resType)
	}
	return handler(ctx, params)
}

// Count returns number of registered resources
func (r *Registry) Count() int {
	return len(r.definitions)
}

func (r *Registry) registerAll() {
	// 1. Clusters list
	r.add(ResourceDefinition{
		URI: "ambari://clusters", Name: "Ambari Clusters",
		Description: "List of all Ambari clusters with basic information", MimeType: "application/json",
	}, "clusters", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, "/clusters", map[string]string{
			"fields": "Clusters/cluster_name,Clusters/version,Clusters/state",
		})
		return r.wrap("ambari://clusters", "clusters", data), err
	})

	// 2. Cluster details
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}", Name: "Cluster Details",
		Description: "Detailed information about a specific cluster", MimeType: "application/json",
	}, "cluster", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s", params["clusterName"]), map[string]string{
			"fields": "Clusters/*,services/ServiceInfo/service_name,services/ServiceInfo/state,hosts/Hosts/host_name,hosts/Hosts/host_status",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"], "cluster-details", data), err
	})

	// 3. Cluster services
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/services", Name: "Cluster Services",
		Description: "All services running in a cluster with their status", MimeType: "application/json",
	}, "services", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/services", params["clusterName"]), map[string]string{
			"fields": "ServiceInfo/service_name,ServiceInfo/state,ServiceInfo/maintenance_state",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/services", "cluster-services", data), err
	})

	// 4. Cluster hosts
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/hosts", Name: "Cluster Hosts",
		Description: "All hosts in a cluster with status and components", MimeType: "application/json",
	}, "hosts", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/hosts", params["clusterName"]), map[string]string{
			"fields": "Hosts/host_name,Hosts/host_status,Hosts/maintenance_state",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/hosts", "cluster-hosts", data), err
	})

	// 5. Cluster alerts
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/alerts", Name: "Cluster Alerts",
		Description: "Current alerts for a cluster grouped by severity", MimeType: "application/json",
	}, "alerts", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/alerts", params["clusterName"]), map[string]string{
			"fields": "Alert/definition_name,Alert/service_name,Alert/host_name,Alert/state,Alert/text",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/alerts", "cluster-alerts", data), err
	})

	// 6. Alert summary
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/alerts/summary", Name: "Alert Summary",
		Description: "Summarized alert information for quick health overview", MimeType: "application/json",
	}, "alerts-summary", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/alerts", params["clusterName"]), map[string]string{
			"format": "groupedSummary",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/alerts/summary", "alerts-summary", data), err
	})

	// 7. Stale configurations
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/services/stale-configs", Name: "Stale Configurations",
		Description: "Services needing restart due to configuration changes", MimeType: "application/json",
	}, "stale-configs", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/host_components", params["clusterName"]), map[string]string{
			"fields":                  "HostRoles/component_name,HostRoles/host_name,HostRoles/service_name,HostRoles/state,HostRoles/stale_configs",
			"HostRoles/stale_configs": "true",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/services/stale-configs", "stale-configs", data), err
	})

	// 8. Service details
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/service/{serviceName}", Name: "Service Details",
		Description: "Detailed information about a specific service", MimeType: "application/json",
	}, "service", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/services/%s", params["clusterName"], params["serviceName"]), map[string]string{
			"fields": "ServiceInfo/*,components/ServiceComponentInfo/*,components/host_components/HostRoles/state",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/service/"+params["serviceName"], "service-details", data), err
	})

	// 9. Service components
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/service/{serviceName}/components", Name: "Service Components",
		Description: "All components of a service with host assignments", MimeType: "application/json",
	}, "service-components", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/services/%s", params["clusterName"], params["serviceName"]), map[string]string{
			"fields": "components/ServiceComponentInfo/component_name,components/ServiceComponentInfo/category,components/host_components/HostRoles/host_name,components/host_components/HostRoles/state",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/service/"+params["serviceName"]+"/components", "service-components", data), err
	})

	// 10. Host details
	r.add(ResourceDefinition{
		URI: "ambari://host/{hostName}", Name: "Host Details",
		Description: "Detailed information about a specific host", MimeType: "application/json",
	}, "host", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/hosts/%s", params["hostName"]), map[string]string{
			"fields": "Hosts/*,host_components/HostRoles/component_name,host_components/HostRoles/state",
		})
		return r.wrap("ambari://host/"+params["hostName"], "host-details", data), err
	})

	// 11. Recent operations
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/requests/recent", Name: "Recent Operations",
		Description: "Recent operations and their status", MimeType: "application/json",
	}, "recent-requests", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/requests", params["clusterName"]), map[string]string{
			"fields": "Requests/id,Requests/request_context,Requests/request_status,Requests/progress_percent",
			"sortBy": "Requests/id.desc", "page_size": "20",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/requests/recent", "recent-requests", data), err
	})

	// 12. Cluster configurations
	r.add(ResourceDefinition{
		URI: "ambari://cluster/{clusterName}/configurations", Name: "Cluster Configurations",
		Description: "Current configuration types for all services", MimeType: "application/json",
	}, "configurations", func(ctx context.Context, params map[string]string) (*ResourceResult, error) {
		data, err := r.client.Get(ctx, fmt.Sprintf("/clusters/%s/configurations", params["clusterName"]), map[string]string{
			"fields": "Config/type,Config/tag,Config/version",
		})
		return r.wrap("ambari://cluster/"+params["clusterName"]+"/configurations", "configurations", data), err
	})

	r.logger.WithField("count", len(r.definitions)).Info("MCP resources registered")
}

func (r *Registry) add(def ResourceDefinition, resType string, handler Handler) {
	r.definitions = append(r.definitions, def)
	r.handlers[resType] = handler
}

func (r *Registry) wrap(uri, resType string, data interface{}) *ResourceResult {
	return &ResourceResult{
		URI: uri, Type: resType,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
}

// parseURI parses an ambari:// URI into resource type and parameters
func (r *Registry) parseURI(uri string) (string, map[string]string, error) {
	if !strings.HasPrefix(uri, "ambari://") {
		return "", nil, fmt.Errorf("invalid URI: %s", uri)
	}
	path := strings.TrimPrefix(uri, "ambari://")
	params := map[string]string{}

	if path == "clusters" {
		return "clusters", params, nil
	}

	if strings.HasPrefix(path, "host/") {
		params["hostName"] = strings.TrimPrefix(path, "host/")
		return "host", params, nil
	}

	if strings.HasPrefix(path, "cluster/") {
		parts := strings.SplitN(strings.TrimPrefix(path, "cluster/"), "/", 2)
		params["clusterName"] = parts[0]
		if len(parts) == 1 {
			return "cluster", params, nil
		}
		sub := parts[1]
		switch {
		case sub == "services":
			return "services", params, nil
		case sub == "hosts":
			return "hosts", params, nil
		case sub == "alerts":
			return "alerts", params, nil
		case sub == "alerts/summary":
			return "alerts-summary", params, nil
		case sub == "services/stale-configs":
			return "stale-configs", params, nil
		case sub == "requests/recent":
			return "recent-requests", params, nil
		case sub == "configurations":
			return "configurations", params, nil
		case strings.HasPrefix(sub, "service/"):
			svcParts := strings.SplitN(strings.TrimPrefix(sub, "service/"), "/", 2)
			params["serviceName"] = svcParts[0]
			if len(svcParts) == 1 {
				return "service", params, nil
			}
			if svcParts[1] == "components" {
				return "service-components", params, nil
			}
		}
	}

	return "", nil, fmt.Errorf("unsupported resource URI: %s", uri)
}

// ToJSON converts a ResourceResult to JSON string
func (r *ResourceResult) ToJSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}
