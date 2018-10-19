package echo_utils

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	q := req.URL.Query()
	filterValue := `[{"date": [{"gte": "2018-10-01T08:10:15.000Z", "lte":"2018-11-01T08:10:15.000Z"}]}, {"q": [{"name": "test"}]}]`
	q.Add("filter", filterValue)
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	cfg := DefaultPaginationConfig
	h := PaginationWithConfig(cfg)(func(c echo.Context) error {
		return nil
	})
	h(c)
	obj := c.Get(cfg.PaginationContextKey)
	assert.NotNil(t, obj)

	fr := obj.(*FilterRequest)
	assert.Equal(t, 2, len(fr.Filters))
	filter := fr.GetFilterByName("date")
	assert.NotNil(t, filter)
	assert.Equal(t, "date", fr.Filters[0].Name)

	input := map[string]interface{}{
		"gte": "2018-10-01T08:10:15.000Z",
		"lte": "2018-11-01T08:10:15.000Z",
	}
	fr.AddFilter(Filter{Name: "date", FilterItems: []map[string]interface{}{input}})

	assert.Equal(t, 3, len(fr.Filters))
}
