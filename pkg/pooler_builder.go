package pooling

import (
	"errors"

	"github.com/panjf2000/ants/v2"
)

// ErrMinimalSizes is the error returned by Build in case the PoolerBuilder doesn't have any size of pool to create.
var ErrMinimalSizes = errors.New("pooler pools size must be at least 1")

// PoolerBuilder is the builder for Pooler.
// It takes input options to tune Pooler behavior.
type PoolerBuilder struct {
	options []ants.Option
	sizes   []int16
}

// NewPoolerBuilder creates a new PoolerBuilder.
func NewPoolerBuilder() *PoolerBuilder {
	return &PoolerBuilder{}
}

// SetOptions takes a slice of ants options for Pooler pools.
func (p *PoolerBuilder) SetOptions(options ...ants.Option) *PoolerBuilder {
	p.options = options
	return p
}

// SetSizes takes a slice of integers where each one will represent the size of an ants pool.
func (p *PoolerBuilder) SetSizes(sizes ...int16) *PoolerBuilder {
	p.sizes = sizes
	return p
}

// Build builds the Pooler associated to PoolerBuilder
// and returns an error in case the ants pools creation fails.
func (p *PoolerBuilder) Build() (*Pooler, error) {
	if len(p.sizes) < 1 {
		return nil, ErrMinimalSizes
	}

	// initialize pooler
	pooler := &Pooler{
		pools: make([]*ants.Pool, 0, len(p.sizes)),
	}

	errs := make([]error, 0, len(p.sizes))
	for _, size := range p.sizes {
		// create pool with each provided size
		pool, err := ants.NewPool(int(size), p.options...)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		pooler.pools = append(pooler.pools, pool)
	}

	if len(errs) > 0 {
		pooler.Close() // call close directly to avoid memory leaks
		return nil, errors.Join(errs...)
	}
	return pooler, nil
}
