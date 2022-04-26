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
		if entities, err := o.OpsForSelect().OrderBy(sg.C("id")).Exec(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().IfWhere(false).IfWhere(true, getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where(getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where().Exec(getTest()); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectExec failed!")
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
		if ee, err := o.OpsForSelect().ExecOne(nil); err != nil {
			t.Fatalf("TestSelectExecOne failed: %v\n", err)
		} else if ee == nil {
			t.Fatalf("TestSelectExecOne failed!")
		}
	}
}

func TestSelectExecOne2(t *testing.T) {
	truncateTestTable()
	o := New(new(userEntity))
	{
		if ee, err := o.OpsForSelect().ExecOne(nil); err != nil {
			t.Fatalf("TestSelectExecOne failed: %v\n", err)
		} else if ee != nil {
			t.Fatalf("TestSelectExecOne failed!")
		}
	}
}

func TestSelectExecPageError(t *testing.T) {
	{
		truncateTestTable()
		o := New(new(userEntity))
		{
			if _, _, err := o.OpsForSelect().Where(sg.And(sg.Eq("1", 1))).ExecPage(nil, pagination.Pg, 0, 2); err == nil {
				t.Fatalf("TestSelectExecPageError failed: %v\n", err)
			}
		}
	}
	{
		o := New(new(userEntity))
		{
			if entities, count, err := o.OpsForSelect().ExecPage(nil, pagination.MySql, 0, 2); err != nil {
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
		if entities, count, err := o.OpsForSelect().ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().IfWhere(true, getTestGes()...).ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().Where(getTestGes()...).ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().ExecPage(getTest(), pagination.MySql, 0, 2); err != nil {
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
		if entities, err := o.OpsForSelect().Exec(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().IfWhere(true, getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where(getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
		}
	}
	{
		if entities, err := o.OpsForSelect().Where().Exec(getNullTest()); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(entities) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
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
		if entities, count, err := o.OpsForSelect().ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().IfWhere(true, getTestGes()...).ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().Where(getTestGes()...).ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if entities, count, err := o.OpsForSelect().ExecPage(getNullTest(), pagination.MySql, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(entities) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
}

func TestSelectJoin(t *testing.T) {
	Register(new(_JoinMaster))
	Register(new(_JoinRel))
	Register(new(_JoinMasterError))
	Register(new(_JoinRelError))

	_, _ = DataSourcePool.Required("master").Exec("truncate table join_master")
	_, _ = DataSourcePool.Required("master").Exec("truncate table join_rel")
	_, _ = DataSourcePool.Required("master").Exec("truncate table join_master_err")
	_, _ = DataSourcePool.Required("master").Exec("truncate table join_rel_err")

	{
		o := New(new(_JoinMaster))
		if es, err := o.OpsForSelect().Join().Exec(nil); err != nil {
			t.Error("TestSelectJoin failed")
		} else if len(es) != 0 {
			t.Error("TestSelectJoin failed")
		}
	}

	{
		o := New(new(_JoinMaster))
		if es, err := o.OpsForSelect().Where(sg.In("xyz", 1, 2, 3, 4)).Exec(nil); err == nil {
			t.Error("TestSelectJoin failed")
		} else if es != nil {
			t.Error("TestSelectJoin failed")
		}
	}

	{
		o := New(new(_JoinMaster))
		jr := &_JoinRel{0, "Rel1"}
		if err := New(new(_JoinRel)).OpsForInsert().Exec(jr); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if err := o.OpsForInsert().Exec(&_JoinMaster{0, "hello", jr.ID, ""}); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if es, total, err := o.OpsForSelect().Join().ExecPage(nil, pagination.MySql, 0, 2); err != nil {
			t.Error("TestSelectJoin failed")
		} else if len(es) <= 0 || total <= 0 {
			t.Error("TestSelectJoin failed")
		}
	}

	{
		o := New(new(_JoinMasterError))
		jr := &_JoinRelError{0, "Rel1"}
		if err := New(new(_JoinRelError)).OpsForInsert().Exec(jr); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if err := o.OpsForInsert().Exec(&_JoinMasterError{0, "hello", jr.ID, ""}); err != nil {
			t.Error("TestSelectJoin failed")
		}
		if _, _, err := o.OpsForSelect().Join().ExecPage(nil, pagination.MySql, 0, 2); err == nil {
			t.Error("TestSelectJoin failed")
		}
	}
}

type (
	_JoinMaster struct {
		ID   int    `orm:"pk{T} c{id} ig{T} def{id int not null auto_increment comment 'ID'}"`
		Name string `orm:"pk{F} c{name} def{name varchar(20) not null comment 'Name'}"`

		RelID   int    `orm:"c{rel_id} def{rel_id int}"`
		RelName string `orm:"ig{T} ug{T} join{left,rel_id,join_rel,id,name}"`
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
