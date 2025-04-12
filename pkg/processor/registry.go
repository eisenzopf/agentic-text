package processor

import (
	"fmt"
	"sync"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// FactoryFunc is a function that creates processors
type FactoryFunc func(provider llm.Provider, options Options) (Processor, error)

// Global processor registry for storing all registered processor factories
var (
	globalRegistry     = make(map[string]FactoryFunc)
	globalRegistryLock sync.RWMutex
)

// Register registers a processor factory with the registry
func Register(name string, factory FactoryFunc) {
	globalRegistryLock.Lock()
	defer globalRegistryLock.Unlock()
	globalRegistry[name] = factory
}

// Create creates a processor by name
func Create(name string, provider llm.Provider, options Options) (Processor, error) {
	globalRegistryLock.RLock()
	factory, ok := globalRegistry[name]
	globalRegistryLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("processor not found: %s", name)
	}
	return factory(provider, options)
}

// ListProcessors returns a list of registered processor names
func ListProcessors() []string {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()
	names := make([]string, 0, len(globalRegistry))
	for name := range globalRegistry {
		names = append(names, name)
	}
	return names
}
