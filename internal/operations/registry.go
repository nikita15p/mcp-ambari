package operations

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Registry holds all registered operations and provides lookup by name and type.
// Implements the Registry pattern with Factory-style creation helpers.
type Registry struct {
	mu         sync.RWMutex
	ops        map[string]Operation
	readonly   []Operation
	actionable []Operation
	logger     *logrus.Logger
}

// NewRegistry creates a new operation registry
func NewRegistry(logger *logrus.Logger) *Registry {
	return &Registry{
		ops:    make(map[string]Operation),
		logger: logger,
	}
}

// Register adds an operation to the registry
func (r *Registry) Register(op Operation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.ops[op.Name()]; exists {
		return fmt.Errorf("operation %s already registered", op.Name())
	}
	r.ops[op.Name()] = op

	switch op.Type() {
	case ReadOnly:
		r.readonly = append(r.readonly, op)
	case Actionable:
		r.actionable = append(r.actionable, op)
	}

	r.logger.WithFields(logrus.Fields{
		"name": op.Name(), "type": op.Type(), "category": op.Category(),
	}).Debug("Operation registered")
	return nil
}

// Get returns an operation by name
func (r *Registry) Get(name string) (Operation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.ops[name]
	return op, ok
}

// All returns all registered operations
func (r *Registry) All() []Operation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Operation, 0, len(r.ops))
	for _, op := range r.ops {
		result = append(result, op)
	}
	return result
}

// ReadOnlyOps returns only read-only operations
func (r *Registry) ReadOnlyOps() []Operation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Operation{}, r.readonly...)
}

// ActionableOps returns only actionable (state-changing) operations
func (r *Registry) ActionableOps() []Operation {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Operation{}, r.actionable...)
}

// Count returns the total, read-only, and actionable counts
func (r *Registry) Count() (total, ro, act int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.ops), len(r.readonly), len(r.actionable)
}

// Definitions returns tool definitions for all operations (for MCP ListTools)
func (r *Registry) Definitions() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	defs := make([]ToolDefinition, 0, len(r.ops))
	for _, op := range r.ops {
		defs = append(defs, op.Definition())
	}
	return defs
}
