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
	"github.com/go-the-way/sg"
	"testing"
)

type (
	_Join struct {
		ID   int    `orm:"c{id} pk{T} def{id int comment 'ID'}"`
		Name string `orm:"def{name varchar(20) comment 'Name'}"`
	}
	_Entity struct {
		ID    int    `orm:"c{id} pk{T} def{id int comment 'ID'}"`
		Name  string `orm:"def{name varchar(20) comment 'Name'}"`
		Name3 string `orm:"c{name3} def{name3 varchar(20) comment 'Name3'}"`

		JoinID   int `orm:"c{join_id} def{join_id int}"`
		JoinName string

		Join2ID   int    `orm:"c{join_2_id} def{join_2_id int}"`
		Join2Name string `orm:"join{left,join_2_id,join,id}"`

		Join3ID   int    `orm:"c{join_3_id} def{join_3_id int}"`
		Join3Name string `orm:"join{left,join_3_id,join,id,name}"`
	}
	_ErrorEntity struct {
		//ID int `orm:"c{id} pk{T} def{id_int comment 'ID'}"`
	}
)

func (_ *_Join) Configure(c *EC) {
}

func (_ *_Entity) Configure(c *EC) {
	c.Migrate = true
	c.PrimaryKeyColumns = []sg.C{"id"}
	c.ColumnDefinitions = []sg.Ge{sg.ColumnDefinition(sg.C("name2"), sg.C("varchar(20)"), true, false, false, "", "Name2")}
	c.UpdateIgnores = []sg.C{"name3"}
	c.JoinRefMap = map[string]*JoinRef{
		"JoinName": {
			Field:      "JoinName",
			Type:       "left",
			SelfColumn: "join_id",
			RelTable:   "join",
			RelID:      "id",
			RelName:    "name",
		},
	}
}

func (_ *_ErrorEntity) Configure(c *EC) {
	c.Migrate = true
}

func TestEntity(t *testing.T) {
	Logger.SetLogLevel(LogLevelDebug)
	defer func() { _ = recover() }()
	Register(new(_Join))
	Register(new(_Entity))
	Register(new(_ErrorEntity))
}

func TestEntity2(t *testing.T) {
	defer func() { _ = recover() }()
	Logger.SetLogLevel(LogLevelDebug)
	Register(new(_ErrorEntity))
}
