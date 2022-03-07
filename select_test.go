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
	o := New(new(userModel))
	{
		if models, err := o.Select().Exec(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
	{
		if models, err := o.Select().IfWhere(true, getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
	{
		if models, err := o.Select().Where(getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
	{
		if models, err := o.Select().Where().Exec(getTest()); err != nil {
			t.Fatalf("TestSelectExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectExec failed!")
		}
	}
}

func TestSelectExecPage(t *testing.T) {
	truncateTestTable()
	c := 10
	for i := 0; i < c; i++ {
		_ = insertTest()
	}
	o := New(new(userModel))
	{
		if models, count, err := o.Select().ExecPage(nil, MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if models, count, err := o.Select().IfWhere(true, getTestGes()...).ExecPage(nil, MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if models, count, err := o.Select().Where(getTestGes()...).ExecPage(nil, MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if models, count, err := o.Select().ExecPage(getTest(), MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
}

func TestSelectNullExec(t *testing.T) {
	truncateTestNullTable()
	if err := insertNullTest(); err != nil {
		t.Fatalf("TestSelectNullExec failed: %v\n", err)
	}
	o := New(new(userModelNull))
	{
		if models, err := o.Select().Exec(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
		}
	}
	{
		if models, err := o.Select().IfWhere(true, getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
		}
	}
	{
		if models, err := o.Select().Where(getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(models) != 1 {
			t.Fatalf("TestSelectNullExec failed!")
		}
	}
	{
		if models, err := o.Select().Where().Exec(getNullTest()); err != nil {
			t.Fatalf("TestSelectNullExec failed: %v\n", err)
		} else if len(models) != 1 {
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
	o := New(new(userModelNull))
	{
		if models, count, err := o.Select().ExecPage(nil, MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if models, count, err := o.Select().IfWhere(true, getTestGes()...).ExecPage(nil, MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if models, count, err := o.Select().Where(getTestGes()...).ExecPage(nil, MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
	{
		if models, count, err := o.Select().ExecPage(getNullTest(), MySQLPagination, 0, 2); err != nil {
			t.Fatalf("TestSelectExecPage failed: %v\n", err)
		} else if len(models) != 2 && int(count) != c {
			t.Fatal("TestSelectExecPage failed!")
		}
	}
}
