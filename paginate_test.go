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

	if p.NumPages() != 1 {
		t.Error("wrong number of pages")
	}

	var s = []int{1}
	if len(p.PageRange()) != len(s) {
		t.Error("wrong page range")
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

	if page.HasOtherPages() {
		t.Error("there should not be any other pages")
	}

	if page.StartIndex() != 1 {
		t.Error("wrong StartIndex")
	}

	if page.EndIndex() != 3 {
		t.Error("wrong EndIndex")
	}

}

type Beatles []string

func (b Beatles) TotalItems() int {
	return len(b)
}

func (b Beatles) ItemRange(offset, count int) []interface{} {
	out := make([]interface{}, len(b))
	for j, v := range b {
		out[j] = v
	}
	return out[offset : offset+count]
}

// duplicate the example code from Django's paginator
// and use that as a nice set of tests.
func Test_DjangoExamples(t *testing.T) {
	// >>> from django.core.paginator import Paginator
	// >>> objects = ['john', 'paul', 'george', 'ringo']
	var objects = Beatles{"john", "paul", "george", "ringo"}

	// >>> p = Paginator(objects, 2)
	p := NewPaginator(objects, 2)

	// >>> p.count
	// 4
	if p.Count() != 4 {
		t.Error("count is off")
	}

	// >>> p.num_pages
	// 2
	if p.NumPages() != 2 {
		t.Error("wrong number of pages")
	}

	// >>> p.page_range
	// [1, 2]
	if len(p.PageRange()) != 2 {
		t.Error("wrong page range")
	}

	// >>> page1 = p.page(1)
	// >>> page1
	// <Page 1 of 2>
	// >>> page1.object_list
	// ['john', 'paul']
	page1 := p.GetPageNumber(1)
	p1_items := page1.Items()
	if len(p1_items) != 2 {
		t.Error("wrong number of items on first page")
	}
	if p1_items[0].(string) != "john" {
		t.Error("not john")
	}
	if p1_items[1].(string) != "paul" {
		t.Error("not paul")
	}

	// >>> page2 = p.page(2)
	// >>> page2.object_list
	// ['george', 'ringo']
	page2 := p.GetPageNumber(2)
	p2_items := page2.Items()
	if len(p2_items) != 2 {
		t.Error("wrong number of items on second page")
	}
	if p2_items[0].(string) != "george" {
		t.Error("not george")
	}
	if p2_items[1].(string) != "ringo" {
		t.Error("not ringo")
	}

	// >>> page2.has_next()
	// False
	if page2.HasNext() {
		t.Error("page 2 is the last")
	}

	// >>> page2.has_previous()
	// True
	if !page2.HasPrev() {
		t.Error("there should be a previous page though")
	}

	// >>> page2.has_other_pages()
	// True
	if !page2.HasOtherPages() {
		t.Error("there should be other pages")
	}

	// >>> page2.next_page_number()
	// Traceback (most recent call last):
	// ...
	// EmptyPage: That page contains no results

	// instead of raising exceptions, let's just have the Go
	// version return the reasonable thing. Use .HasNext()
	// to test.
	if page2.NextPage() != 2 {
		t.Error("limit it to two pages")
	}

	// >>> page2.previous_page_number()
	// 1
	if page2.PrevPage() != 1 {
		t.Error("wrong prev page")
	}

	// >>> page2.start_index() # The 1-based index of the first item on this page
	// 3
	if page2.StartIndex() != 3 {
		t.Error("wrong start index")
	}

	// >>> page2.end_index() # The 1-based index of the last item on this page
	// 4
	if page2.EndIndex() != 4 {
		t.Error("wrong end index")
	}

	// haven't yet decided what we should do
	// when asked for invalid page numbers
	// django's approach is to raise exceptions:

	// >>> p.page(0)
	// Traceback (most recent call last):
	// ...
	// EmptyPage: That page number is less than 1

	page0 := p.GetPageNumber(0)
	if page0.PageNumber != 1 {
		t.Error("invalid page number")
	}

	// >>> p.page(3)
	// Traceback (most recent call last):
	// ...
	// EmptyPage: That page contains no results

	page3 := p.GetPageNumber(3)
	if page3.PageNumber != 2 {
		t.Error("invalid page number")
	}

}
