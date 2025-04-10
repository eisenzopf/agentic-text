package data

import (
	"context"
	"runtime"
	"sync"
)

// DefaultWorkers is the default number of parallel workers
const DefaultWorkers = 4

// ParallelProcessor processes data using multiple goroutines
type ParallelProcessor struct {
	batchProcessor *BatchProcessor
	maxWorkers     int
}

// NewParallelProcessor creates a new parallel processor
func NewParallelProcessor(source Source, batchSize, maxWorkers int) *ParallelProcessor {
	if maxWorkers <= 0 {
		maxWorkers = DefaultWorkers
	}

	// Cap the number of workers to the number of CPU cores
	if maxWorkers > runtime.NumCPU() {
		maxWorkers = runtime.NumCPU()
	}

	return &ParallelProcessor{
		batchProcessor: NewBatchProcessor(source, batchSize),
		maxWorkers:     maxWorkers,
	}
}

// Close closes the underlying batch processor
func (p *ParallelProcessor) Close() error {
	return p.batchProcessor.Close()
}

// ProcessBatch processes a batch of data in parallel
func (p *ParallelProcessor) ProcessBatch(ctx context.Context, processor func(ctx context.Context, item *TextItem) (*TextItem, error)) ([]*TextItem, error) {
	batch, err := p.batchProcessor.NextBatch(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*TextItem, len(batch))
	errs := make([]error, len(batch))

	// Use a semaphore to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, p.maxWorkers)
	var wg sync.WaitGroup

	for i, item := range batch {
		wg.Add(1)
		go func(i int, item *TextItem) {
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

// ProcessAll processes all data in parallel
func (p *ParallelProcessor) ProcessAll(ctx context.Context, processor func(ctx context.Context, item *TextItem) (*TextItem, error)) ([]*TextItem, error) {
	var allResults []*TextItem
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
