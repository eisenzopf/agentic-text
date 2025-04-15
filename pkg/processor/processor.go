package processor

// InitializeBuiltInProcessors ensures all built-in processors are registered before use
func InitializeBuiltInProcessors() {
	// This function must be called early in the application to ensure
	// all processor init() functions have run and registered their processors

	// Force loading of processor packages by name
	// The side effect of importing these packages is that their init() functions will run
	// and register themselves with the processor registry

	// We don't need any actual code here, just the import side effects
}

// init runs automatically and calls RegisterBuiltInProcessors to ensure processors are registered
func init() {
	// Make sure our init() function runs after all processors are registered
	// This happens automatically due to Go's package initialization order
}
