package future

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/goleak"
)

var sleepDur = time.Millisecond * 100

func TestSuccessful(t *testing.T) {
	defer goleak.VerifyNone(t)

	fut := InvokeAsync(func() (interface{}, error) {
		return "Hello, world!", nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := fut.Block(ctx)
	if err != nil {
		t.Errorf("expected future to return a successful result but got err=%+v", err)
	}

	expect := "Hello, world!"

	if res != "Hello, world!" {
		t.Errorf("expected result to be %s but got %s", expect, res)
	}
}

func TestFutureError(t *testing.T) {
	defer goleak.VerifyNone(t)

	fh := InvokeAsync(func() (interface{}, error) {
		return nil, errors.New("Some error")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := fh.Block(ctx)
	if err == nil {
		t.Errorf("expected Block(ctx) to return an error but was nil-")
	}
}

func TestCancelFuture(t *testing.T) {
	defer goleak.VerifyNone(t)

	fut := InvokeAsync(func() (interface{}, error) {
		return "Hello, world!", nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fut.Cancel()

	_, err := fut.Block(ctx)

	if err != ErrFutureCancelled {
		t.Errorf("expected Block to return ErrFutureCanceled after cancellation but didn't")
	}
}

func TestTimeoutFuture(t *testing.T) {
	defer func() {
		time.Sleep(sleepDur)
		goleak.VerifyNone(t)
	}()

	fut := InvokeAsync(func() (interface{}, error) {
		time.Sleep(sleepDur) // heavy computation
		return "Hello, world!", nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	_, err := fut.Block(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected err to be 'context deadline exceeded' but was %+v", err)
	}
}

func TestCancelWithContext(t *testing.T) {
	defer func() {
		time.Sleep(sleepDur)
		goleak.VerifyNone(t)
	}()
	fut := InvokeAsync(func() (interface{}, error) {
		time.Sleep(sleepDur)
		return "Hello, world!", nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := fut.Block(ctx)
	if err != ErrFutureCancelled {
		t.Errorf("expected context cancellation to cause ErrFutureCancelled but was: %+v", err)
	}
}
