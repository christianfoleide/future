#### An implementation of ``java.util.concurrent.Future`` in Go

[![Build Status](https://app.travis-ci.com/christianfoleide/future.svg?branch=master)](https://app.travis-ci.com/christianfoleide/future) [![codecov](https://codecov.io/gh/christianfoleide/future/branch/master/graph/badge.svg?token=Q0JNZRULUS)](https://codecov.io/gh/christianfoleide/future)

##### Installation
```
$ go get -u github.com/christianfoleide/future
```

##### Usage

```go

import "github.com/christianfoleide/future"

func main() {

      // Invoke and create handle
      fh := future.InvokeAsync(func() (interface{}, error) {
            // some computation
            return res, nil
      })

      // some other computation...

      ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 200)
      defer cancel()

      // Block(ctx) blocks until: 
      // 1. result is ready
      // 2. FutureFunc returns an error, 
      // 3. ctx receives on its Done channel
      result, err := fh.Block(ctx)
      if err != nil {
            // handle
      }
}
```

**Cancellation**
```go
fh.Cancel() // cancels the future
fh.IsCancelled() // true
fh.IsDone() // true
```

