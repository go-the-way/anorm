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
	SingleSelectable[T any] interface {
		One(ps ...any) (T, error)
		OneTemplate(data any) (T, error)
	}
	singleSelectableImpl[T any] struct {
		ds              string
		db              *sql.DB
		sqlStr, tSqlStr string
	}
)

func SingleSelect[T any](namespace, id string) SingleSelectable[T] {
	datasource, db, sqlStr := getNodeParams(namespace, id, selectType)
	return &singleSelectableImpl[T]{datasource, db, sqlStr, ""}
}

func (s *singleSelectableImpl[T]) One(ps ...any) (t T, err error) {
	sqlStr := s.tSqlStr
	if sqlStr == "" {
		sqlStr = s.sqlStr
	}
	queryLog("SingleSelectable.One", sqlStr, ps...)
	row := s.db.QueryRow(s.sqlStr, ps...)
	if err2 := row.Err(); err2 != nil {
		err = err2
		return
	}
	if err2 := row.Scan(&t); err2 != nil {
		err = err2
	}
	return
}

func (s *singleSelectableImpl[T]) OneTemplate(data any) (t T, err error) {
	if temp, err2 := template.New("QUERY").Parse(s.sqlStr); err2 != nil {
		err = err2
	} else {
		var buf = bytes.Buffer{}
		if err2 := temp.Execute(&buf, data); err2 != nil {
			err = err2
		}
		s.tSqlStr = buf.String()
		t, err = s.One()
	}
	return
}
