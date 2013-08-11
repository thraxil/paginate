# paginate

Simple, basic pagination for Go.

The interface should be as familiar as possible for someone who's
worked with [Django's
pagination](https://docs.djangoproject.com/en/dev/topics/pagination/)
functionality.

## installing

The usual:

    $ go get github.com/thraxil/paginate

    import (
        "github.com/thraxil/paginate"
        ...
    )

## example usage

    // we'll just wrap a slice of ints
    type items struct {
        Stuff []int
    }
    // two methods that need to implemented
    // to satisify the paginate.Pagable interface
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
    // make a paginator. 20 items per page.
    var p = NewPaginator(itemset, 20)
    // get the first page of results
    page := p.GetPageNumber(1)
    // what you get from it:
    items_on_page := page.Items()
    prev := page.PrevPage()
    next := page.NextPage()

And so on. That's the basic idea.

You have some kind of type which represents a collection of items. You
need to implement an interface so the Paginator can get a count of the
total number of items and can ask for a contiguous slice of items from
your collection. Then it handles figuring out the offsets and
providing handy methods for (probably in a template) making next/prev
links and such.

The biggest difference from the Django version is because of Go's
strong typing and lack of generics. The .Items() method has to return
`[]interface{}` and you'll have to use type assertions to get your
item type back. Similarly, your .ItemRange() method has to convert
your items to `[]interface{}`. If anyone has a better idea for how
to do pagination on a generic collection in Go, I'm all ears.

The other difference is that Django's paginator handles some things by
raising exceptions. Eg, if you ask for a page that doesn't exist. Go,
of course doesn't have exceptions, and the idiomatic Go approach of
returning a (value, error) pair doesn't really work well if you are
working with the paginator in a template (which is a pretty key
use-case). I think the most reasonable thing for the Go version to do
is to try to always return something valid, and to provide functions
to validate the inputs. So, eg, `.PrevPage()` called on the first page
will just give you `1` instead of trying to send you to a `0` page. It
is your responsibility to make use of `.HasPrevPage()` to make the
interface clean for users.

See the [API
documentation](http://godoc.org/github.com/thraxil/paginate) for more
details.

Look in `paginate_test.go` for an example session translating the full
example code from the Django documentation to the Go equivalent.

Another small thing to note is that the `Paginator.GetPage()` function
takes an `http/request`, looks for a `page` parameter, and gives you
that page. I find that convenient.
