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
	"text/template"
)

type (
	RowsSelectable interface {
		Row(ps ...any) *sql.Row
		Rows(ps ...any) (*sql.Rows, error)
		RowTemplate(data any) (*sql.Row, error)
		RowsTemplate(data any) (*sql.Rows, error)
	}
	rowsSelectableImpl struct {
		ds              string
		db              *sql.DB
		sqlStr, tSqlStr string
	}
)

func RowsSelect(namespace, id string) RowsSelectable {
	datasource, db, sqlStr := getNodeParams(namespace, id, selectType)
	return &rowsSelectableImpl{datasource, db, sqlStr, ""}
}

func (q *rowsSelectableImpl) Row(ps ...any) *sql.Row {
	sqlStr := q.tSqlStr
	if sqlStr == "" {
		sqlStr = q.sqlStr
	}
	queryLog("RowsSelectable.Row", sqlStr, ps...)
	return q.db.QueryRow(sqlStr, ps...)
}

func (q *rowsSelectableImpl) Rows(ps ...any) (*sql.Rows, error) {
	sqlStr := q.tSqlStr
	if sqlStr == "" {
		sqlStr = q.sqlStr
	}
	queryLog("RowsSelectable.Rows", sqlStr, ps...)
	return q.db.Query(sqlStr, ps...)
}

func (q *rowsSelectableImpl) RowTemplate(data any) (*sql.Row, error) {
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, err
		} else {
			q.tSqlStr = buf.String()
		}
		return q.Row(), nil
	}
}

func (q *rowsSelectableImpl) RowsTemplate(data any) (*sql.Rows, error) {
	if temp, err := template.New("QUERY").Parse(q.sqlStr); err != nil {
		return nil, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return nil, err
		} else {
			q.tSqlStr = buf.String()
		}
		return q.Rows()
	}
}
