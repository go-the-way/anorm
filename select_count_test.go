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

func TestSelectCountExec(t *testing.T) {
	truncateTestTable()
	_ = insertTest()
	o := New(new(userEntity))
	{
		if c, err := o.OpsForSelectCount().Count(nil); err != nil {
			t.Fatalf("TestSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestSelectCountExec failed!")
		}
	}
	{
		if c, err := o.OpsForSelectCount().Count(getTest()); err != nil {
			t.Fatalf("TestSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestSelectCountExec failed!")
		}
	}
	{
		if c, err := o.OpsForSelectCount().IfWhere(false).IfWhere(true, getTestGes()...).Count(nil); err != nil {
			t.Fatalf("TestSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestSelectCountExec failed!")
		}
	}
}

func TestNullSelectCountExec(t *testing.T) {
	truncateTestNullTable()
	_ = insertNullTest()
	o := New(new(userEntityNull))
	{
		if c, err := o.OpsForSelectCount().Count(nil); err != nil {
			t.Fatalf("TestNullSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestNullSelectCountExec failed!")
		}
	}
	{
		if c, err := o.OpsForSelectCount().Count(getNullTest()); err != nil {
			t.Fatalf("TestNullSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestNullSelectCountExec failed!")
		}
	}
	{
		if c, err := o.OpsForSelectCount().IfWhere(true, getTestGes()...).Join().Count(nil); err != nil {
			t.Fatalf("TestNullSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestNullSelectCountExec failed!")
		}
	}
}

func TestSelectCountError(t *testing.T) {
	if _, err := New(new(userEntityNull)).OpsForSelectCount().Where(sg.Eq("ta", 100)).Count(nil); err == nil {
		t.Error("TestSelectCountError failed")
	}
}

type testSelectCountE struct{ ID int }

func (t *testSelectCountE) Configure(*EC) {}

func TestSelectCount(t *testing.T) {
	Register(new(testSelectCountE))
	d := SelectCount(new(testSelectCountE))
	if d == nil {
		t.Fatal("test failed")
	}
}

type testSelectCountWithDsE struct{ ID int }

func (t *testSelectCountWithDsE) Configure(*EC) {}

func TestSelectCountWithDs(t *testing.T) {
	Register(new(testSelectCountWithDsE))
	d := SelectCountWithDs(new(testSelectCountWithDsE), "_")
	if d == nil {
		t.Fatal("test failed")
	}
}

type testSelectCountBeginTxE struct{ ID int }

func (t *testSelectCountBeginTxE) Configure(*EC) {}

func TestSelectCountBeginTx(t *testing.T) {
	Register(new(testSelectCountBeginTxE))
	d := SelectCount(new(testSelectCountBeginTxE))
	if err := d.BeginTx(NewTxManager()); err != nil {
		t.Fatal("test failed!")
	}
}
