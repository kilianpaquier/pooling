package pooling_test

import (
	"errors"
	"testing"

	"github.com/panjf2000/ants/v2"

	pooling "github.com/kilianpaquier/pooling/pkg"
)

func TestBuild(t *testing.T) {
	t.Run("error_no_sizes", func(t *testing.T) {
		// Act
		_, err := pooling.NewPoolerBuilder().Build()

		// Assert
		if !errors.Is(err, pooling.ErrMinimalSizes) {
			t.Fatal(err)
		}
	})

	t.Run("error_pool_creation", func(t *testing.T) {
		// Act
		_, err := pooling.NewPoolerBuilder().
			SetSizes(0, 5).
			SetOptions(ants.WithPreAlloc(true)).
			Build()

		// Assert
		if !errors.Is(err, ants.ErrInvalidPreAllocSize) {
			t.Fatal(err)
		}
	})

	t.Run("success", func(t *testing.T) {
		// Act & Assert
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 2, 3).
			Build()
		if err != nil {
			t.Fatal(err)
		}
		pooler.Close()
	})
}
