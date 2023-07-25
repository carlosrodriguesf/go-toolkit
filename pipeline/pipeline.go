//go:generate mockery --all --case snake --outpkg ${GOPACKAGE}mocks --exported

package pipeline

import (
	"context"
	"errors"
)

type (
	Next[T any] func(context.Context, T) (T, error)
	Step[T any] interface {
		Run(context.Context, T, Next[T]) (T, error)
		Name() string
	}
	Pipeline[T any] []Step[T]
)

func (p Pipeline[T]) Run(ctx context.Context, data T) (T, error) {
	return p.runStep(ctx, 0, data)
}

func (p Pipeline[T]) runStep(ctx context.Context, stepIndex int, data T) (T, error) {
	if stepIndex >= len(p) {
		return data, nil
	}

	step := p[stepIndex]

	select {
	case <-ctx.Done():
		return data, p.checkError(step, ctx.Err())
	default:
		next := func(ctx context.Context, e T) (T, error) { return p.runStep(ctx, stepIndex+1, e) }
		data, err := step.Run(ctx, data, next)
		if err != nil {
			return *new(T), p.checkError(step, err)
		}
		return data, nil
	}
}

func (p Pipeline[T]) checkError(step Step[T], err error) error {
	var runErr *RunError
	if errors.As(err, &runErr) {
		return err
	}
	return &RunError{
		StepName: step.Name(),
		Err:      err,
	}
}
