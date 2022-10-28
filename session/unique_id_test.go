package session

import (
	"sync"
	"testing"
)

func TestUniqueSessionId(t *testing.T) {
	concurrent := 1000
	count := 100000
	m := sync.Map{}
	wg := sync.WaitGroup{}
	for concurrent > 0 {
		concurrent--
		wg.Add(1)
		go func(n int) {
			for n > 0 {
				n--
				id := UniqueSessionId()
				//fmt.Println(id)
				if _, ok := m.LoadOrStore(id, nil); ok {
					panic(id)
				}
				//time.Sleep(time.Microsecond)
			}
			wg.Done()
		}(count)
	}
	wg.Wait()
}
