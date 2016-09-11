package datatables

import (
	"errors"
	"net/url"

	"github.com/astaxie/beego/orm"
	"fmt"
)

type response struct {
	draws           int
	recordsTotal    int64
	recordsFiltered int64
	data            []interface{}
}

func (r *response) makeMap() map[string]interface{} {
	return map[string]interface{}{
		"draw":            r.draws,
		"recordsTotal":    r.recordsTotal,
		"recordsFiltered": r.recordsFiltered,
		"data":            r.data,
	}
}

type Data struct {
	UrlValues    url.Values //get args
	db           orm.Ormer
	TableName    string     //table name
	Columns      []string   //select column
}

func (d *Data) dbQuery(urlParams *urlParams, records []interface{}) (*response, error) {
	if urlParams == nil {
		return nil, errors.New("invalid url params argument")
	}

	// construct a string which is then passed to Select of orm.
	var selectStr string
	for k, v := range d.Columns {
		if k != 0 {
			selectStr += ","
		}
		selectStr += v
	}

	isSearch := urlParams.search != ""

	colOffset := urlParams.start
	qb, err := orm.NewQueryBuilder(d.db.Driver().Name())
	if err != nil {
		return nil, fmt.Errorf("failed to construct a new query builder. Error: %v", err)
	}

	qb.Select(selectStr).From(d.TableName)

	var whereStr string
	if isSearch {
		for k, v := range d.Columns {
			if k != 0 {
				whereStr += " OR "
			}
			whereStr += v + " LIKE " + "\"%" + urlParams.search + "%\"" //like
		}
		qb.Where(whereStr)

	}

	qb.OrderBy(d.Columns[urlParams.orderColumn])

	if urlParams.orderDir == "asc" {
		qb.Asc()
	} else {
		qb.Desc()
	}

	qb.Limit(urlParams.length).Offset(colOffset)

	fmt.Println(qb.String())

	_, err = d.db.Raw(qb.String()).QueryRows(records)
	if err != nil {
		return nil, fmt.Errorf("query failed. Error: %v", err)
	}

	recordsTotal, err := d.db.QueryTable(d.TableName).Count() //data sum

	var recordsFiltered int64 //search data sum

	if isSearch {
		qbFilterCount, _ := orm.NewQueryBuilder(d.db.Driver().Name())
		qbFilterCount.Select("Count(*) AS cnt ").From(d.TableName).Where(whereStr)
		sqlFilter := qbFilterCount.String()

		var rcount struct {
			Cnt int64
		}

		err = d.db.Raw(sqlFilter).QueryRow(&rcount)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch the number of filtered records. Error: %v", err)
		}

		recordsFiltered = rcount.Cnt
	} else {
		recordsFiltered = recordsTotal
	}

	return &response{
		draws:           urlParams.draws,
		recordsFiltered: recordsFiltered,
		recordsTotal:    recordsTotal,
		data:            records,
	}, nil
}

func (d *Data) Table(records []interface{}) (map[string]interface{}, error) {
	if d.db == nil {
		return nil, errors.New("invalid orm: nil value")
	}

	if records == nil || len(records) == 0 {
		return nil, errors.New("invalid resultColumns: nil or zero length slice")
	}

	urlParams, err := parseUrlQuery(d.UrlValues, len(d.Columns))
	if err != nil {
		return nil, err
	}

	response, err := d.dbQuery(urlParams, records)

	return response.makeMap(), err

}
