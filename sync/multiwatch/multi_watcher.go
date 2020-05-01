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

func (mw *MultiWatcher) FetchNext() <-chan interface{} {
	return mw.resultant
}

func (mw *MultiWatcher) QueueChannel(chnl <-chan interface{}) {
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
		var dataChan <-chan interface{}
		for {
			select {
			case data, ok := <-dataChan:
				if ok {
					mw.resultant <- data
				} else {
					dataChan = nil
				}
			case chnl := <-mw.input:
				if dataChan == nil {
					dataChan = chnl
				} else {
					newChan := make(chan interface{})
					var newReadChan <-chan interface{} = newChan
					oldChan := dataChan
					dataChan = newReadChan

					go func(ctx context.Context) {
						for {
							select {
							case d1, ok1 := <-oldChan:
								if ok1 {
									newChan <- d1
								} else {
									oldChan = nil
								}
							case d2, ok2 := <-chnl:
								if ok2 {
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
