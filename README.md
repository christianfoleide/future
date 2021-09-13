### An implementation of ``java.util.concurrent.Future`` in Go

[![Build Status](https://app.travis-ci.com/christianfoleide/future.svg?branch=master)](https://app.travis-ci.com/christianfoleide/future) [![codecov](https://codecov.io/gh/christianfoleide/future/branch/master/graph/badge.svg?token=Q0JNZRULUS)](https://codecov.io/gh/christianfoleide/future)

#### Usage

```go
func main() {

      // Invoke and create handle
      fh := future.InvokeAsync(func() (interface{}, error) {
            // some computation
            return res, nil
      })

      // some other computation...

      ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 200)
      defer cancel()

      // get the result
      res, err := fh.Block(ctx)
      if err != nil {
            // handle
      }
}
```