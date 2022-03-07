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

func TestSelectCountExec(t *testing.T) {
	truncateTestTable()
	_ = insertTest()
	o := New(new(userModel))
	{
		if c, err := o.SelectCount().Exec(nil); err != nil {
			t.Fatalf("TestSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestSelectCountExec failed!")
		}
	}
	{
		if c, err := o.SelectCount().Exec(getTest()); err != nil {
			t.Fatalf("TestSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestSelectCountExec failed!")
		}
	}
	{
		if c, err := o.SelectCount().IfWhere(true, getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestSelectCountExec failed!")
		}
	}
}

func TestNullSelectCountExec(t *testing.T) {
	truncateTestNullTable()
	_ = insertNullTest()
	o := New(new(userModelNull))
	{
		if c, err := o.SelectCount().Exec(nil); err != nil {
			t.Fatalf("TestNullSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestNullSelectCountExec failed!")
		}
	}
	{
		if c, err := o.SelectCount().Exec(getNullTest()); err != nil {
			t.Fatalf("TestNullSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestNullSelectCountExec failed!")
		}
	}
	{
		if c, err := o.SelectCount().IfWhere(true, getTestGes()...).Exec(nil); err != nil {
			t.Fatalf("TestNullSelectCountExec failed: %v\n", err)
		} else if c != 1 {
			t.Fatal("TestNullSelectCountExec failed!")
		}
	}
}
