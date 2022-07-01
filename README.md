# anorm

```
 _____  ____    ___    ____  ____  
(____ ||  _ \  / _ \  / ___)|    \ 
/ ___ || | | || |_| || |    | | | |
\_____||_| |_| \___/ |_|    |_|_|_|

::anorm:: 

An another generic ORM framework implementation using the new way for Go.

{{ Version @VER }}

{{ Powered by go-the-way }}

{{ https://github.com/go-the-way/anorm }}

```

[![CircleCI](https://circleci.com/gh/go-the-way/anorm/tree/main.svg?style=shield)](https://circleci.com/gh/go-the-way/anorm/tree/main)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/go-the-way/anorm)
[![codecov](https://codecov.io/gh/go-the-way/anorm/branch/main/graph/badge.svg?token=8MAR3J959H)](https://codecov.io/gh/go-the-way/anorm)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-the-way/anorm)](https://goreportcard.com/report/github.com/go-the-way/anorm)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-the-way/anorm?status.svg)](https://pkg.go.dev/github.com/go-the-way/anorm?tab=doc)
[![Release](https://img.shields.io/github/release/go-the-way/anorm.svg?style=flat-square)](https://github.com/go-the-way/anorm/releases)


### Features
- DataSourcePool manage
- Pager implementation
- More levels logger
- Support joins
- Null fields
- sql.Null types
- Migrate DDL
- Column definition
- Index definition
- Simple Tx manager
- Insert batch
- Support XmlQuery

### Quickstart

```go
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
	count, err := o.OpsForSelectCount().Count(nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("count = ", count)
	}
}
```

### ORM Struct Tag

| TagName | PropertyName | Description            | Default      | Example                                              |
|---------|--------------|------------------------|--------------|------------------------------------------------------|
| pk      | PK           | Table primary key      | false        | pk{1},pk{t},pk{T},pk{true},pk{TRUE},pk{True}         |
| ig      | InsertIgnore | Ignore when inserts    | false        | ig{1},ig{t},ig{T},ig{true},ig{TRUE},ig{True}         |
| ug      | UpdateIgnore | Ignore when updates    | false        | ug{1},ug{t},ug{T},ug{true},ug{TRUE},ug{True}         |
| c       | Column       | Struct property column | propertyName | c{id},c{hello_world},c{halo_1234},c{WorldHa}         |
| def     | Definition   | Column definition SQL  |              | def{address varchar(100) not null comment 'Address'} |
| join    | JoinRef      | Join Ref definition    |              | join{inner,self_id,rel_table,rel_id,rel_name}        |
