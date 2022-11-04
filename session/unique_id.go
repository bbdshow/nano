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

// UniqueSessionId 可通过sessionId 得出创建时间与进程  不完全保证不重复，1ms内生成10000可保证不重复
func UniqueSessionId() int64 {
	i := atomic.AddInt64(&num, 1) % 10000
	milli := time.Now().UnixMilli()
	id := fmt.Sprintf("%d%02d%04d", milli, pid, i)
	v, _ := strconv.ParseInt(id, 10, 64)
	return v
}
