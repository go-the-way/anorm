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
	"github.com/go-the-way/anorm"
	"html/template"
)

var (
	ErrSelectTooManyResult = errors.New("xmlquery: query one return too many result")
)

type (
	Selectable[E anorm.Entity] interface {
		List(entity E, ps ...any) ([]E, error)
		One(entity E, ps ...any) (E, error)
		ListTemplate(entity E, data any) ([]E, error)
		OneTemplate(entity E, data any) (E, error)
	}
	selectableImpl[E anorm.Entity] struct {
		ds              string
		db              *sql.DB
		sqlStr, tSqlStr string
	}
)

func Select[E anorm.Entity](namespace, id string) Selectable[E] {
	datasource, db, sqlStr := getNodeParams(namespace, id, selectType)
	return &selectableImpl[E]{datasource, db, sqlStr, ""}
}

func (q *selectableImpl[E]) List(e E, ps ...any) ([]E, error) {
	sqlStr := q.tSqlStr
	if sqlStr == "" {
		sqlStr = q.sqlStr
	}
	queryLog("Selectable.List", sqlStr, ps...)
	rows, err := q.db.Query(sqlStr, ps...)
	if err != nil {
		queryErrorLog(err, "Selectable.List", sqlStr, ps...)
		return nil, err
	}
	return anorm.ScanStruct(rows, e, nil)
}

func (q *selectableImpl[E]) One(ie E, ps ...any) (oe E, err error) {
	if es, err2 := q.List(ie, ps...); err2 != nil {
		err = err2
	} else if len(es) > 1 {
		err = ErrSelectTooManyResult
	} else if len(es) == 1 {
		oe = es[0]
	}
	return
}

func (q *selectableImpl[E]) ListTemplate(e E, data any) ([]E, error) {
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, err
		} else {
			q.tSqlStr = buf.String()
		}
		return q.List(e)
	}
}

func (q *selectableImpl[E]) OneTemplate(ie E, data any) (oe E, err error) {
	if es, err2 := q.ListTemplate(ie, data); err2 != nil {
		err = err2
	} else if len(es) > 1 {
		err = ErrSelectTooManyResult
	} else if len(es) == 1 {
		oe = es[0]
	}
	return
}
