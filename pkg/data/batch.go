package data

import (
	"context"
	"io"
)

// DefaultBatchSize is the default number of items to batch together
const DefaultBatchSize = 10

// BatchProcessor processes data in batches
type BatchProcessor struct {
	source       Source
	batchSize    int
	currentBatch []*TextItem
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(source Source, batchSize int) *BatchProcessor {
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}

	return &BatchProcessor{
		source:       source,
		batchSize:    batchSize,
		currentBatch: make([]*TextItem, 0, batchSize),
	}
}

// NextBatch returns the next batch of text items
func (b *BatchProcessor) NextBatch(ctx context.Context) ([]*TextItem, error) {
	// If we already have items in the current batch, return those
	if len(b.currentBatch) > 0 {
		batch := b.currentBatch
		b.currentBatch = make([]*TextItem, 0, b.batchSize)
		return batch, nil
	}

	// Otherwise, fetch a new batch
	batch := make([]*TextItem, 0, b.batchSize)
	for i := 0; i < b.batchSize; i++ {
		item, err := b.source.Next(ctx)
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
func (b *BatchProcessor) Close() error {
	return b.source.Close()
}

// Process applies a processor function to each item in a batch
func (b *BatchProcessor) Process(ctx context.Context, processor func(ctx context.Context, item *TextItem) (*TextItem, error)) ([]*TextItem, error) {
	batch, err := b.NextBatch(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*TextItem, len(batch))
	for i, item := range batch {
		result, err := processor(ctx, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ProcessAll applies a processor function to all remaining items
func (b *BatchProcessor) ProcessAll(ctx context.Context, processor func(ctx context.Context, item *TextItem) (*TextItem, error)) ([]*TextItem, error) {
	var allResults []*TextItem

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
