package notify

import "time"

type TickerChannel struct {
	size       int
	tickerTime time.Duration

	inputChan  chan interface{}
	outputChan chan interface{}
	sleepChan  chan struct{}
}

func NewTickerChannel(size int, tickerTime time.Duration) *TickerChannel {
	cc := &TickerChannel{
		inputChan:  make(chan interface{}, size),
		outputChan: make(chan interface{}, 1),
		sleepChan:  make(chan struct{}, 1),
		tickerTime: tickerTime,
	}

	go cc.ticker()
	go cc.output()
	return cc
}

func (tc *TickerChannel) Store(data interface{}) {
	tc.inputChan <- data
}

func (tc *TickerChannel) Range(f func(data interface{})) {
	for data := range tc.outputChan {
		f(data)
	}
}

func (tc *TickerChannel) ticker() {
	ticker := time.NewTicker(tc.tickerTime)
	go func() {
		for range ticker.C {
			tc.sleepChan <- struct{}{}
		}
	}()
}

func (tc *TickerChannel) output() {
	var dataset []interface{}
	for {
		select {
		case data := <-tc.inputChan:
			dataset = append(dataset, data)
		case <-tc.sleepChan:
			for _, d := range dataset {
				tc.outputChan <- d
			}
			dataset = []interface{}{}
			continue
		}
	}
}
