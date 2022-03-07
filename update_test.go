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

func init() {
	testInit()
}

func TestUpdateExec(t *testing.T) {
	testFn := func(o *orm) {
		if c, err := o.SelectCount().IfWhere(true, sg.Eq("id", 1), sg.Eq("name", "yoyo")).Exec(nil); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestUpdateExec failed!")
		}
	}
	{
		truncateTestTable()
		if err := insertTest(); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		}
		o := New(new(userModel))
		m := getTest()
		m.ID = 1
		m.Name = "yoyo"
		if _, err := o.Update().Exec(m); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		}
		testFn(o)
	}
	{
		truncateTestTable()
		if err := insertTest(); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		}
		o := New(new(userModel))
		m := getTest()
		m.ID = 1
		m.Name = "yoyo"
		if _, err := o.Update().IfWhere(true, sg.Eq("id", 1)).IfWhere(true, getTestGes()...).Exec(m); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		}
		testFn(o)
	}
	{
		truncateTestTable()
		if err := insertTest(); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		}
		o := New(new(userModel))
		m := getTest()
		m.ID = 1
		m.Name = "yoyo"
		if _, err := o.Update().Where(sg.Eq("id", 1)).Where(getTestGes()...).Exec(m); err != nil {
			t.Fatalf("TestUpdateExec failed: %v\n", err)
		}
		testFn(o)
	}
}

func TestUpdate_Set(t *testing.T) {
	truncateTestTable()
	if err := insertTest(); err != nil {
		t.Fatalf("TestUpdateSet failed: %v\n", err)
	}
	o := New(new(userModel))
	m := getTest()
	m.ID = 1
	m.Name = "yoyo"
	m.Age = testAge + 1
	if _, err := o.Update().Set("name").Exec(m); err != nil {
		t.Fatalf("TestUpdateSet failed: %v\n", err)
	}
	m.Age = testAge
	if c, err := o.SelectCount().Exec(m); err != nil {
		t.Fatalf("TestUpdateSet failed: %v\n", err)
	} else if c != 1 {
		t.Fatal("TestUpdateSet failed")
	}
}

func TestUpdate_Ignore(t *testing.T) {
	truncateTestTable()
	if err := insertTest(); err != nil {
		t.Fatalf("TestUpdateSet failed: %v\n", err)
	}
	o := New(new(userModel))
	m := getTest()
	m.ID = 1
	m.Name = "yoyo"
	m.Age = testAge + 1
	if _, err := o.Update().Ignore("age").Exec(m); err != nil {
		t.Fatalf("TestUpdate_Ignore failed: %v\n", err)
	}
	m.Age = testAge
	if c, err := o.SelectCount().Exec(m); err != nil {
		t.Fatalf("TestUpdate_Ignore failed: %v\n", err)
	} else if c != 1 {
		t.Fatal("TestUpdate_Ignore failed")
	}
}
