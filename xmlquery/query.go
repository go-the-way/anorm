// Copyright 2022 anorm Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xmlquery

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-the-way/anorm"
	"github.com/go-the-way/anorm/pagination"
	"html/template"
)

var (
	ErrQueryTooManyResult = errors.New("xmlquery: query one return too many result")
)

type (
	Queryable[E anorm.Entity] interface {
		Query(entity E, ps ...any) ([]E, error)
		QueryOne(entity E, ps ...any) (E, error)
		QueryTemplate(entity E, data any) ([]E, error)
		QueryOneTemplate(entity E, data any) (E, error)
	}
	RawQueryable interface {
		Query(ps ...any) *sql.Row
		QueryRows(ps ...any) (*sql.Rows, error)
		QueryTemplate(data any) (*sql.Row, error)
		QueryRowsTemplate(data any) (*sql.Rows, error)
	}
	PageQueryable[E anorm.Entity] interface {
		Query(entity E, pager pagination.Pager, offset, size int, ps ...any) ([]E, int, error)
		QueryTemplate(entity E, pager pagination.Pager, offset, size int, data any) ([]E, int, error)
	}
	queryableImpl[E anorm.Entity] struct {
		ds     string
		db     *sql.DB
		sqlStr string
	}
	rawQueryableImpl struct {
		ds     string
		db     *sql.DB
		sqlStr string
	}
	pageQueryableImpl[E anorm.Entity] struct {
		ds     string
		db     *sql.DB
		sqlStr string
	}
)

func getDS(rn *rootNode, nd node) string {
	if ds := nd.getDatasource(); ds != "" {
		return ds
	}
	if ds := rn.Datasource; ds != "" {
		return ds
	}
	return ""
}

func Query[E anorm.Entity](namespace, id string) Queryable[E] {
	rn, nd := getNode(namespace, id, selectType)
	datasource := getDS(rn, nd)
	if datasource == "" {
		datasource = "_"
	}
	return &queryableImpl[E]{datasource, anorm.DataSourcePool.Required(datasource), nd.GetInnerXml()}
}

func (q *queryableImpl[E]) Query(entity E, ps ...any) ([]E, error) {
	queryLog("Query.Query", q.sqlStr, ps)
	rows, err := q.db.Query(q.sqlStr, ps...)
	if err != nil {
		queryErrorLog(err, "Query.QueryPage", q.sqlStr, ps)
		return nil, err
	}
	return anorm.ScanStruct(rows, entity, nil)
}

func (q *queryableImpl[E]) QueryOne(entity E, ps ...any) (e E, err error) {
	if es, err2 := q.Query(entity, q.sqlStr, ps); err2 != nil {
		err = err2
	} else if len(es) > 1 {
		err = ErrQueryTooManyResult
	} else if len(es) == 1 {
		e = es[0]
	}
	return
}

func (q *queryableImpl[E]) QueryTemplate(entity E, data any) ([]E, error) {
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, err
		}
		return q.Query(entity, buf.String())
	}
}

func (q *queryableImpl[E]) QueryOneTemplate(entity E, data any) (e E, err error) {
	if es, err2 := q.QueryTemplate(entity, data); err != nil {
		err = err2
	} else if len(es) > 1 {
		err = ErrQueryTooManyResult
	} else if len(es) == 1 {
		e = es[0]
	}
	return
}

func RawQuery(namespace, id string) RawQueryable {
	rn, nd := getNode(namespace, id, selectType)
	datasource := getDS(rn, nd)
	if datasource == "" {
		datasource = "_"
	}
	return &rawQueryableImpl{datasource, anorm.DataSourcePool.Required(datasource), nd.GetInnerXml()}
}

func (q *rawQueryableImpl) Query(ps ...any) *sql.Row {
	queryLog("RawQuery.Query", q.sqlStr, ps)
	return q.db.QueryRow(q.sqlStr, ps...)
}
func (q *rawQueryableImpl) QueryRows(ps ...any) (*sql.Rows, error) {
	queryLog("RawQuery.QueryRows", q.sqlStr, ps)
	return q.db.Query(q.sqlStr, ps...)
}

func (q *rawQueryableImpl) QueryTemplate(data any) (*sql.Row, error) {
	queryLog("RawQuery.QueryTemplate", q.sqlStr, nil)
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, err
		}
		return q.db.QueryRow(buf.String()), nil
	}
}

func (q *rawQueryableImpl) QueryRowsTemplate(data any) (*sql.Rows, error) {
	queryLog("RawQuery.QueryRowsTemplate", q.sqlStr, nil)
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, err
		}
		return q.db.Query(buf.String())
	}
}

func PageQuery[E anorm.Entity](namespace, id string) PageQueryable[E] {
	rn, nd := getNode(namespace, id, selectType)
	datasource := getDS(rn, nd)
	if datasource == "" {
		datasource = "_"
	}
	return &pageQueryableImpl[E]{datasource, anorm.DataSourcePool.Required(datasource), nd.GetInnerXml()}
}

func (q *pageQueryableImpl[E]) Query(entity E, pager pagination.Pager, offset, size int, ps ...any) ([]E, int, error) {
	pageSqlStr := fmt.Sprintf("select count(0) from (%s) as _t", q.sqlStr)
	queryLog("PageQuery.QueryCount", pageSqlStr, ps)
	c := 0
	row := q.db.QueryRow(pageSqlStr)
	if err := row.Scan(&c); err != nil {
		queryErrorLog(err, "PageQuery.QueryCount", pageSqlStr, ps)
		return nil, c, err
	}
	if c <= 0 {
		return []E{}, c, nil
	}
	sqlStr, ps := pager.Page(q.sqlStr, offset, size)
	newPs := make([]any, 0)
	newPs = append(newPs, ps...)
	rows, err := q.db.Query(sqlStr, newPs...)
	queryLog("PageQuery.Query", sqlStr, newPs)
	if err != nil {
		queryErrorLog(err, "PageQuery.Query", q.sqlStr, newPs)
		return nil, 0, err
	}
	es, err := anorm.ScanStruct(rows, entity, nil)
	return es, c, err
}

func (q *pageQueryableImpl[E]) QueryTemplate(entity E, pager pagination.Pager, offset, size int, data any) ([]E, int, error) {
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, 0, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, 0, err
		}
		return q.Query(entity, pager, offset, size)
	}
}
