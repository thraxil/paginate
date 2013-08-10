package paginate

import (
	"math"
	"net/http"
	"strconv"
)

type Paginator struct {
	ItemList Pagable
	PerPage  int
}

// make a new Paginator from a collection of items and number of
// items per page
func NewPaginator(itemset Pagable, pp int) *Paginator {
	return &Paginator{ItemList: itemset, PerPage: pp}
}

// get a Page from an HTTP request. For now, it's hard
// coded to expect a 'page' parameter
func (p Paginator) GetPage(r *http.Request) Page {
	pagen, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		// can't parse as int? just default to one
		pagen = 1
	}
	return p.GetPageNumber(pagen)
}

// get a Page by number
func (p Paginator) GetPageNumber(n int) Page {
	return Page{Paginator: p, PageNumber: n}
}

// total number of items
func (p Paginator) Count() int {
	return p.ItemList.TotalItems()
}

type Page struct {
	Paginator  Paginator
	PageNumber int
}

// return the items on the page. Returns them as a slice
// of interface{}, so you'll need to cast them back
func (p Page) Items() []interface{} {
	return p.Paginator.ItemList.ItemRange(p.Offset(), p.NumItems())
}

// starting offset for items on the Page
func (p Page) Offset() int {
	total_items := p.Paginator.Count()
	offset := (p.PageNumber - 1) * p.Paginator.PerPage
	// a couple reasonable boundaries
	offset = minint(offset, total_items)
	offset = maxint(offset, 0)
	return offset
}

// Returns the 1-based index of the first object on the page,
// relative to all of the objects in the paginator’s list.
// For example, when paginating a list of 5 objects with 2
// objects per page, the second page’s start_index() would return 3.
func (p Page) StartIndex() int {
	return p.Offset() + 1
}

// number of items on the Page
func (p Page) NumItems() int {
	total_items := p.Paginator.Count()
	cnt := p.Paginator.PerPage
	if p.Offset() >= (total_items - p.Paginator.PerPage) {
		cnt = total_items % p.Paginator.PerPage
	}
	return minint(p.Paginator.PerPage, cnt)
}

// page number for the page before this
// bottoms out at the first page
func (p Page) PrevPage() int {
	return maxint(p.PageNumber-1, 1)
}

// does this Page have one before it?
func (p Page) HasPrev() bool {
	return p.PageNumber > 1
}

// page number for the next page. won't go past the end
func (p Page) NextPage() int {
	total_items := p.Paginator.Count()
	return minint(p.PageNumber+1, int(total_items/p.Paginator.PerPage)+1)
}

// is there a page after this one?
func (p Page) HasNext() bool {
	total_items := p.Paginator.Count()
	return p.Offset() < (total_items - p.Paginator.PerPage)
}

type Pagable interface {
	TotalItems() int
	ItemRange(offset, count int) []interface{}
}

func minint(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func maxint(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}
