package pooling_test

import (
	"testing"

	"github.com/panjf2000/ants/v2"

	"github.com/kilianpaquier/pooling/internal/testutils"
	pooling "github.com/kilianpaquier/pooling/pkg"
)

func TestBuild(t *testing.T) {
	t.Run("error_no_sizes", func(t *testing.T) {
		// Act
		_, err := pooling.NewPoolerBuilder().Build()

		// Assert
		testutils.ErrorIs(t, err, pooling.ErrMinimalSizes)
	})

	t.Run("error_pool_creation", func(t *testing.T) {
		// Act
		_, err := pooling.NewPoolerBuilder().
			SetSizes(0, 5).
			SetOptions(ants.WithPreAlloc(true)).
			Build()

		// Assert
		testutils.ErrorIs(t, err, ants.ErrInvalidPreAllocSize)
	})

	t.Run("success", func(t *testing.T) {
		// Act & Assert
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 2, 3).
			Build()
		testutils.NoError(testutils.Require(t), err)
		pooler.Close()
	})
}
