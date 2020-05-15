package multiwatch

import (
	"context"
)

type MultiWatcher struct {
	input     chan (<-chan interface{})
	resultant chan interface{}
	ctx       context.Context
	cancel    func()
}

func NewMultiWatcher() *MultiWatcher {
	ctx, cfunc := context.WithCancel(context.Background())
	return &MultiWatcher{
		input:     make(chan (<-chan interface{})),
		resultant: make(chan interface{}),
		ctx:       ctx,
		cancel:    cfunc,
	}
}

func (mw *MultiWatcher) Channel() <-chan interface{} {
	return mw.resultant
}

func (mw *MultiWatcher) EnqueueWatcher(chnl <-chan interface{}) {
	mw.input <- chnl
}

func (mw *MultiWatcher) Close() {
	mw.resultant = nil
}

func (mw *MultiWatcher) Stop() {
	mw.cancel()
}

func (mw *MultiWatcher) Watch() {
	go func() {
		var rootChan <-chan interface{}
		for {
			select {
			case data, ok := <-rootChan:
				if ok {
					mw.resultant <- data
				} else {
					rootChan = nil
				}
			case chnl := <-mw.input:
				if rootChan == nil {
					rootChan = chnl
				} else {
					newChan := make(chan interface{})
					oldChan := rootChan
					rootChan = newChan

					go func(ctx context.Context) {
						for {
							select {
							case d1, ok := <-oldChan:
								if ok {
									newChan <- d1
								} else {
									oldChan = nil
								}
							case d2, ok := <-chnl:
								if ok {
									newChan <- d2
								} else {
									chnl = nil
								}
							case <-ctx.Done():
								close(newChan)
								return
							}
						}
					}(mw.ctx)
				}
			case <-mw.ctx.Done():
				return
			}
		}
	}()
}
