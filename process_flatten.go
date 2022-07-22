package pipeline

import (
	"context"
	"github.com/Sora233/pipeline/v2/semaphore"
)

func ProcessFlatten[Input, Output any](ctx context.Context, processor Processor[Input, []Output], in <-chan Input) <-chan Output {
	out := make(chan []Output)
	go func() {
		for i := range in {
			process(ctx, processor, i, out)
		}
		close(out)
	}()
	return Split(out)
}

func ProcessConcurrentlyFlattern[Input, Output any](ctx context.Context, concurrently int, p Processor[Input, []Output], in <-chan Input) <-chan Output {
	// Create the out chan
	out := make(chan []Output)
	go func() {
		// Perform Process concurrently times
		sem := semaphore.New(concurrently)
		for i := range in {
			sem.Add(1)
			go func(i Input) {
				process(ctx, p, i, out)
				sem.Done()
			}(i)
		}
		// Close the out chan after all of the Processors finish executing
		sem.Wait()
		close(out)
	}()
	return Split(out)
}
