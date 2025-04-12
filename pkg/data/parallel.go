package data

import (
	"context"
	"runtime"
	"sync"
)

// DefaultWorkers is the default number of parallel workers
const DefaultWorkers = 4

// ProcessItemParallelProcessor processes ProcessItems using multiple goroutines
type ProcessItemParallelProcessor struct {
	batchProcessor *ProcessItemBatchProcessor
	maxWorkers     int
}

// NewProcessItemParallelProcessor creates a new parallel processor for ProcessItems
func NewProcessItemParallelProcessor(source ProcessItemSource, batchSize, maxWorkers int) *ProcessItemParallelProcessor {
	if maxWorkers <= 0 {
		maxWorkers = DefaultWorkers
	}

	// Cap the number of workers to the number of CPU cores
	if maxWorkers > runtime.NumCPU() {
		maxWorkers = runtime.NumCPU()
	}

	return &ProcessItemParallelProcessor{
		batchProcessor: NewProcessItemBatchProcessor(source, batchSize),
		maxWorkers:     maxWorkers,
	}
}

// Close closes the underlying batch processor
func (p *ProcessItemParallelProcessor) Close() error {
	return p.batchProcessor.Close()
}

// ProcessBatch processes a batch of ProcessItems in parallel
func (p *ProcessItemParallelProcessor) ProcessBatch(ctx context.Context, processor func(ctx context.Context, item *ProcessItem) (*ProcessItem, error)) ([]*ProcessItem, error) {
	batch, err := p.batchProcessor.NextBatch(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*ProcessItem, len(batch))
	errs := make([]error, len(batch))

	// Use a semaphore to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, p.maxWorkers)
	var wg sync.WaitGroup

	for i, item := range batch {
		wg.Add(1)
		go func(i int, item *ProcessItem) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Process the item
			result, err := processor(ctx, item)
			results[i] = result
			errs[i] = err
		}(i, item)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ProcessAll processes all ProcessItems in parallel
func (p *ProcessItemParallelProcessor) ProcessAll(ctx context.Context, processor func(ctx context.Context, item *ProcessItem) (*ProcessItem, error)) ([]*ProcessItem, error) {
	var allResults []*ProcessItem
	var mu sync.Mutex

	// Process batches sequentially, but items within each batch in parallel
	for {
		results, err := p.ProcessBatch(ctx, processor)
		if err == nil {
			mu.Lock()
			allResults = append(allResults, results...)
			mu.Unlock()
			continue
		}

		// Break on EOF
		if err.Error() == "EOF" {
			break
		}

		return nil, err
	}

	return allResults, nil
}
