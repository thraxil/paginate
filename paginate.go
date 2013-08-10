package paginate

import (
	"math"
	"net/http"
	"strconv"
)

func minint(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func maxint(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

type Paginator struct {
	ItemList Pagable
	PerPage  int
}

func (p Paginator) GetPage(r *http.Request) Page {
	pagen, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		// can't parse as int? just default to one
		pagen = 1
	}
	return Page{Paginator: p, PageNumber: pagen}
}

func (p Paginator) Count() int {
	return p.ItemList.TotalItems()
}

type Page struct {
	Paginator  Paginator
	PageNumber int
}

func (p Page) Items() []interface{} {
	return p.Paginator.ItemList.ItemRange(p.Offset(), p.NumItems())
}

func (p Page) Offset() int {
	total_items := p.Paginator.Count()
	offset := (p.PageNumber - 1) * p.Paginator.PerPage
	// a couple reasonable boundaries
	offset = minint(offset, total_items)
	offset = maxint(offset, 0)
	return offset
}

func (p Page) NumItems() int {
	total_items := p.Paginator.Count()
	cnt := p.Paginator.PerPage
	if p.Offset() >= (total_items - p.Paginator.PerPage) {
		cnt = total_items % p.Paginator.PerPage
	}
	return minint(p.Paginator.PerPage, cnt)
}

func (p Page) PrevPage() int {
	return maxint(p.PageNumber-1, 1)
}

func (p Page) HasPrev() bool {
	return p.PageNumber > 1
}

func (p Page) NextPage() int {
	total_items := p.Paginator.Count()
	return minint(p.PageNumber+1, int(total_items/p.Paginator.PerPage)+1)
}

func (p Page) HasNext() bool {
	total_items := p.Paginator.Count()
	return p.Offset() < (total_items - p.Paginator.PerPage)
}

type Pagable interface {
	TotalItems() int
	ItemRange(offset, count int) []interface{}
}
