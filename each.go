package flowmatic

import (
	"errors"
	"iter"
)

// Each starts numWorkers concurrent workers (or GOMAXPROCS workers if numWorkers < 1)
// and processes each item as a task.
// Errors returned by a task do not halt execution,
// but are joined into a multierror return value.
// If a task panics during execution,
// the panic will be caught and rethrown in the parent Goroutine.
func Each[Input any](numWorkers int, seq iter.Seq[Input], task func(Input) error) error {
	type void struct{}

	inch, ouch := TaskPool(numWorkers, func(in Input) (void, error) {
		return void{}, task(in)
	})

	var errs []error

	_ = Do(
		func() error {
			defer close(inch)

			for in := range seq {
				inch <- in
			}
			return nil
		},
		func() error {
			for out := range ouch {
				if out.Panic != nil {
					panic(out.Panic)
				}
				if err := out.Err; err != nil {
					errs = append(errs, err)
				}
			}
			return nil
		},
	)
	return errors.Join(errs...)
}
