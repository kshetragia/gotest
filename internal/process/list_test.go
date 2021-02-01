package process

import (
	"testing"
)

func TestList(t *testing.T) {
	l := &List{}
	l.Init()

	l.Add(&Process{Name: "One"})
	l.Add(&Process{Name: "Two"})
	l.Add(&Process{Name: "Three"})

	testarr := [3]string{"One", "Two", "Three"}
	pos := 0
	for ok := l.First(); ok != false; ok = l.Next() {
		var data *Process
		if data = l.Read(); data == nil {
			t.Fatalf("Data not found")
		}
		if data.Name != testarr[pos] {
			t.Fatalf("Wrong Data in list: Got %v, Should be %v", data.Name, testarr[pos])
		}
		pos++
	}
}
