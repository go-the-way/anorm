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

An another ORM framework implementation using the new way for Go.

::: quickstart :::

package main

import (
  "database/sql"
  "fmt"
  . "github.com/go-the-way/anorm"
  "github.com/go-the-way/sg"
)

const (
  testDriverName = "mysql"
  testDSN        = "root:123456@tcp(localhost:3306)/test"
)

var (
  testDB, _ = sql.Open(testDriverName, testDSN)
)

type userModel struct {
  ID         int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
  Name       string `orm:"pk{F} c{name} def{name varchar(50) not null default 'hello world' comment 'Name'}"`
  Age        int    `orm:"pk{F} c{age} def{age int not null default '20' comment 'Age'}"`
  Address    string `orm:"pk{F} c{address} def{address varchar(100) not null comment 'Address'}"`
  Phone      string `orm:"pk{F} c{phone} def{phone varchar(11) not null default '13900000000' comment 'Phone'}"`
  CreateTime string `orm:"pk{F} c{create_time} ig{T} ug{T} def{create_time datetime not null default current_timestamp comment 'CreateTime'}"`
  XYZ        string `orm:"pk{F} c{xyz} def{xyz varchar(50) not null default 'xyz' comment 'XYZ'}"`
}

func (u *userModel) MetaData() *ModelMeta {
  return &ModelMeta{
    Migrate: true,
    Comment: "The userModel Table",
    IndexDefinitions: []sg.Ge{
      sg.IndexDefinition(false, sg.P("idx_name"), sg.C("name")),
    },
    InsertIgnores: []sg.C{"id", "create_time"},
  }
}

func main() {
  DS(testDB)
  Register(new(userModel))
  o := New(new(userModel))
  count, err := o.Select().Count(&userModel{ID: 1})
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("count = ", count)
  }
}
*/
package anorm