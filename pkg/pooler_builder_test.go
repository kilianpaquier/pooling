package pooling_test

import (
	"testing"

	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"

	pooling "github.com/kilianpaquier/pooling/pkg"
)

func TestBuild(t *testing.T) {
	t.Run("error_no_sizes", func(t *testing.T) {
		// Act
		_, err := pooling.NewPoolerBuilder().Build()

		// Assert
		assert.Equal(t, pooling.ErrMinimalSizes, err)
	})

	t.Run("error_pool_creation", func(t *testing.T) {
		// Act
		_, err := pooling.NewPoolerBuilder().
			SetSizes(0, 5).
			SetOptions(ants.WithPreAlloc(true)).
			Build()

		// Assert
		assert.ErrorContains(t, err, ants.ErrInvalidPreAllocSize.Error())
	})

	t.Run("success", func(t *testing.T) {
		// Act
		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(1, 2, 3).
			Build()

		// Assert
		assert.NoError(t, err)
		pooler.Close()
	})
}
