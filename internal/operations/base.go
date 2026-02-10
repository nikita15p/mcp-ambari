// Package operations defines the Template Method pattern base interfaces
// for separating read-only and actionable (state-changing) Ambari operations.
//
// Design Patterns Used:
//   - Template Method: BaseOperation defines the execution skeleton;
//     ReadOnlyOperation and ActionableOperation provide hooks.
//   - Strategy: Each concrete operation implements the Operation interface.
//   - Factory/Registry: OperationRegistry auto-registers and resolves operations.
package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mcp-ambari/internal/auth"
	"mcp-ambari/internal/client"
	"github.com/sirupsen/logrus"
)

// OperationType distinguishes read-only from actionable operations
type OperationType string

const (
	ReadOnly   OperationType = "readonly"
	Actionable OperationType = "actionable"
)

// ToolSchema describes an MCP tool's input schema (JSON Schema subset)
type ToolSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required"`
}

// ToolDefinition describes a tool for MCP registration
type ToolDefinition struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema ToolSchema `json:"inputSchema"`
}

// OperationResult wraps the result of executing an operation
type OperationResult struct {
	Tool          string      `json:"tool"`
	OperationType string      `json:"operation_type"`
	ExecutionMs   int64       `json:"execution_ms"`
	Timestamp     string      `json:"timestamp"`
	Result        interface{} `json:"result"`
}

// ---------- Operation Interface (Strategy pattern) ----------

// Operation is the core interface every Ambari operation must implement.
type Operation interface {
	// Metadata
	Name() string
	Description() string
	Type() OperationType
	Category() string

	// Schema for MCP tool registration
	Definition() ToolDefinition

	// Permissions required to execute this operation
	RequiredPermissions() []auth.Permission

	// Validate input arguments before execution
	Validate(args map[string]interface{}) error

	// Execute the operation (called by the template method)
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// ---------- Template Method executor ----------

// Executor runs operations through a standard lifecycle:
//
//	authenticate → authorise → validate → execute → audit
type Executor struct {
	client client.AmbariClient
	logger *logrus.Logger
}

// NewExecutor creates a new operation executor
func NewExecutor(c client.AmbariClient, logger *logrus.Logger) *Executor {
	return &Executor{client: c, logger: logger}
}

// Run applies the Template Method: auth-check → validate → execute → wrap result
func (e *Executor) Run(ctx context.Context, op Operation, args map[string]interface{}, authCtx *auth.AuthContext) (*OperationResult, error) {
	start := time.Now()

	// Step 1: Authorization check
	if err := e.checkPermissions(op, authCtx); err != nil {
		return nil, err
	}

	// Step 2: Extra safety for actionable operations
	if op.Type() == Actionable {
		e.logger.WithFields(logrus.Fields{
			"user": authCtx.Username, "tool": op.Name(), "type": "actionable",
		}).Info("Actionable operation requested")
	}

	// Step 3: Validate arguments
	if err := op.Validate(args); err != nil {
		return nil, fmt.Errorf("validation failed for %s: %w", op.Name(), err)
	}

	// Step 4: Execute
	result, err := op.Execute(ctx, args)
	if err != nil {
		e.logger.WithFields(logrus.Fields{"tool": op.Name(), "error": err}).Error("Operation failed")
		return nil, fmt.Errorf("operation %s failed: %w", op.Name(), err)
	}

	// Step 5: Wrap result with metadata
	elapsed := time.Since(start).Milliseconds()
	return &OperationResult{
		Tool:          op.Name(),
		OperationType: string(op.Type()),
		ExecutionMs:   elapsed,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Result:        result,
	}, nil
}

func (e *Executor) checkPermissions(op Operation, authCtx *auth.AuthContext) error {
	required := op.RequiredPermissions()
	if len(required) == 0 {
		return nil
	}
	if !authCtx.HasAllPermissions(required...) {
		return fmt.Errorf("insufficient permissions for %s (requires %v)", op.Name(), required)
	}
	return nil
}

// ResultJSON is a helper to marshal OperationResult to JSON string
func (r *OperationResult) JSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

// ---------- ReadOnlyBase provides common logic for read-only operations ----------

// ReadOnlyBase is embedded by all read-only operations
type ReadOnlyBase struct {
	OpName        string
	OpDescription string
	OpCategory    string
	Permissions   []auth.Permission
	Client        client.AmbariClient
	Logger        *logrus.Logger
}

func (b *ReadOnlyBase) Name() string                           { return b.OpName }
func (b *ReadOnlyBase) Description() string                    { return b.OpDescription }
func (b *ReadOnlyBase) Type() OperationType                    { return ReadOnly }
func (b *ReadOnlyBase) Category() string                       { return b.OpCategory }
func (b *ReadOnlyBase) RequiredPermissions() []auth.Permission { return b.Permissions }

// ---------- ActionableBase provides common logic for state-changing operations ----------

// ActionableBase is embedded by all actionable (write/mutate) operations
type ActionableBase struct {
	OpName        string
	OpDescription string
	OpCategory    string
	Permissions   []auth.Permission
	Dangerous     bool // true for stop/delete style operations
	Client        client.AmbariClient
	Logger        *logrus.Logger
}

func (b *ActionableBase) Name() string                           { return b.OpName }
func (b *ActionableBase) Description() string                    { return b.OpDescription }
func (b *ActionableBase) Type() OperationType                    { return Actionable }
func (b *ActionableBase) Category() string                       { return b.OpCategory }
func (b *ActionableBase) RequiredPermissions() []auth.Permission { return b.Permissions }

// IsDangerous returns true if the operation can cause data loss or downtime
func (b *ActionableBase) IsDangerous() bool { return b.Dangerous }
