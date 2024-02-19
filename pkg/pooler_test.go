package pooling_test

import (
	"testing"

	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pooling "github.com/kilianpaquier/pooling/pkg"
)

func TestRead(t *testing.T) {
	t.Run("error_panic", func(t *testing.T) {
		// Arrange
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1).
			SetOptions(ants.WithNonblocking(true)).
			Build()
		require.NoError(t, err)
		t.Cleanup(pooler.Close)

		input := make(chan pooling.PoolerFunc)

		// Act
		go func() {
			defer close(input)
			triggerPanic(input)
		}()
		// & Assert
		assert.Panics(t, func() { pooler.Read(input) })
	})

	t.Run("error_panic_pool", func(t *testing.T) {
		// Arrange
		panicked := false
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 1).
			SetOptions(
				ants.WithNonblocking(true),
				ants.WithPanicHandler(func(_ any) { panicked = true }),
			).
			Build()
		require.NoError(t, err)
		t.Cleanup(pooler.Close)

		input := make(chan pooling.PoolerFunc)

		// Act
		go func() {
			defer close(input)
			input <- func(funcs chan<- pooling.PoolerFunc) {
				triggerPanic(funcs)
			}
		}()
		// & Assert
		assert.NotPanics(t, func() { pooler.Read(input) })

		// Assert
		assert.True(t, panicked)
	})

	t.Run("success_less_pools", func(t *testing.T) {
		// Arrange
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 1, 1). // less pools than limit
			Build()
		require.NoError(t, err)
		t.Cleanup(pooler.Close)

		limit := 5    // do 5 levels of recursive calls
		subcalls := 0 // subcalls should be 5 at the end

		input := make(chan pooling.PoolerFunc)

		// Act
		go func() {
			defer close(input)
			input <- poolerCaller(&subcalls, limit)
		}()
		pooler.Read(input)

		// Assert
		assert.Equal(t, limit, subcalls)
	})

	t.Run("success_same_pools", func(t *testing.T) {
		// Arrange
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 1, 1, 1, 1). // same number of pools as limit
			Build()
		require.NoError(t, err)
		t.Cleanup(pooler.Close)

		limit := 5    // do 5 levels of recursive calls
		subcalls := 0 // subcalls should be 5 at the end

		input := make(chan pooling.PoolerFunc)

		// Act
		go func() {
			defer close(input)
			input <- poolerCaller(&subcalls, limit)
		}()
		pooler.Read(input)

		// Assert
		assert.Equal(t, limit, subcalls)
	})

	t.Run("success_more_pools", func(t *testing.T) {
		// Arrange
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 1, 1, 1, 1, 1, 1). // more pools than limit
			Build()
		require.NoError(t, err)
		t.Cleanup(pooler.Close)

		limit := 5    // do 5 levels of recursive calls
		subcalls := 0 // subcalls should be 5 at the end

		input := make(chan pooling.PoolerFunc)

		// Act
		go func() {
			defer close(input)
			input <- poolerCaller(&subcalls, limit)
		}()
		pooler.Read(input)

		// Assert
		assert.Equal(t, limit, subcalls)
	})
}

// poolerCaller takes as input a ptr calls and a limit integer and send recursively poolerCaller until calls reaches limit.
func poolerCaller(calls *int, limit int) pooling.PoolerFunc {
	return func(funcs chan<- pooling.PoolerFunc) {
		if *calls < limit {
			*calls++
			funcs <- poolerCaller(calls, limit)
		}
	}
}

// triggerPanic only triggers a panic with Pooler Read function if the pool size associated to input channel is one.
func triggerPanic(funcs chan<- pooling.PoolerFunc) {
	blocker := make(chan struct{})
	defer close(blocker)

	// push a first blocking function
	funcs <- func(_ chan<- pooling.PoolerFunc) {
		blocker <- struct{}{}
	}

	// push a second function to expose panic
	funcs <- func(_ chan<- pooling.PoolerFunc) {}

	<-blocker // consume blocking message to stop execution
}
