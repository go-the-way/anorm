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
	"github.com/go-the-way/anorm"
	"text/template"
)

type (
	Executable interface {
		Exec(ps ...any) (int64, error)
		ExecTemplate(data any) (int64, error)
	}
	executableImpl struct {
		ds     string
		db     *sql.DB
		sqlStr string
	}
)

func executable(namespace, id string, nodeType nodeType) Executable {
	rn, nd := getNode(namespace, id, nodeType)
	datasource := getDS(rn, nd)
	if datasource == "" {
		datasource = "_"
	}
	return &executableImpl{datasource, anorm.DataSourcePool.Required(datasource), nd.GetInnerXml()}
}

func Insert(namespace, id string) Executable { return executable(namespace, id, insertType) } // Insert return Executable
func Delete(namespace, id string) Executable { return executable(namespace, id, deleteType) } // Delete return Executable
func Update(namespace, id string) Executable { return executable(namespace, id, updateType) } // Update return Executable

func (e *executableImpl) Exec(ps ...any) (int64, error) {
	var (
		result sql.Result
		err    error
	)
	queryLog("Executable.Exec", e.sqlStr, ps)
	result, err = e.db.Exec(e.sqlStr, ps...)
	queryErrorLog(err, "Executable.Exec", e.sqlStr, ps)
	if err != nil {
		return 0, err
	}
	if result != nil {
		return result.RowsAffected()
	}
	return 0, err
}
func (e *executableImpl) ExecTemplate(data any) (int64, error) {
	if temp, err := template.New("QUERY").Parse(e.sqlStr); err != nil {
		return 0, err
	} else {
		var buf = bytes.Buffer{}
		if err := temp.Execute(&buf, data); err != nil {
			return 0, err
		}
		return e.Exec()
	}
}
