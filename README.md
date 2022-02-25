# anox

```
 _____  ____    ___    ___    ____  ____  
(____ ||  _ \  / _ \  / _ \  / ___)|    \ 
/ ___ || | | || |_| || |_| || |    | | | |
\_____||_| |_| \___/  \___/ |_|    |_|_|_|

::anox:: 

An another ORM framework implementation using the new way for Go.

{{ Version @VER }}

{{ Powered by go-the-way }}

{{ https://github.com/go-the-way/anox }}

```

[![CircleCI](https://circleci.com/gh/go-the-way/anox/tree/main.svg?style=shield)](https://circleci.com/gh/go-the-way/anox/tree/main)
[![codecov](https://codecov.io/gh/go-the-way/anox/branch/main/graph/badge.svg?token=8MAR3J959H)](https://codecov.io/gh/go-the-way/anox)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-the-way/anox)](https://goreportcard.com/report/github.com/go-the-way/anox)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-the-way/anox?status.svg)](https://pkg.go.dev/github.com/go-the-way/anox?tab=doc)
[![Release](https://img.shields.io/github/release/go-the-way/anox.svg?style=flat-square)](https://github.com/go-the-way/anox/releases)
[![License](https://img.shields.io/github/license/go-the-way/anox)](https://raw.githubusercontent.com/go-the-way/anox/main/LICENSE)

### Quickstart
```go
package main

import (
  "database/sql"
  "fmt"
  
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-the-way/anox"
  "github.com/go-the-way/sg"
)

const (
  testDriverName = "mysql"
  testDSN        = "root:123456@tcp(127.0.0.1:3306)/test"
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

```

# Docs
- [wiki](https://github.com/go-the-way/anox/wikis)