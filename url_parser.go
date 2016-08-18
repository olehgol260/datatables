package datatables

import (
	"net/url"
	"strconv"
	"errors"
	"fmt"
)

const  (
	urlStart = "start"
	urlLength = "length"
	urlSearchValue = "search[value]"
	urlOrderColumn = "order[0][column]"
	urlOrderDir = "order[0][dir]"
	urlDraw = "draw"

	urlColBase = "columns[%d]"
	urlColData = "[data]"
	urlColName = "[name]"
	urlColSearchable = "[searchable]"
	urlColOrderable = "[orderable]"
	urlColSearchValue = "[search][value]"
	urlColSearchRegex = "[search][regex]"
)

type columnURLFields struct {
	data       string
	name       string
	searchable string
	orderable  string
	searchValue string
	searchRegex string
}

func newColumnUrlFields(colI int) (*columnURLFields) {
	columnBase := fmt.Sprintf(urlColBase, colI)
	col := new(columnURLFields)
	col.data = columnBase + urlColData
	col.name = columnBase + urlColName
	col.searchable = columnBase + urlColSearchable
	col.orderable = columnBase + urlColOrderable
	col.searchValue = columnBase + urlColSearchValue
	col.searchRegex = columnBase + urlColSearchRegex
	return col
}

type UrlParams struct {
	start          int
	length         int
	search         string
	orderColumn    int
	orderDirection string
	draws          int
	columns        []column
}

type column struct {
	data       int
	name       string
	searchable bool
	orderable  bool
	search     struct {
			   value string
			   regex bool
		   }
}

// ParseUrlQuery parse urlValues according to DataTables url parameter specification
func ParseUrlQuery(urlValues url.Values, columnsCount int) (*UrlParams, error) {
	if urlValues == nil {
		return nil, errors.New("urlValues argument equals to nil")
	}
	if columnsCount <= 0 {
		return nil, errors.New("columns count cannot be less than or equal to zero")
	}

	urlP := new(UrlParams)
	var err error

	urlP.start, err = strconv.Atoi(urlValues.Get(urlStart))
	if err != nil {
		return nil, errors.New("invalid url parameter: " + urlStart + " must be a valid integer")
	}

	urlP.length, err = strconv.Atoi(urlValues.Get(urlLength))
	if err != nil {
		return nil, errors.New("invalid url parameter: " + urlLength + " must be a valid integer")
	}

	urlP.search = urlValues.Get(urlSearchValue)

	urlP.orderColumn, err = strconv.Atoi(urlValues.Get(urlOrderColumn))
	if err != nil {
		return nil, errors.New("invalid url parameter: " + urlOrderColumn + " must be a valid integer")
	}

	urlP.orderDirection = urlValues.Get(urlOrderDir)
	if urlP.orderDirection != "asc" && urlP.orderDirection != "desc" {
		return nil, errors.New("invalid url parameter: " + urlOrderDir + " may be either 'asc' or 'desc'")
	}

	urlP.draws, err = strconv.Atoi(urlValues.Get(urlDraw))
	if err != nil {
		return nil, errors.New("invalid url parameter: " + urlDraw + " must be a valid integer")
	}

	urlP.columns = make([]column, columnsCount)

	for i := 0; i < columnsCount; i++ {
		colPtr := &urlP.columns[i]
		colFields := newColumnUrlFields(i)

		colPtr.data, err = strconv.Atoi(urlValues.Get(colFields.data))
		if err != nil {
			return nil, errors.New("invalid field " + colFields.data)
		}

		colPtr.name = urlValues.Get(colFields.name)

		colPtr.searchable, err = strconv.ParseBool(urlValues.Get(colFields.searchable))
		if err != nil {
			return nil, errors.New("failed to parse string field '" + colFields.searchable + "' to bool")
		}

		colPtr.orderable, err = strconv.ParseBool(urlValues.Get(colFields.orderable))
		if err != nil {
			return nil, errors.New("failed to parse string field '" + (colFields.orderable) + "' to bool")
		}

		colPtr.search.value = urlValues.Get(colFields.searchValue)

		colPtr.search.regex, err = strconv.ParseBool(urlValues.Get(colFields.searchRegex))
		if err != nil {
			return nil, errors.New("failed to parse string field '" + colFields.searchRegex + "' to bool")
		}
	}

	return urlP, nil
}