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

func NewMultiWatcher(input chan (<-chan interface{})) *MultiWatcher {
	ctx, cfunc := context.WithCancel(context.Background())
	return &MultiWatcher{
		input:     input,
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
	var otherChan <-chan interface{}

	go func() {
		var chanPtr <-chan interface{} = otherChan
		for {
			select {
			case data, ok := <-chanPtr:
				if ok {
					mw.resultant <- data
				} else {
					chanPtr = nil
				}
			case chnl := <-mw.input:
				if chanPtr == nil {
					chanPtr = chnl
				} else {
					newChan := make(chan interface{})
					var newReadChan <-chan interface{} = newChan
					oldChan := chanPtr
					chanPtr = newReadChan

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
