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
	"errors"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-the-way/sg"
)

var (
	testDB                *sql.DB
	testDriverName        = "mysql"
	errTestDSNNotProvided = errors.New("anox.test: the test DSN not provided")
)

type (
	userModel struct {
		ID         int       `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name       string    `orm:"pk{F} c{name} def{name varchar(50) not null default 'hello world' comment 'Name'}"`
		Age        int       `orm:"pk{F} c{age} def{age int not null default '20' comment 'Age'}"`
		Address    string    `orm:"pk{F} c{address} def{address varchar(100) not null comment 'Address'}"`
		Phone      string    `orm:"pk{F} c{phone} def{phone varchar(11) not null default '13900000000' comment 'Phone'}"`
		CreateTime time.Time `orm:"pk{F} c{create_time} ig{T} def{create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'CreateTime'}" json:"create_time"`
	}
	userModelNull struct {
		ID         int            `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name       sql.NullString `orm:"pk{F} c{name} def{name varchar(50) null comment 'Name'}"`
		Age        sql.NullInt32  `orm:"pk{F} c{age} def{age int null comment 'Age'}"`
		Address    sql.NullString `orm:"pk{F} c{address} def{address varchar(100) null comment 'Address'}"`
		Phone      sql.NullString `orm:"pk{F} c{phone} def{phone varchar(11) null comment 'Phone'}"`
		CreateTime sql.NullTime   `orm:"pk{F} c{create_time} ig{T} def{create_time datetime null comment 'CreateTime'}" json:"create_time"`
	}
)

func (u *userModel) MetaData() *ModelMeta {
	return &ModelMeta{
		Migrate:           true,
		Comment:           "The userModel Table",
		ColumnDefinitions: []sg.Ge{},
		IndexDefinitions:  []sg.Ge{sg.IndexDefinition(false, sg.P("idx_name"), sg.C("name"))},
		InsertIgnores:     []sg.C{"id", "create_time"},
	}
}

func (u *userModelNull) MetaData() *ModelMeta {
	return &ModelMeta{
		Migrate:           true,
		Comment:           "The userModelNull Table",
		ColumnDefinitions: []sg.Ge{},
		IndexDefinitions:  []sg.Ge{},
		InsertIgnores:     []sg.C{},
	}
}

var envInit = false

func testInit() {
	if envInit {
		return
	}
	if !envInit {
		envInit = true
	}
	if dsn := os.Getenv("ANOX_TEST_DSN"); dsn == "" {
		panic(errTestDSNNotProvided)
	} else {
		if tdb, err := sql.Open(testDriverName, dsn); err != nil {
			panic(err)
		} else {
			testDB = tdb
			DS(tdb)
		}
	}
	Register(new(userModel))
	Register(new(userModelNull))
}

func insertUserModel(name string, age int, address, phone string) error {
	_, err := testDB.Exec("insert into user_model(name,age,address,phone) values (?,?,?,?)", name, age, address, phone)
	return err
}

func selectUserModelCount(id int) int {
	c := 0
	testDB.QueryRow("select count(0) from user_model where id = ?", id).Scan(&c)
	return c
}

func truncateTestTable() {
	_, _ = testDB.Exec("truncate table user_model")
	_, _ = testDB.Exec("truncate table user_model_null")
}
