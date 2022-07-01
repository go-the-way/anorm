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

func TestDeleteDel(t *testing.T) {
	truncateTestTable()
	o := New(new(userEntity))
	{
		if err := insertTest(); err != nil {
			t.Fatalf("TestDeleteDel prepare insert failed: %v\n", err)
		}
		if _, err := o.OpsForDelete().IfWhere(false).IfWhere(true, getTestGes()...).Del(nil); err != nil {
			t.Fatalf("TestDeleteDel failed: %v\n", err)
		}
		if selectUserEntityCount() != 0 {
			t.Fatal("TestDeleteDel failed!")
		}
	}
	{
		if err := insertTest(); err != nil {
			t.Fatalf("TestDeleteDel prepare insert failed: %v\n", err)
		}
		if _, err := o.OpsForDelete().Where(getTestGes()...).Del(nil); err != nil {
			t.Fatalf("TestDeleteDel failed: %v\n", err)
		}
		if selectUserEntityCount() != 0 {
			t.Fatal("TestDeleteDel failed!")
		}
	}
	{
		if err := insertTest(); err != nil {
			t.Fatalf("TestDeleteDel prepare insert failed: %v\n", err)
		}

		Logger.SetLogLevel(LogLevelDebug)

		if _, err := o.OpsForDelete().Del(getTest()); err != nil {
			t.Fatalf("TestDeleteDel failed: %v\n", err)
		}
		if c := selectUserEntityCount(); c != 0 {
			t.Fatal("TestDeleteDel failed!")
		}
	}
}

func TestDeleteDel2(t *testing.T) {
	truncateTestNullTable()
	o := New(new(userEntityNull))
	{
		if err := insertNullTest(); err != nil {
			t.Fatalf("TestDeleteDel prepare insert failed: %v\n", err)
		}
		if _, err := o.OpsForDelete().IfWhere(true, getTestGes()...).Del(nil); err != nil {
			t.Fatalf("TestDeleteDel failed: %v\n", err)
		}
		if selectUserEntityNullCount() > 0 {
			t.Fatal("TestDeleteDel failed!")
		}
	}
	{
		if err := insertNullTest(); err != nil {
			t.Fatalf("TestDeleteDel prepare insert failed: %v\n", err)
		}
		if _, err := o.OpsForDelete().Where(getTestGes()...).Del(nil); err != nil {
			t.Fatalf("TestDeleteDel failed: %v\n", err)
		}
		if selectUserEntityNullCount() > 0 {
			t.Fatal("TestDeleteDel failed!")
		}
	}
	{
		if err := insertNullTest(); err != nil {
			t.Fatalf("TestDeleteDel prepare insert failed: %v\n", err)
		}
		if _, err := o.OpsForDelete().Del(getNullTest()); err != nil {
			t.Fatalf("TestDeleteDel failed: %v\n", err)
		}
		if selectUserEntityNullCount() > 0 {
			t.Fatal("TestDeleteDel failed!")
		}
	}
}

func TestDeleteDelTX(t *testing.T) {
	truncateTestTable()
	o := New(new(userEntity))
	if err := insertTest(); err != nil {
		t.Fatalf("TestDeleteDelTX prepare insert failed: %v\n", err)
	}
	if err := o.BeginTx(NewTxManager()); err != nil {
		t.Fatalf("TestDeleteDelTX failed: %v\n", err)
	}
	if _, err := o.OpsForDelete().IfWhere(false).IfWhere(true, getTestGes()...).IfOnlyWhere(false).IfOnlyWhere(true, getTestGes()...).Del(nil); err != nil {
		t.Fatalf("TestDeleteDelTX failed: %v\n", err)
	}
	if err := o.txm.Commit(); err != nil {
		t.Fatalf("TestDeleteDelTX failed: %v\n", err)
	}
	if selectUserEntityCount() != 0 {
		t.Fatal("TestDeleteDelTX failed!")
	}
}

type testDeleteE struct{ ID int }

func (t *testDeleteE) Configure(*EC) {}

func TestDelete(t *testing.T) {
	Register(new(testDeleteE))
	d := Delete(new(testDeleteE))
	if d == nil {
		t.Fatal("test failed")
	}
}

type testDeleteWithDsE struct{ ID int }

func (t *testDeleteWithDsE) Configure(*EC) {}

func TestDeleteWithDs(t *testing.T) {
	Register(new(testDeleteWithDsE))
	d := DeleteWithDs(new(testDeleteWithDsE), "_")
	if d == nil {
		t.Fatal("test failed")
	}
}

type testDeleteBeginTxE struct{ ID int }

func (t *testDeleteBeginTxE) Configure(*EC) {}

func TestDeleteBeginTx(t *testing.T) {
	Register(new(testDeleteBeginTxE))
	d := Delete(new(testDeleteBeginTxE))
	if err := d.BeginTx(NewTxManager()); err != nil {
		t.Fatal("test failed!")
	}
}
