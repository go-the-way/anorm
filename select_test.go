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
	"errors"
	"github.com/go-the-way/anorm/pagination"
	"github.com/go-the-way/sg"
	"testing"
)

func init() {
	testInit()
}

func TestSelectExec(t *testing.T) {
	truncateTestTable()
	if err := insertTest(); err != nil {
		t.Fatalf("TestSelectExec failed: %v\n", err)
	}
	o := New(new(userEntity))
	{
		if entities, err := o.OpsForSelect().OrderBy(sg.C("id")).List(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().IfWhere(false).IfWhere(true, getTestGes()...).List(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where(getTestGes()...).List(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where().List(getTest()); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectExec failed!")
		}
	}
}

func TestSelectExecOne(t *testing.T) {
	truncateTestTable()
	if err := insertTest(); err != nil {
		t.Fatalf("TestSelectExecOne failed: %v\n", err)
	}
	o := New(new(userEntity))
	{
		if ee, err := o.OpsForSelect().One(nil); err != nil {
			t.Fatalf("TestSelectExecOne failed: %v\n", err)
		} else if ee == nil {
			t.Fatal("TestSelectExecOne failed!")
		}
	}
}

func TestSelectExecOne2(t *testing.T) {
	truncateTestTable()
	o := New(new(userEntity))
	{
		if ee, err := o.OpsForSelect().One(nil); err != nil {
			t.Fatalf("TestSelectExecOne failed: %v\n", err)
		} else if ee != nil {
			t.Fatal("TestSelectExecOne failed!")
		}
	}
}

func TestSelectExecPageError(t *testing.T) {
	{
		truncateTestTable()
		o := New(new(userEntity))
		{
			if _, _, err := o.OpsForSelect().Where(sg.And(sg.Eq("1", 1))).Page(nil, pagination.Pg, 0, 2); err == nil {
				t.Fatalf("TestSelectExecPageError failed: %v\n", err)
			}
		}
	}
	{
		o := New(new(userEntity))
		{
			if entities, count, err := o.OpsForSelect().Page(nil, pagination.MySql, 0, 2); err != nil {
				t.Fatalf("TestSelectExecPageError failed: %v\n", err)
			} else if len(entities) != 0 && int(count) != 0 {
				t.Fatal("TestSelectExecPageError failed!")
			}
		}
	}

}

func TestSelectExecPage(t *testing.T) {
	truncateTestTable()
	c := 10
	for i := 0; i < c; i++ {
		_ = insertTest()
	}
	o := New(new(userEntity))
	{
		if entities, count, err := o.OpsForSelect().Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().IfWhere(true, getTestGes()...).Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().Where(getTestGes()...).Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().Page(getTest(), pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
}

func TestSelectNullExec(t *testing.T) {
	truncateTestNullTable()
	if err := insertNullTest(); err != nil {
		t.Fatalf("TestSelectNullExec failed: %v\n", err)
	}
	o := New(new(userEntityNull))
	{
		if entities, err := o.OpsForSelect().List(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectNullExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().IfWhere(true, getTestGes()...).List(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectNullExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where(getTestGes()...).List(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectNullExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where().List(getNullTest()); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatal("TestSelectNullExec failed!")
		}
	}
}

func TestSelectNullExecPage(t *testing.T) {
	truncateTestNullTable()
	c := 10
	for i := 0; i < c; i++ {
		_ = insertNullTest()
	}
	o := New(new(userEntityNull))
	{
		if entities, count, err := o.OpsForSelect().Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().IfWhere(true, getTestGes()...).Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().Where(getTestGes()...).Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().Page(getNullTest(), pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
}

func TestSelectJoin(t *testing.T) {
	_, _ = DataSourcePool.Required("master").Exec("drop table join_master")
	_, _ = DataSourcePool.Required("master").Exec("drop table join_rel")
	_, _ = DataSourcePool.Required("master").Exec("drop table join_master_err")
	_, _ = DataSourcePool.Required("master").Exec("drop table join_rel_err")

	Register(new(_JoinMaster))
	Register(new(_JoinRel))
	Register(new(_JoinMasterError))
	Register(new(_JoinRelError))

	{
		o := New(new(_JoinMaster))
		if es, err := o.OpsForSelect().Join().CountJoin().List(nil); err != nil {
			t.Error("TestSelectJoin failed")
		} else if len(es) != 0 {
			t.Error("TestSelectJoin failed")
		}
	}

	{
		o := New(new(_JoinMaster))
		if es, err := o.OpsForSelect().Where(sg.In("xyz", 1, 2, 3, 4)).List(nil); err == nil {
			t.Error("TestSelectJoin failed")
		} else if es != nil {
			t.Error("TestSelectJoin failed")
		}
	}

	{
		o := New(new(_JoinMaster))
		jr := &_JoinRel{0, "Rel1"}
		if err := New(new(_JoinRel)).OpsForInsert().One(jr); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if err := o.OpsForInsert().One(&_JoinMaster{0, "hello", "hello", jr.ID, "", ""}); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if es, total, err := o.OpsForSelect().Join().CountJoin().Page(nil, pagination.MySql, 0, 2); err != nil {
			t.Error("TestSelectJoin failed")
		} else if len(es) <= 0 || total <= 0 {
			t.Error("TestSelectJoin failed")
		}
	}

	{
		o := New(new(_JoinMasterError))
		jr := &_JoinRelError{0, "Rel1"}
		if err := New(new(_JoinRelError)).OpsForInsert().One(jr); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if err := o.OpsForInsert().One(&_JoinMasterError{0, "hello", jr.ID, ""}); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if _, _, err := o.OpsForSelect().Join().Page(nil, pagination.MySql, 0, 2); err == nil {
			t.Error("TestSelectJoin failed")
		}
	}
}

type (
	_JoinMaster struct {
		ID       int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name     string `orm:"pk{F} c{name} def{name varchar(20) not null comment 'Name'}"`
		Name2    string `orm:"pk{F} c{name2} def{name2 varchar(20) not null comment 'Name2'}"`
		RelID    int    `orm:"c{rel_id} def{rel_id int}"`
		RelName  string `orm:"ig{T} ug{T} join{left,rel_id,join_rel,id,name}"`
		RelName2 string `orm:"ig{T} ug{T} join{left,rel_id,join_rel,id,name}"`
	}
	_JoinRel struct {
		ID   int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name string `orm:"pk{F} c{name} def{name varchar(20) not null comment 'Name'}"`
	}
	_JoinMasterError struct {
		ID   int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name string `orm:"pk{F} c{name} def{name varchar(20) not null comment 'Name'}"`

		RelID   int    `orm:"c{rel_id} def{rel_id int}"`
		RelName string `orm:"ig{T} ug{T} join{left,rel_id,join_rel,id,name2}"`
	}
	_JoinRelError struct {
		ID   int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name string `orm:"pk{F} c{name} def{name varchar(20) not null comment 'Name'}"`
	}
)

func (_ *_JoinMaster) Configure(c *EC) {
	c.Migrate = true
	c.Table = "join_master"
	c.NullFields = map[string]*NullField{
		"Name":     {"IFNULL", "", true},
		"Name2":    {"IFNULL", "''", false},
		"RelName":  {"IFNULL", "", true},
		"RelName2": {"IFNULL", "''", false},
	}
	c.JoinNullFields = map[string]*NullField{
		"Name":     {"IFNULL", "", true},
		"Name2":    {"IFNULL", "''", false},
		"RelName":  {"IFNULL", "", true},
		"RelName2": {"IFNULL", "''", false},
	}
}

func (_ *_JoinRel) Configure(c *EC) {
	c.Migrate = true
	c.Table = "join_rel"
}

func (_ *_JoinMasterError) Configure(c *EC) {
	c.Migrate = true
	c.Table = "join_master_err"
}

func (_ *_JoinRelError) Configure(c *EC) {
	c.Migrate = true
	c.Table = "join_rel_err"
}

type testSelectE struct{ ID int }

func (t *testSelectE) Configure(*EC) {}

func TestSelect(t *testing.T) {
	Register(new(testSelectE))
	d := Select(new(testSelectE))
	if d == nil {
		t.Fatal("test failed")
	}
}

type testSelectWithDsE struct{ ID int }

func (t *testSelectWithDsE) Configure(*EC) {}

func TestSelectWithDs(t *testing.T) {
	Register(new(testSelectWithDsE))
	d := SelectWithDs(new(testSelectWithDsE), "_")
	if d == nil {
		t.Fatal("test failed")
	}
}

type testSelectBeginTxE struct{ ID int }

func (t *testSelectBeginTxE) Configure(*EC) {}

func TestSelectBeginTx(t *testing.T) {
	type eE = testSelectBeginTxE
	Register(new(eE))
	d := Select(new(eE))
	if err := d.BeginTx(NewTxManager()); err != nil {
		t.Fatal("test failed!")
	}
}

type testSelectOneErrE struct {
	ID int `orm:"c{id} def{id int}"`
}

func (t *testSelectOneErrE) Configure(c *EC) {
	c.Migrate = true
}

func TestSelectOneErr(t *testing.T) {
	type eE = testSelectOneErrE
	defer func() { _, _ = DataSourcePool.Required("_").Exec("drop table testSelectOneErrE") }()
	Register(new(eE))
	if _, err := Select(new(eE)).Where(sg.C("c")).One(&eE{1}); err == nil {
		t.Fatal("test failed!")
	}
}

type testSelectOneErrSelectTooManyResult struct {
	ID int `orm:"c{id} def{id int}"`
}

func (t *testSelectOneErrSelectTooManyResult) Configure(c *EC) {
	c.Migrate = true
}

func TestSelectOneErrSelectTooManyResult(t *testing.T) {
	type eE = testSelectOneErrSelectTooManyResult
	defer func() { _, _ = DataSourcePool.Required("_").Exec("drop table testSelectOneErrE") }()
	Register(new(eE))
	iF := func() {
		if err := Insert(new(eE)).One(&eE{1}); err != nil {
			t.Fatal("test failed!")
		}
	}
	iF()
	iF()
	if _, err := Select(new(eE)).One(&eE{1}); !errors.Is(err, ErrSelectTooManyResult) {
		t.Fatal("test failed!")
	}
}
