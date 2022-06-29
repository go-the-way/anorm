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

package anorm

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

type dataSourcePool struct {
	mu  *sync.Mutex
	dbM map[string]*sql.DB
}

var (
	// DataSourcePool global datasource pool
	DataSourcePool = &dataSourcePool{mu: &sync.Mutex{}, dbM: make(map[string]*sql.DB, 0)}
)

// Push master datasource
func (p *dataSourcePool) Push(db *sql.DB) {
	p.PushDB("_", db)
	p.PushDB("master", db)
}

// PushDB push name datasource
func (p *dataSourcePool) PushDB(name string, db *sql.DB) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dbM[name] = db
}

// Required required name datasource
func (p *dataSourcePool) Required(name string) *sql.DB {
	db, have := p.dbM[name]
	if !have {
		panic(errors.New(fmt.Sprintf("anorm: required named[%s] data source", name)))
	}
	return db
}
