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
	return out[offset : offset+count]
}

var itemset = items{Stuff: []int{1, 2, 3}}

func Test_Count(t *testing.T) {
	var p = NewPaginator(itemset, 20)
	if p.Count() != len(itemset.Stuff) {
		t.Error("wrong count")
	}
}

func Test_Page(t *testing.T) {
	var p = NewPaginator(itemset, 20)
	page := p.GetPageNumber(1)
	if len(page.Items()) != 3 {
		t.Error("not the right number of items returned")
	}

	if page.Offset() != 0 {
		t.Error("wrong offset")
	}

	if page.NumItems() != 3 {
		t.Error("NumItems is wrong")
	}

	if page.PrevPage() != 1 {
		t.Error("PrevPage wrong for first page")
	}

	if page.HasPrev() {
		t.Error("first page should not have a previous")
	}

	if page.NextPage() != 1 {
		t.Error("nextpage wrong on single page")
	}

	if page.HasNext() {
		t.Error("single page should not have a next")
	}

	if page.StartIndex() != 1 {
		t.Error("wrong StartIndex")
	}

	if page.EndIndex() != 3 {
		t.Error("wrong EndIndex")
	}


}
