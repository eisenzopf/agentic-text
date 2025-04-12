package data

import (
	"context"
	"io"
)

// DefaultBatchSize is the default number of items to batch together
const DefaultBatchSize = 10

// ProcessItemBatchProcessor processes ProcessItems in batches
type ProcessItemBatchProcessor struct {
	source       ProcessItemSource
	batchSize    int
	currentBatch []*ProcessItem
}

// NewProcessItemBatchProcessor creates a new batch processor for ProcessItems
func NewProcessItemBatchProcessor(source ProcessItemSource, batchSize int) *ProcessItemBatchProcessor {
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}

	return &ProcessItemBatchProcessor{
		source:       source,
		batchSize:    batchSize,
		currentBatch: make([]*ProcessItem, 0, batchSize),
	}
}

// NextBatch returns the next batch of ProcessItems
func (b *ProcessItemBatchProcessor) NextBatch(ctx context.Context) ([]*ProcessItem, error) {
	// If we already have items in the current batch, return those
	if len(b.currentBatch) > 0 {
		batch := b.currentBatch
		b.currentBatch = make([]*ProcessItem, 0, b.batchSize)
		return batch, nil
	}

	// Otherwise, fetch a new batch
	batch := make([]*ProcessItem, 0, b.batchSize)
	for i := 0; i < b.batchSize; i++ {
		item, err := b.source.NextProcessItem(ctx)
		if err == io.EOF {
			if len(batch) == 0 {
				return nil, io.EOF
			}
			break
		}
		if err != nil {
			return nil, err
		}
		batch = append(batch, item)
	}

	return batch, nil
}

// Close closes the underlying source
func (b *ProcessItemBatchProcessor) Close() error {
	return b.source.Close()
}

// Process applies a processor function to each ProcessItem in a batch
func (b *ProcessItemBatchProcessor) Process(ctx context.Context, processor func(ctx context.Context, item *ProcessItem) (*ProcessItem, error)) ([]*ProcessItem, error) {
	batch, err := b.NextBatch(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*ProcessItem, len(batch))
	for i, item := range batch {
		result, err := processor(ctx, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ProcessAll applies a processor function to all remaining ProcessItems
func (b *ProcessItemBatchProcessor) ProcessAll(ctx context.Context, processor func(ctx context.Context, item *ProcessItem) (*ProcessItem, error)) ([]*ProcessItem, error) {
	var allResults []*ProcessItem

	for {
		results, err := b.Process(ctx, processor)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}
