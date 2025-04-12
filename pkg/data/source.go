package data

import (
	"context"
	"io"
)

// ProcessItemSource defines an interface for sources that can directly provide ProcessItems
type ProcessItemSource interface {
	// NextProcessItem returns the next ProcessItem or error when exhausted
	NextProcessItem(context.Context) (*ProcessItem, error)
	// Close releases any resources used by the source
	Close() error
}

// ProcessItemSliceSource implements ProcessItemSource for a slice of ProcessItems
type ProcessItemSliceSource struct {
	items []*ProcessItem
	index int
}

// NewProcessItemSliceSource creates a new source from a slice of ProcessItems
func NewProcessItemSliceSource(items []*ProcessItem) *ProcessItemSliceSource {
	return &ProcessItemSliceSource{
		items: items,
		index: 0,
	}
}

// NextProcessItem implements the ProcessItemSource interface
func (s *ProcessItemSliceSource) NextProcessItem(_ context.Context) (*ProcessItem, error) {
	if s.index >= len(s.items) {
		return nil, io.EOF
	}

	item := s.items[s.index]
	s.index++
	return item, nil
}

// Close implements the ProcessItemSource interface
func (s *ProcessItemSliceSource) Close() error {
	return nil
}

// TextStringsProcessItemSource is a convenience wrapper for a slice of strings as ProcessItems
type TextStringsProcessItemSource struct {
	items []*ProcessItem
	index int
}

// NewTextStringsProcessItemSource creates a new ProcessItemSource from a slice of strings
func NewTextStringsProcessItemSource(strings []string) *TextStringsProcessItemSource {
	items := make([]*ProcessItem, len(strings))
	for i, str := range strings {
		items[i] = NewTextProcessItem("", str, nil)
	}

	return &TextStringsProcessItemSource{
		items: items,
		index: 0,
	}
}

// NextProcessItem implements the ProcessItemSource interface
func (s *TextStringsProcessItemSource) NextProcessItem(_ context.Context) (*ProcessItem, error) {
	if s.index >= len(s.items) {
		return nil, io.EOF
	}

	item := s.items[s.index]
	s.index++
	return item, nil
}

// Close implements the ProcessItemSource interface
func (s *TextStringsProcessItemSource) Close() error {
	return nil
}
