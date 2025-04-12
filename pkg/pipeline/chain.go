package pipeline

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Chain represents a pipeline of processors
type Chain struct {
	processors []processor.Processor
	name       string
}

// NewChain creates a new processor chain
func NewChain(name string, processors ...processor.Processor) *Chain {
	return &Chain{
		processors: processors,
		name:       name,
	}
}

// Process processes a ProcessItem through the entire chain
func (c *Chain) Process(ctx context.Context, item *data.ProcessItem) (*data.ProcessItem, error) {
	if len(c.processors) == 0 {
		return nil, fmt.Errorf("empty processor chain")
	}

	// Process with the first processor
	result, err := c.processors[0].Process(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("processor '%s' error: %w", c.processors[0].GetName(), err)
	}

	// Process with remaining processors, using the result from the previous step
	for i := 1; i < len(c.processors); i++ {
		proc := c.processors[i]
		result, err = proc.Process(ctx, result)
		if err != nil {
			return nil, fmt.Errorf("processor '%s' error: %w", proc.GetName(), err)
		}
	}

	return result, nil
}

// ProcessSource processes a data source through the chain
func (c *Chain) ProcessSource(ctx context.Context, source data.ProcessItemSource, batchSize, workers int) ([]*data.ProcessItem, error) {
	if len(c.processors) == 0 {
		return nil, fmt.Errorf("empty processor chain")
	}

	// Use the first processor to process the source
	firstResults, err := c.processors[0].ProcessSource(ctx, source, batchSize, workers)
	if err != nil {
		return nil, err
	}

	// If there's only one processor, return the results
	if len(c.processors) == 1 {
		return firstResults, nil
	}

	// Process the results through the remaining processors
	currentResults := firstResults
	for i := 1; i < len(c.processors); i++ {
		proc := c.processors[i]

		// Process with the next processor
		nextResults, err := proc.ProcessBatch(ctx, currentResults)
		if err != nil {
			return nil, err
		}

		currentResults = nextResults
	}

	return currentResults, nil
}

// GetName returns the chain name
func (c *Chain) GetName() string {
	return c.name
}

// ProcessBatch processes a batch of items through the chain
func (c *Chain) ProcessBatch(ctx context.Context, items []*data.ProcessItem) ([]*data.ProcessItem, error) {
	if len(c.processors) == 0 {
		return nil, fmt.Errorf("empty processor chain")
	}

	// Process with the first processor
	currentResults, err := c.processors[0].ProcessBatch(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("processor '%s' error: %w", c.processors[0].GetName(), err)
	}

	// Process with remaining processors
	for i := 1; i < len(c.processors); i++ {
		proc := c.processors[i]
		currentResults, err = proc.ProcessBatch(ctx, currentResults)
		if err != nil {
			return nil, fmt.Errorf("processor '%s' error: %w", proc.GetName(), err)
		}
	}

	return currentResults, nil
}
