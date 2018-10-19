package echo_utils

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// PaginationConfig defines the config for Pagination middleware.
	PaginationConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// PageParameter is a string to define a query parameter name used for pagination.
		// Default value "page".
		PageParameter            string
		PageParameterDefault     int
		PageParameterMax         int
		PageSizeParameter        string
		PageSizeParameterDefault int
		PageSizeParameterMin     int
		PageSizeParameterMax     int
		SortParameter            string
		FilterParameter          string
		PaginationContextKey     string
	}
)

type PageRequest struct {
	Page int
	Size int
	Sort string
}

type FilterRequest struct {
	Filters     []Filter
	PageRequest *PageRequest
}

func (f *FilterRequest) AddFilter(fs Filter) {
	f.Filters = append(f.Filters, fs)
}

func (f *FilterRequest) GetFilterByName(name string) *Filter {
	for _, filter := range f.Filters {
		if filter.Name == name {
			return &filter
		}
	}
	return nil
}

type Filter struct {
	Name        string
	FilterItems []map[string]interface{}
}

const (
	pageParameter            = "page"
	pageParameterDefault     = 1
	pageParameterMax         = 9999
	pageSizeParameter        = "pageSize"
	pageSizeParameterDefault = 20
	pageSizeParameterMin     = 2
	pageSizeParameterMax     = 1000
	sortParameter            = "sort"
	filterParameter          = "filter"
	paginationContextKey     = "pagination"
)

var (
	// DefaultPaginationConfig is the default Pagination middleware config.
	DefaultPaginationConfig = PaginationConfig{
		Skipper:                  middleware.DefaultSkipper,
		PageParameter:            pageParameter,
		PageParameterDefault:     pageParameterDefault,
		PageParameterMax:         pageParameterMax,
		PageSizeParameter:        pageSizeParameter,
		PageSizeParameterDefault: pageSizeParameterDefault,
		PageSizeParameterMin:     pageSizeParameterMin,
		PageSizeParameterMax:     pageSizeParameterMax,
		SortParameter:            sortParameter,
		FilterParameter:          filterParameter,
		PaginationContextKey:     paginationContextKey,
	}
)

// PaginationWithConfig returns an Pagination middleware with config.
func PaginationWithConfig(config PaginationConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultPaginationConfig.Skipper
	}
	if config.PageParameter == "" {
		config.PageParameter = pageParameter
	}
	if config.PageParameterDefault == 0 {
		config.PageParameterDefault = pageParameterDefault
	}
	if config.PageParameterMax == 0 {
		config.PageParameterMax = pageParameterMax
	}
	if config.PageSizeParameter == "" {
		config.PageSizeParameter = pageSizeParameter
	}
	if config.PageSizeParameterDefault == 0 {
		config.PageSizeParameterDefault = pageSizeParameterDefault
	}
	if config.PageSizeParameterMin == 0 {
		config.PageSizeParameterMin = pageSizeParameterMin
	}
	if config.PageSizeParameterMax == 0 {
		config.PageSizeParameterMax = pageSizeParameterMax
	}
	if config.SortParameter == "" {
		config.SortParameter = sortParameter
	}
	if config.FilterParameter == "" {
		config.FilterParameter = filterParameter
	}
	if config.PaginationContextKey == "" {
		config.PaginationContextKey = paginationContextKey
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			queryParams := c.QueryParams()
			pageParameter := queryParams.Get(config.PageParameter)
			pageSizeParameter := queryParams.Get(config.PageSizeParameter)
			sortParameter := queryParams.Get(config.SortParameter)

			var filters []Filter
			for queryParamName := range queryParams {
				if queryParamName == config.FilterParameter {
					filters = parseFilters(c.QueryParam(queryParamName))
				}
			}

			page := config.PageParameterDefault
			val, errPage := strconv.Atoi(pageParameter)
			if errPage == nil && val >= config.PageParameterDefault && val <= config.PageParameterMax {
				page = val
			}

			pageSize := config.PageSizeParameterDefault
			if val, errPageSize := strconv.Atoi(pageSizeParameter); errPageSize == nil {
				pageSize = val
				if val < config.PageSizeParameterMin {
					pageSize = config.PageSizeParameterMin
				}
				if val > config.PageSizeParameterMax {
					pageSize = config.PageSizeParameterMax
				}
			}

			filterRequest := &FilterRequest{
				Filters: filters,
				PageRequest: &PageRequest{
					Page: page,
					Size: pageSize,
					Sort: sortParameter,
				},
			}

			c.Set(config.PaginationContextKey, filterRequest)

			return next(c)
		}
	}
}

func parseFilters(filter string) []Filter {
	filters := make([]Filter, 0)
	var f []map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(filter), &f); err != nil {
		return filters
	}
	for _, filterItem := range f {
		for k, v := range filterItem {
			filters = append(filters, Filter{Name: k, FilterItems: v})
		}
	}
	return filters
}
