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
	"fmt"
	"github.com/go-the-way/anorm"
	"github.com/go-the-way/anorm/pagination"
	"text/template"
)

type (
	PageSelectable[E anorm.Entity] interface {
		List(entity E, pager pagination.Pager, offset, size int, ps ...any) ([]E, int, error)
		ListTemplate(entity E, pager pagination.Pager, offset, size int, data any) ([]E, int, error)
	}
	pageSelectableImpl[E anorm.Entity] struct {
		ds              string
		db              *sql.DB
		sqlStr, tSqlStr string
	}
)

func PageSelect[E anorm.Entity](namespace, id string) PageSelectable[E] {
	datasource, db, sqlStr := getNodeParams(namespace, id, selectType)
	return &pageSelectableImpl[E]{datasource, db, sqlStr, ""}
}

func (q *pageSelectableImpl[E]) List(entity E, pager pagination.Pager, offset, size int, ps ...any) ([]E, int, error) {
	sqlStr := q.tSqlStr
	if sqlStr == "" {
		sqlStr = q.sqlStr
	}
	pageSqlStr := fmt.Sprintf("select count(0) from (%s) as _t", sqlStr)
	queryLog("PageSelect.SelectCount", pageSqlStr, ps...)
	c := 0
	row := q.db.QueryRow(pageSqlStr)
	if err := row.Scan(&c); err != nil {
		queryErrorLog(err, "PageSelect.SelectCount", pageSqlStr, ps...)
		return nil, c, err
	}
	if c <= 0 {
		return []E{}, c, nil
	}
	nSqlStr, ps2 := pager.Page(sqlStr, offset, size)
	newPs := make([]any, 0)
	newPs = append(newPs, ps...)
	newPs = append(newPs, ps2...)
	rows, err := q.db.Query(nSqlStr, newPs...)
	queryLog("PageSelect.List", nSqlStr, newPs...)
	if err != nil {
		queryErrorLog(err, "PageSelect.List", nSqlStr, newPs...)
		return nil, 0, err
	}
	es, err := anorm.ScanStruct(rows, entity, nil)
	return es, c, err
}

func (q *pageSelectableImpl[E]) ListTemplate(entity E, pager pagination.Pager, offset, size int, data any) ([]E, int, error) {
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, 0, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, 0, err
		}
		q.tSqlStr = buf.String()
		return q.List(entity, pager, offset, size)
	}
}
