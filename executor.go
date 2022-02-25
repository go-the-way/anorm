// Copyright 2022 anox Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anox

import (
	"database/sql"
	"sync"
)

type executor struct {
	mu         *sync.Mutex
	db         *sql.DB
	tx         *sql.Tx
	openTX     bool
	autoCommit bool
}

func newExecutorFromOrm(o *orm) *executor {
	return newExecutor(o.db, o.openTX, o.autoCommit)
}

func newExecutor(db *sql.DB, openTX, autoCommit bool) *executor {
	etr := &executor{mu: &sync.Mutex{}, db: db, openTX: openTX, autoCommit: autoCommit}
	if openTX {
		etr.Begin()
	}
	return etr
}

func (e *executor) Begin() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.openTX && e.db != nil {
		tx, _ := e.db.Begin()
		e.tx = tx
	}
}

func (e *executor) Commit() error {
	if e.tx != nil {
		return e.tx.Commit()
	}
	return nil
}

func (e *executor) Rollback() error {
	if e.tx != nil {
		return e.tx.Rollback()
	}
	return nil
}

func (e *executor) exec(sqlStr string, ps ...interface{}) (sql.Result, error) {
	if e.tx != nil {
		if e.autoCommit {
			defer func() { _ = e.Commit() }()
		}
		return e.tx.Exec(sqlStr, ps...)
	}
	return e.db.Exec(sqlStr, ps...)
}

func (e *executor) query(sqlStr string, ps ...interface{}) (*sql.Rows, error) {
	if e.tx != nil {
		return e.tx.Query(sqlStr, ps...)
	}
	return e.db.Query(sqlStr, ps...)
}

func (e *executor) queryRow(sqlStr string, ps ...interface{}) *sql.Row {
	if e.tx != nil {
		return e.tx.QueryRow(sqlStr, ps...)
	}
	return e.db.QueryRow(sqlStr, ps...)
}
