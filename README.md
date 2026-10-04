# pooling <!-- omit in toc -->

<div align="center">
  <a href="https://gitlab.com/kilianpaquier/pooling/-/releases">
    <img alt="GitLab Release" src="https://img.shields.io/gitlab/v/release/kilianpaquier%2Fpooling?gitlab_url=https%3A%2F%2Fgitlab.com&include_prereleases&sort=semver&style=for-the-badge">
  </a>
  <a href="https://gitlab.com/kilianpaquier/pooling/-/work_items">
    <img alt="GitLab Issues" src="https://img.shields.io/gitlab/issues/open/kilianpaquier%2Fpooling?gitlab_url=https%3A%2F%2Fgitlab.com&style=for-the-badge">
  </a>
  <a href="https://gitlab.com/kilianpaquier/pooling/-/blob/HEAD/LICENSE">
    <img alt="GitLab License" src="https://img.shields.io/gitlab/license/kilianpaquier%2Fpooling?gitlab_url=https%3A%2F%2Fgitlab.com&style=for-the-badge">
  </a>
  <a href="https://gitlab.com/kilianpaquier/pooling/-/pipelines?ref=main">
    <img alt="GitLab CICD" src="https://img.shields.io/gitlab/pipeline-status/kilianpaquier%2Fpooling?gitlab_url=https%3A%2F%2Fgitlab.com&branch=main&style=for-the-badge">
  </a>
  <a href="https://gitlab.com/kilianpaquier/pooling/-/blob/HEAD/go.mod">
    <img alt="Go Version" src="https://img.shields.io/gitlab/go-mod/go-version/kilianpaquier/pooling?style=for-the-badge">
  </a>
  <a href="https://score.getplumber.io/gitlab.com/kilianpaquier/pooling">
    <img alt="Plumber Score" src="https://img.shields.io/endpoint?url=https%3A%2F%2Fscore.getplumber.io%2Fgitlab.com%2Fkilianpaquier%2Fpooling.json&style=for-the-badge">
  </a>
</div>

---

- [How to use ?](#how-to-use-)
- [Documentation](#documentation)

## How to use ?

```sh
go get -u github.com/kilianpaquier/pooling@latest
```

## Documentation

Can be found here in a better format: <https://pkg.go.dev/github.com/kilianpaquier/pooling/pkg>.

```go
/*
Package pooling allows one to dispatch an infinite number of functions to be
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
```
