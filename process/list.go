package process

import "container/list"

// List Process list container to aggregate information about working processes
type List struct {
	elemCurrent *list.Element
	list        *list.List
}

// Init creates a new process list
func (l *List) Init() {
	l.list = list.New()
}

// Add new process info into process list
func (l *List) Add(entry *Process) {
	l.list.PushBack(entry)
}

func (l *List) Read() *Process {
	return l.elemCurrent.Value.(*Process)
}

func (l *List) First() bool {
	if l.list.Front() == nil {
		return false
	}
	l.elemCurrent = l.list.Front()
	return true
}

func (l *List) Next() bool {
	if l.elemCurrent.Next() == nil {
		return false
	}
	l.elemCurrent = l.elemCurrent.Next()
	return true
}
