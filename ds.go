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
)

type dsm map[string]*sql.DB

var (
	errRequiredNamedDSFunc = func(name string) error { return errors.New(fmt.Sprintf("required named[%s] DS", name)) }
	errDSIsNil             = errors.New("DS is nil")

	dsMap = make(dsm, 0)
)

func DS(db *sql.DB) {
	if db == nil {
		panic(errDSIsNil)
	}
	DSWithName("_", db)
	DSWithName("master", db)
}

func DSWithName(name string, db *sql.DB) {
	dsMap[name] = db
}

func (d *dsm) required(name string) *sql.DB {
	db := (*d)[name]
	if db == nil {
		panic(errRequiredNamedDSFunc(name))
	}
	return db
}
