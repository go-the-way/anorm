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

// Package anorm
/*

An another generic ORM framework implementation using the new way for Go.

::: quickstart :::


package main

import (
	"database/sql"
	"fmt"

	a "github.com/go-the-way/anorm"
)

const (
	testDriverName = "mysql"
	testDSN        = "root:123456@tcp(localhost:3310)/test"
)

type Person struct {
	ID   int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
	Name string `orm:"pk{F} c{name} def{name varchar(20) not null default 'Coco' comment 'Name'}"`
}

func (p *Person) Configure(c *a.EC) {
	c.Migrate = true
	c.Table = "Table of Person"
}

func main() {
	db, err := sql.Open(testDriverName, testDSN)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.DataSourcePool.Push(db)
	a.Register(new(Person))
	o := a.New(new(Person))
	count, err := o.OpsForSelectCount().Query(nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("count = ", count)
	}
}
*/
package anorm
