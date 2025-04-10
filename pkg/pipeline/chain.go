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

// Process processes a text through the entire chain
func (c *Chain) Process(ctx context.Context, text string) (*processor.Result, error) {
	if len(c.processors) == 0 {
		return nil, fmt.Errorf("empty processor chain")
	}

	var result *processor.Result
	var err error

	// Process with the first processor
	result, err = c.processors[0].Process(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("processor '%s' error: %w", c.processors[0].GetName(), err)
	}

	// Process with remaining processors, using the processed text from the previous step
	for i := 1; i < len(c.processors); i++ {
		proc := c.processors[i]
		result, err = proc.Process(ctx, result.Processed)
		if err != nil {
			return nil, fmt.Errorf("processor '%s' error: %w", proc.GetName(), err)
		}
	}

	return result, nil
}

// ProcessItem processes a data.TextItem through the entire chain
func (c *Chain) ProcessItem(ctx context.Context, item *data.TextItem) (*processor.Result, error) {
	if len(c.processors) == 0 {
		return nil, fmt.Errorf("empty processor chain")
	}

	var result *processor.Result
	var err error

	// Process with the first processor
	result, err = c.processors[0].ProcessItem(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("processor '%s' error: %w", c.processors[0].GetName(), err)
	}

	// Process with remaining processors, using the processed text from the previous step
	for i := 1; i < len(c.processors); i++ {
		proc := c.processors[i]
		// Create a new TextItem with the processed text
		nextItem := &data.TextItem{
			ID:       item.ID,
			Content:  result.Processed,
			Metadata: item.Metadata,
		}
		result, err = proc.ProcessItem(ctx, nextItem)
		if err != nil {
			return nil, fmt.Errorf("processor '%s' error: %w", proc.GetName(), err)
		}
	}

	return result, nil
}

// ProcessSource processes a data source through the chain
func (c *Chain) ProcessSource(ctx context.Context, source data.Source, batchSize, workers int) ([]*processor.Result, error) {
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

		// Convert results to TextItems for the next processor
		items := make([]*data.TextItem, len(currentResults))
		for j, result := range currentResults {
			// Get the ID and metadata from the original if available
			var id string
			var metadata map[string]interface{}

			if item, ok := result.Original.(*data.TextItem); ok {
				id = item.ID
				metadata = item.Metadata
			}

			items[j] = &data.TextItem{
				ID:       id,
				Content:  result.Processed,
				Metadata: metadata,
			}
		}

		// Create a source from the items
		nextSource := data.NewSliceSource(items)

		// Process with the next processor
		nextResults, err := proc.ProcessSource(ctx, nextSource, batchSize, workers)
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
