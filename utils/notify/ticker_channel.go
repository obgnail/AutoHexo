package notify

import "time"

type TickerChannel struct {
	size       int
	tickerTime time.Duration

	inputChan  chan interface{}
	outputChan chan interface{}
	tickerChan chan struct{}
}

func NewTickerChannel(size int, tickerTime time.Duration) *TickerChannel {
	cc := &TickerChannel{
		inputChan:  make(chan interface{}, size),
		outputChan: make(chan interface{}, 1),
		tickerChan: make(chan struct{}, 1),
		tickerTime: tickerTime,
	}

	go cc.output()
	go cc.ticker()
	return cc
}

func (tc *TickerChannel) Store(data interface{}) {
	tc.inputChan <- data
}

func (tc *TickerChannel) Range(f func(data interface{})) {
	for d := range tc.outputChan {
		f(d)
	}
}

func (tc *TickerChannel) ticker() {
	ticker := time.NewTicker(tc.tickerTime)
	go func() {
		for range ticker.C {
			tc.tickerChan <- struct{}{}
		}
	}()
}

func (tc *TickerChannel) output() {
	var dataset []interface{}
	for {
		select {
		case d := <-tc.inputChan:
			dataset = append(dataset, d)
		case <-tc.tickerChan:
			for _, d := range dataset {
				tc.outputChan <- d
			}
			dataset = []interface{}{}
		}
	}
}
