package data

import (
	"encoding/json"
	"fmt"
)

// ProcessItem represents a standard item flowing through processors
type ProcessItem struct {
	// ID for tracking the item
	ID string `json:"id"`

	// Content holds the actual data (could be string, object, etc.)
	Content interface{} `json:"content"`

	// ContentType indicates how to interpret the content
	ContentType string `json:"content_type"`

	// Metadata for additional information
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// ProcessingInfo contains history and context of processing operations
	ProcessingInfo map[string]interface{} `json:"processing_info,omitempty"`
}

// NewTextProcessItem creates a new ProcessItem from a string
func NewTextProcessItem(id string, text string, metadata map[string]interface{}) *ProcessItem {
	return &ProcessItem{
		ID:             id,
		Content:        text,
		ContentType:    "text",
		Metadata:       metadata,
		ProcessingInfo: make(map[string]interface{}),
	}
}

// GetTextContent extracts the content as a string if it's text type
func (p *ProcessItem) GetTextContent() (string, error) {
	if p.ContentType != "text" {
		return "", fmt.Errorf("content type is not text: %s", p.ContentType)
	}

	if content, ok := p.Content.(string); ok {
		return content, nil
	}

	return "", fmt.Errorf("content cannot be converted to string")
}

// AddProcessingInfo adds information about a processing step
func (p *ProcessItem) AddProcessingInfo(processorName string, info interface{}) {
	if p.ProcessingInfo == nil {
		p.ProcessingInfo = make(map[string]interface{})
	}
	p.ProcessingInfo[processorName] = info
}

// Clone creates a deep copy of the ProcessItem
func (p *ProcessItem) Clone() (*ProcessItem, error) {
	// Convert to JSON and back to create a deep copy
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	var clone ProcessItem
	if err := json.Unmarshal(data, &clone); err != nil {
		return nil, err
	}

	return &clone, nil
}
