package nano

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/lonng/nano/session"
)

func TestChannel_Add(t *testing.T) {
	c := NewGroup("test_add")

	var paraCount = 100
	w := make(chan bool, paraCount)
	for i := 0; i < paraCount; i++ {
		go func(id int) {
			s := session.NewSession(nil)
			s.Bind(fmt.Sprintf("%d", id+1))
			c.Add(s)
			w <- true
		}(i)
	}

	for i := 0; i < paraCount; i++ {
		<-w
	}

	if c.Count() != paraCount {
		t.Fatalf("count expect: %d, got: %d", paraCount, c.Count())
	}

	n := fmt.Sprintf("%d", rand.Int63n(int64(paraCount))+1)
	if !c.Contains(n) {
		t.Fail()
	}

	// leave
	c.LeaveAll()
	if c.Count() != 0 {
		t.Fail()
	}
}
