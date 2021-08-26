package notify

import (
	"testing"
	"time"
)

func TestTickerChannel(t *testing.T) {
	tc := NewTickerChannel(100, 5*time.Second)
	go func() {
		i := 1
		for {
			tc.Store(i)
			i++
			time.Sleep(time.Second)
		}
	}()

	tc.Range(func(data interface{}){
		t.Log(data)
	})
}
