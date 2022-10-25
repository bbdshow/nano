package session

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var pid = os.Getpid() % 100
var num int64

func UniqueSessionId() int64 {
	i := atomic.AddInt64(&num, 1) % 10000
	milli := time.Now().UnixMilli()
	id := fmt.Sprintf("%d%02d%04d", milli, pid, i)
	v, _ := strconv.ParseInt(id, 10, 64)
	return v
}
