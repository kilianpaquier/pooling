package pooling

import (
	"sync"

	"github.com/panjf2000/ants/v2"
)

// PoolerFunc represents a function to be given into Pooler channel consumption.
//
// Input channel closing is handled by Pooler so don't close it yourself.
type PoolerFunc func(funcs chan<- PoolerFunc)

// Pooler represents a slice of pools alongside a waitgroup to handle functions (like a queue).
//
// A Pooler contains a slice of pools to allow each functions (given as channel in Read)
// to send functions (recursively and indefinitely) into the next pool.
type Pooler struct {
	pools []*ants.Pool
	wg    sync.WaitGroup
}

// Close waits for all Pooler funcs to be ended and then closes all Pooler pools.
func (p *Pooler) Close() {
	p.wg.Wait() // wait should give back hand directly because Read also waits for all pools to be free
	for _, pool := range p.pools {
		pool.Release()
	}
}

// Read reads indefinitely (until closed) the input channel.
// It will wait at the end (when closed) for all functions executions to be ended before giving back the hand.
func (p *Pooler) Read(funcs <-chan PoolerFunc) {
	p.readAt(0, funcs)
	p.wg.Wait()
}

// readAt reads the input channel functions in current thread
// and executes them in index's pool.
//
// if index's pool is outside of Pooler pool size, then execution is done in current thread.
func (p *Pooler) readAt(index int, funcs <-chan PoolerFunc) {
	// execute funcs in current thread if the index is outside of pool size
	if index > len(p.pools)-1 {
		for f := range funcs {
			p.runFunc(f, index)
		}
		return
	}

	for f := range funcs {
		// add one to gloal Pooler running routines to ensure at the end of Read
		// that all running routines end
		p.wg.Add(1)

		// submit f execution to the index's pool
		err := p.pools[index].Submit(func() {
			defer p.wg.Done()
			p.runFunc(f, index)
		})
		if err != nil { // f could not be pushed, panic and set f as done
			p.wg.Done()
			panic(err)
		}
	}
}

// runFunc adds one to pooler waitgroup and returns the function handling the input PoolerFunc.
func (p *Pooler) runFunc(f PoolerFunc, index int) {
	// un-buffered channel to avoid waiting tasks in pools
	children := make(chan PoolerFunc)

	go func() {
		// close children after f sent them all
		defer close(children)
		f(children) // execute f sending elements to children
	}()

	p.readAt(index+1, children) // read children elements
}
