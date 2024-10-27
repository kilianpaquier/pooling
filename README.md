<!-- This file is safe to edit. Once it exists it will not be overwritten. -->

# pooling <!-- omit in toc -->

<p align="center">
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/kilianpaquier/pooling?include_prereleases&sort=semver&style=for-the-badge">
  <img alt="GitHub Issues" src="https://img.shields.io/github/issues-raw/kilianpaquier/pooling?style=for-the-badge">
  <img alt="GitHub License" src="https://img.shields.io/github/license/kilianpaquier/pooling?style=for-the-badge">
  <img alt="Coverage" src="https://img.shields.io/codecov/c/github/kilianpaquier/pooling/main?style=for-the-badge">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/kilianpaquier/pooling/main?style=for-the-badge&label=Go+Version">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kilianpaquier/pooling?style=for-the-badge">
</p>

---

- [How to use ?](#how-to-use-)
- [Features](#features)

## How to use ?

```sh
go get -u github.com/kilianpaquier/pooling@latest
```

## Features

The pooling package allows one to dispatch an infinite number of functions to be executed in parallel while still limiting the number of routines.

For that, pooling package takes advantage of ants pool library. A pooling Pooler can have multiple pools (with builder SetSizes) to dispatch sub functions into different pools of routines.

When sending a function into the pooler (with the appropriate channel), this function can itself send other functions into the pooler. It allows one to "split" functions executions (like iterating over a slice and each element handled in parallel).

```go
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
```
