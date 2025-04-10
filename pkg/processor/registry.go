package processor

import (
	"fmt"
	"sync"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

var (
	registry = make(map[string]ProcessorFactory)
	mutex    sync.RWMutex
)

// ProcessorFactory is a function that creates a Processor
type ProcessorFactory func(provider llm.Provider, options Options) (Processor, error)

// Register registers a processor factory with the registry
func Register(name string, factory ProcessorFactory) {
	mutex.Lock()
	defer mutex.Unlock()
	registry[name] = factory
}

// GetProcessor retrieves a processor by name
func GetProcessor(name string, provider llm.Provider, options Options) (Processor, error) {
	mutex.RLock()
	factory, exists := registry[name]
	mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("processor not found: %s", name)
	}

	return factory(provider, options)
}

// ListProcessors returns a list of all registered processor names
func ListProcessors() []string {
	mutex.RLock()
	defer mutex.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}

	return names
}

// init registers built-in processors
func init() {
	// Register built-in processors here
	// This will be populated as we implement individual processors
}
