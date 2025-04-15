package processor

// NewDefaultOptions creates a new Options instance with default settings
func NewDefaultOptions() Options {
	return Options{
		PreProcessOptions:  make(map[string]interface{}),
		LLMOptions:         make(map[string]interface{}),
		PostProcessOptions: make(map[string]interface{}),
	}
}

// Clone creates a deep copy of the Options
func (o Options) Clone() Options {
	result := NewDefaultOptions()

	// Copy pre-process options
	for k, v := range o.PreProcessOptions {
		result.PreProcessOptions[k] = v
	}

	// Copy LLM options
	for k, v := range o.LLMOptions {
		result.LLMOptions[k] = v
	}

	// Copy post-process options
	for k, v := range o.PostProcessOptions {
		result.PostProcessOptions[k] = v
	}

	return result
}

// WithLLMOption adds an LLM option and returns the updated Options
func (o Options) WithLLMOption(key string, value interface{}) Options {
	result := o.Clone()
	result.LLMOptions[key] = value
	return result
}

// WithPreProcessOption adds a pre-process option and returns the updated Options
func (o Options) WithPreProcessOption(key string, value interface{}) Options {
	result := o.Clone()
	result.PreProcessOptions[key] = value
	return result
}

// WithPostProcessOption adds a post-process option and returns the updated Options
func (o Options) WithPostProcessOption(key string, value interface{}) Options {
	result := o.Clone()
	result.PostProcessOptions[key] = value
	return result
}

// WithDebug sets the debug mode for the processor
func (o Options) WithDebug(debug bool) Options {
	result := o.Clone()
	result.LLMOptions["debug"] = debug
	return result
}

// GetDebugEnabled returns whether debug mode is enabled
func (o Options) GetDebugEnabled() bool {
	if o.LLMOptions == nil {
		return false
	}

	if debug, ok := o.LLMOptions["debug"].(bool); ok {
		return debug
	}
	return false
}
