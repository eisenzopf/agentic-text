package data

import (
	"context"
	"io"
)

// TextItem represents a single text snippet with optional metadata
type TextItem struct {
	// ID is an optional identifier for the text
	ID string
	// Content is the actual text content
	Content string
	// Metadata is optional key-value pairs associated with the text
	Metadata map[string]interface{}
}

// Source defines the interface for data sources
type Source interface {
	// Next returns the next text item or error when exhausted
	Next(context.Context) (*TextItem, error)
	// Close releases any resources used by the source
	Close() error
}

// SliceSource implements Source for a slice of text items
type SliceSource struct {
	items []*TextItem
	index int
}

// NewSliceSource creates a new source from a slice of text items
func NewSliceSource(items []*TextItem) *SliceSource {
	return &SliceSource{
		items: items,
		index: 0,
	}
}

// Next implements the Source interface
func (s *SliceSource) Next(_ context.Context) (*TextItem, error) {
	if s.index >= len(s.items) {
		return nil, io.EOF
	}

	item := s.items[s.index]
	s.index++
	return item, nil
}

// Close implements the Source interface
func (s *SliceSource) Close() error {
	return nil
}

// StringsSource is a convenience wrapper for a slice of strings
type StringsSource struct {
	source *SliceSource
}

// NewStringsSource creates a new source from a slice of strings
func NewStringsSource(strings []string) *StringsSource {
	items := make([]*TextItem, len(strings))
	for i, str := range strings {
		items[i] = &TextItem{
			Content: str,
		}
	}

	return &StringsSource{
		source: NewSliceSource(items),
	}
}

// Next implements the Source interface
func (s *StringsSource) Next(ctx context.Context) (*TextItem, error) {
	return s.source.Next(ctx)
}

// Close implements the Source interface
func (s *StringsSource) Close() error {
	return s.source.Close()
}
