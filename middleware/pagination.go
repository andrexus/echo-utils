package echo_utils

import (
	"strconv"

	"strings"

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
		FilterParameterPrefix    string
		PaginationContextKey     string
	}
)

type PageRequest struct {
	Page int
	Size int
	Sort string
}

type FilterRequest struct {
	FilterName  string
	FilterValue string
	PageRequest *PageRequest
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
	filterParameterPrefix    = "filter."
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
		FilterParameterPrefix:    filterParameterPrefix,
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
	if config.FilterParameterPrefix == "" {
		config.FilterParameterPrefix = filterParameterPrefix
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

			filterName := ""
			filterValue := ""
			for queryParamName := range queryParams {
				queryParamValue := c.QueryParam(queryParamName)
				if strings.HasPrefix(queryParamName, config.FilterParameterPrefix) {
					filterName = queryParamName[len(config.FilterParameterPrefix):]
					filterValue = queryParamValue
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
				FilterName:  filterName,
				FilterValue: filterValue,
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
