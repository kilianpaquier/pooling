/*
The pooling package allows one to dispatch an infinite number of functions to be
executed in parallel while still limiting the number of routines.

For that, pooling package takes advantage of ants pool library.
A pooling Pooler can have multiple pools (with builder SetSizes) to dispatch sub functions into different pools of routines.

When sending a function into the pooler (with the appropriate channel), this function can itself send other functions into the pooler.
It allows one to "split" functions executions (like iterating over a slice and each element handled in parallel).

	func main() {
		log := logrus.WithContext(context.Background())

		pooler, err := pooling.NewPoolerBuilder().
			SetSizes(10, 500, ...). // each size will initialize a pool with given size
			SetOptions(ants.WithLogger(log)).
			Build()
		if err != nil {
			panic(err)
		}
		defer pooler.Close()

		input := ReadFrom()

		// Read function is blocking until input is closed
		// and all running routines have ended
		pooler.Read(input)
	}

	func ReadFrom() <-chan pooling.PoolerFunc {
		input := make(chan pooling.PoolerFunc)

		go func() {
			// close input to stop blocking function Read once all elements are sent to input
			defer close(input)

			// do something populating input channel
			for i := range 100 {
				input <- HandleInt(i)
			}
		}()

		return input
	}

	func HandleInt(i int) pooling.PoolerFunc {
		return func(funcs chan<- pooling.PoolerFunc) {
			// you may handle the integer whichever you want
			// funcs channel is present to dispatch again some elements into a channel handled by the pooler
		}
	}
*/
package pooling
