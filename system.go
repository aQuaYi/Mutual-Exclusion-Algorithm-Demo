package mutual

import (
	"time"
)

type system struct {
	processes []*process
}

// size: process 的数量
func newSystem(size int) *system {
	chans := make([]chan *message, size)
	for i := range chans {
		// TODO: chan 可以带缓冲吗？
		chans[i] = make(chan *message)
	}

	ps := make([]*process, size)
	for i := range ps {
		ps[i] = newProcess(i, chans)
	}

	return &system{
		processes: ps,
	}
}

func (s *system) kill() {

	// TODO: 添加退出机制

	time.Sleep(20 * time.Second)
}