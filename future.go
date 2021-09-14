package future

import (
	"context"
	"errors"
	"sync"
)

// FutureFunc is the function passed to the future.InvokeAsync method.
type FutureFunc func() (interface{}, error)

// ErrFutureCancelled is an error indicating that
// a future was cancelled. Calling future.Block() will return
// this error if the future was cancelled either by context cancellation,
// or by calling future.Cancel().
var ErrFutureCancelled = errors.New("future cancelled")

type futureHandle struct {
	fn         FutureFunc
	resChan    chan interface{}
	errChan    chan error
	cancelChan chan struct{}

	mu          sync.RWMutex
	isCancelled bool
	isDone      bool
}

// InvokeAsync invokes the given FutureFunc on a new goroutine, and immediately returns
// a handle to the result.
func InvokeAsync(fn FutureFunc) *futureHandle {
	fh := &futureHandle{
		fn:         fn,
		resChan:    make(chan interface{}, 1),
		errChan:    make(chan error, 1),
		cancelChan: make(chan struct{}, 1),
		mu:         sync.RWMutex{},
	}

	go func(h *futureHandle) {

		resPipe := make(chan interface{}, 1)
		errPipe := make(chan error, 1)

		defer func() {
			close(h.errChan)
			close(h.resChan)
		}()

		go func(rPipe chan interface{}, ePipe chan error) {

			res, err := h.fn()
			if err != nil {
				errPipe <- err
				close(errPipe)
				close(resPipe)
				return
			}
			resPipe <- res
			close(errPipe)
			close(resPipe)
		}(resPipe, errPipe)

		select {
		case <-h.cancelChan:
			h.setCancelled()
			h.errChan <- ErrFutureCancelled
		case err := <-errPipe:
			h.setDone()
			h.errChan <- err
		case res := <-resPipe:
			h.setDone()
			h.resChan <- res
		}

	}(fh)
	return fh
}

// Cancel cancels the receiving future, causing future.Block to return an error and also
// release resources.
func (h *futureHandle) Cancel() {
	h.setCancelled()
	h.cancelChan <- struct{}{}
}

// Block blocks until a result or an error returned by the receiving future is resolved, or
// returns an error if the given context receives on its Done channel.
func (h *futureHandle) Block(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		if ctx.Err() == context.Canceled {
			h.setCancelled()
			return nil, ErrFutureCancelled
		}
		return nil, ctx.Err() // context.DeadlineExceeded
	case err := <-h.errChan:
		return nil, err
	case res := <-h.resChan:
		return res, nil
	}
}

func (h *futureHandle) IsCancelled() bool {
	h.mu.RLocker().Lock()
	defer h.mu.RLocker().Unlock()
	return h.isCancelled
}

func (h *futureHandle) IsDone() bool {
	h.mu.RLocker().Lock()
	defer h.mu.RLocker().Unlock()
	return h.isDone
}

func (h *futureHandle) setCancelled() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.isCancelled = true
	h.isDone = true
}

func (h *futureHandle) setDone() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.isDone = true
}
