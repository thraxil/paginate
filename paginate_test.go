package paginate

import (
	"testing"
)

// a boring test class
type items struct {
	Stuff []int
}

func (i items) TotalItems() int {
	return len(i.Stuff)
}

func (i items) ItemRange(offset, count int) []interface{} {
	out := make([]interface{}, len(i.Stuff))
	for j, v := range i.Stuff {
		out[j] = v
	}
	return out
}

func Test_Count(t *testing.T) {
	var i = items{Stuff: []int{1, 2, 3}}
	var p = Paginator{ItemList: i, PerPage: 20}
	if p.Count() != len(i.Stuff) {
		t.Error("wrong count")
	}
}
