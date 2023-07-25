package pipeline_test

import (
	"context"
	"errors"
	"github.com/carlosrodriguesf/go-toolkit/pipeline"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPipeline_Run(t *testing.T) {
	step1 := getMockedStep("step1", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})
	step2 := getMockedStep("step2", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})
	step3 := getMockedStep("step3", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})

	p := pipeline.Pipeline[int]{
		step1,
		step2,
		step3,
	}
	result, err := p.Run(context.TODO(), 0)
	require.NoError(t, err)
	require.Equal(t, 3, result)
}

func TestPipeline_Run_StopsOn2(t *testing.T) {
	step1 := getMockedStep("step1", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})
	step2 := getMockedStep("step2", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return -1, nil
	})
	step3 := getMockedStep("step3", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return 0, errors.New("test error")
	})
	step4 := getMockedStep("step4", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return i + 1, nil
	})
	p := pipeline.Pipeline[int]{
		step1,
		step2,
		step3,
		step4,
	}
	data, err := p.Run(context.TODO(), 0)
	require.NoError(t, err)
	require.Equal(t, -1, data)
}

func TestPipeline_Run_Error(t *testing.T) {
	step1 := getMockedStep("step1", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})
	step2 := getMockedStep("step2", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})
	step3 := getMockedStep("step3", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return 0, errors.New("test error")
	})
	step4 := getMockedStep("step4", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return i + 1, nil
	})
	p := pipeline.Pipeline[int]{
		step1,
		step2,
		step3,
		step4,
	}
	_, err := p.Run(context.TODO(), 0)
	require.EqualError(t, err, "error running step 'step3': test error")
}

func TestPipeline_Run_Error_ContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	step1 := getMockedStep("step1", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return next(ctx, i+1)
	})
	step2 := getMockedStep("step2", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		cancel()
		return next(ctx, i+1)
	})
	step3 := getMockedStep("step3", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return 0, errors.New("test error")
	})
	step4 := getMockedStep("step4", func(ctx context.Context, i int, next pipeline.Next[int]) (int, error) {
		return i + 1, nil
	})
	p := pipeline.Pipeline[int]{
		step1,
		step2,
		step3,
		step4,
	}
	_, err := p.Run(ctx, 0)
	require.EqualError(t, err, "error running step 'step3': context canceled")
}

type (
	mockStepCall[T any] func(context.Context, T, pipeline.Next[T]) (T, error)
	mockStep[T any]     struct {
		name string
		call mockStepCall[T]
	}
)

func (t mockStep[T]) Run(ctx context.Context, data T, next pipeline.Next[T]) (T, error) {
	return t.call(ctx, data, next)
}

func (t mockStep[T]) Name() string {
	return t.name
}

func getMockedStep[T any](name string, call mockStepCall[T]) mockStep[T] {
	return mockStep[T]{
		name: name,
		call: call,
	}
}
