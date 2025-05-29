package processor

import (
	"context"
	"fmt"
	"strings"
)

// ProcessorBuilder provides a fluent interface for creating processors
type ProcessorBuilder struct {
	name            string
	resultStruct    interface{}
	contentTypes    []string
	role            string
	objective       string
	instructions    []string
	customSections  map[string]string
	customPromptGen PromptGenerator
	customInit      func(*GenericProcessor) error
	validateStruct  bool
}

// NewBuilder creates a new processor builder
func NewBuilder(name string) *ProcessorBuilder {
	return &ProcessorBuilder{
		name:           name,
		contentTypes:   []string{"text"}, // sensible default
		customSections: make(map[string]string),
		validateStruct: false, // sensible default
	}
}

// WithStruct sets the result structure (required)
func (b *ProcessorBuilder) WithStruct(resultStruct interface{}) *ProcessorBuilder {
	b.resultStruct = resultStruct
	return b
}

// WithContentTypes sets supported content types
func (b *ProcessorBuilder) WithContentTypes(types ...string) *ProcessorBuilder {
	b.contentTypes = types
	return b
}

// WithRole sets the AI role for the prompt
func (b *ProcessorBuilder) WithRole(role string) *ProcessorBuilder {
	b.role = role
	return b
}

// WithObjective sets the main objective for the prompt
func (b *ProcessorBuilder) WithObjective(objective string) *ProcessorBuilder {
	b.objective = objective
	return b
}

// WithInstructions sets step-by-step instructions
func (b *ProcessorBuilder) WithInstructions(instructions ...string) *ProcessorBuilder {
	b.instructions = instructions
	return b
}

// WithCustomSection adds a custom section to the prompt
func (b *ProcessorBuilder) WithCustomSection(name, content string) *ProcessorBuilder {
	b.customSections[name] = content
	return b
}

// WithCustomPrompt replaces the auto-generated prompt with a custom one
func (b *ProcessorBuilder) WithCustomPrompt(promptGen PromptGenerator) *ProcessorBuilder {
	b.customPromptGen = promptGen
	return b
}

// WithCustomInit sets a custom initialization function
func (b *ProcessorBuilder) WithCustomInit(initFunc func(*GenericProcessor) error) *ProcessorBuilder {
	b.customInit = initFunc
	return b
}

// WithValidation enables struct-level validation
func (b *ProcessorBuilder) WithValidation() *ProcessorBuilder {
	b.validateStruct = true
	return b
}

// Register creates and registers the processor
func (b *ProcessorBuilder) Register() {
	if b.resultStruct == nil {
		panic(fmt.Sprintf("processor %s: result struct is required", b.name))
	}

	var promptGen PromptGenerator
	if b.customPromptGen != nil {
		// Use custom prompt generator
		promptGen = b.customPromptGen
	} else {
		// Create auto-generated prompt generator
		promptGen = &BuilderPromptGenerator{
			resultStruct:   b.resultStruct,
			role:           b.role,
			objective:      b.objective,
			instructions:   b.instructions,
			customSections: b.customSections,
		}
	}

	RegisterGenericProcessor(
		b.name,
		b.contentTypes,
		b.resultStruct,
		promptGen,
		b.customInit,
		b.validateStruct,
	)
}

// BuilderPromptGenerator generates prompts based on builder configuration
type BuilderPromptGenerator struct {
	resultStruct   interface{}
	role           string
	objective      string
	instructions   []string
	customSections map[string]string
}

// GeneratePrompt implements PromptGenerator interface
func (p *BuilderPromptGenerator) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	jsonExample := GenerateJSONExample(p.resultStruct)

	var promptParts []string

	// Add role if specified
	if p.role != "" {
		promptParts = append(promptParts, fmt.Sprintf("**Role:** %s", p.role))
	}

	// Add objective if specified
	if p.objective != "" {
		promptParts = append(promptParts, fmt.Sprintf("**Objective:** %s", p.objective))
	}

	// Add input text
	promptParts = append(promptParts, fmt.Sprintf("**Input Text:**\n%s", text))

	// Add instructions if specified
	if len(p.instructions) > 0 {
		instructionText := "**Instructions:**\n"
		for i, instruction := range p.instructions {
			instructionText += fmt.Sprintf("%d. %s\n", i+1, instruction)
		}
		promptParts = append(promptParts, instructionText)
	}

	// Add custom sections
	for name, content := range p.customSections {
		promptParts = append(promptParts, fmt.Sprintf("**%s:**\n%s", name, content))
	}

	// Always add JSON structure requirement
	promptParts = append(promptParts, fmt.Sprintf("**Required JSON Output Structure:**\n%s", jsonExample))

	// Always add critical JSON-only instruction
	promptParts = append(promptParts, "*** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***")

	return strings.Join(promptParts, "\n\n"), nil
}
