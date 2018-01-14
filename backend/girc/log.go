package girc

import (
	"sync"

	"github.com/cceckman/discoirc/data"
)

type Log struct {
	data.Scope

	data.EventList
	sync.RWMutex
}

func (l *Log) Append(e data.Event) {
	n := e
	n.Scope = l.Scope

	l.Lock()
	defer l.Unlock()


	if len(l.EventList) != 0 {
		n.Seq = l.EventList[l.EventList.Len()].Seq + 1
	}

	l.EventList = append(l.EventList, n)
}

func (l *Log) EventsBefore(n int, max data.Seq) data.EventList {
	l.RLock()
	defer l.RUnlock()

	return l.EventList.SelectSizeMax(n, max)
}
