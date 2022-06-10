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

func TestInsertExec(t *testing.T) {
	truncateTestTable()
	user := userEntity{
		Name:    "hugo",
		Age:     20,
		Address: "wuhan",
		Phone:   "13900110121",
	}
	if err := New(new(userEntity)).OpsForInsert().Ignore("ig").Exec(&user); err != nil {
		t.Fatalf("TestInsertExec failed: %v\n", err)
	}
	if c := selectUserEntityCount(); c != 1 {
		t.Fatal("TestInsertExec failed!")
	}
}

type insertExecError struct {
	ID   int    `orm:"c{id} pk{T} def{id int auto_increment}"`
	Name string `orm:"c{name} def{name varchar(5)}"`
}

func (i *insertExecError) Configure(c *EC) {
	c.Table = "insert_exec_error"
	c.Migrate = true
	c.IndexDefinitions = []sg.Ge{sg.IndexDefinition(true, sg.C("uidx_name"), sg.C("name"))}
}

func init() {
	Register(new(insertExecError))
}

func TestInsertExecList(t *testing.T) {
	truncateTestTable()
	user1 := userEntity{Name: "hugo1", Age: 20, Address: "wuhan", Phone: "13900110121"}
	user2 := userEntity{Name: "hugo2", Age: 21, Address: "wuhan", Phone: "13900110122"}
	user3 := userEntity{Name: "hugo3", Age: 22, Address: "wuhan", Phone: "13900110123"}
	if err := New(new(userEntity)).OpsForInsert().ExecList(true, &user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertExecList failed: %v\n", err)
	}
	if c := selectUserEntityCount(); c != 3 {
		t.Fatal("TestInsertExecList failed!")
	}
}

func TestInsertExecListError(t *testing.T) {
	user1 := insertExecError{Name: "hugo1"}
	user2 := insertExecError{Name: "hugo1"}
	user3 := insertExecError{Name: "hugo1"}
	if err := New(new(insertExecError)).OpsForInsert().ExecList(false, &user1, &user2, &user3); err == nil {
		t.Fatalf("TestInsertExecList failed: %v\n", err)
	}
	if c, _ := New(new(insertExecError)).OpsForSelectCount().Exec(nil); c != 1 {
		t.Fatal("TestInsertExecList failed!")
	}
}

func TestInsertExecBatchError(t *testing.T) {
	user1 := insertExecError{Name: "hugo1"}
	user2 := insertExecError{Name: "hugo1"}
	user3 := insertExecError{Name: "hugo1"}
	o := New(new(insertExecError))
	txm := NewTxManager()
	o.BeginTx(txm)
	_ = txm.Rollback()
	if _, err := o.OpsForInsert().ExecBatch(&user1, &user2, &user3); err == nil {
		t.Fatalf("TestInsertExecList failed: %v\n", err)
	}
}

func TestInsertExecBatch(t *testing.T) {
	truncateTestTable()
	user1 := userEntity{Name: "hugo1", Age: 20, Address: "wuhan", Phone: "13900110122"}
	user2 := userEntity{Name: "hugo2", Age: 21, Address: "wuhan", Phone: "13900110122"}
	user3 := userEntity{Name: "hugo3", Age: 22, Address: "wuhan", Phone: "13900110123"}
	if c, err := New(new(userEntity)).OpsForInsert().ExecBatch(); !(err == nil && c == 0) {
		t.Fatal("TestInsertExecBatch failed!")
	}
	if _, err := New(new(userEntity)).OpsForInsert().ExecBatch(&user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertExecBatch failed: %v\n", err)
	}
	if c := selectUserEntityCount(); c != 3 {
		t.Fatal("TestInsertExecBatch failed!")
	}
}

func TestInsertNullExec(t *testing.T) {
	truncateTestNullTable()
	user := userEntityNull{
		Name:    NullString("hugo"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110121"),
	}
	if err := New(new(userEntityNull)).OpsForInsert().Exec(&user); err != nil {
		t.Fatalf("TestInsertNullExec failed: %v\n", err)
	}
	if c := selectUserEntityNullCount(); c != 1 {
		t.Fatal("TestInsertNullExec failed!")
	}
}

func TestInsertNullExecList(t *testing.T) {
	truncateTestNullTable()
	user1 := userEntityNull{
		Name:    NullString("hugo"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110121"),
	}
	user2 := userEntityNull{
		Name:    NullString("hugo2"),
		Age:     NullInt32(21),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user3 := userEntityNull{
		Name:    NullString("hugo3"),
		Age:     NullInt32(22),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110123"),
	}
	if err := New(new(userEntityNull)).OpsForInsert().ExecList(true, &user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertNullExecList failed: %v\n", err)
	}
	if c := selectUserEntityNullCount(); c != 3 {
		t.Fatalf("TestInsertNullExecList failed!")
	}
}

func TestInsertNullExecBatch(t *testing.T) {
	truncateTestNullTable()
	user1 := userEntityNull{
		Name:    NullString("hugo1"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user2 := userEntityNull{
		Name:    NullString("hugo2"),
		Age:     NullInt32(21),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user3 := userEntityNull{
		Name:    NullString("hugo1"),
		Age:     NullInt32(22),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110123"),
	}
	if _, err := New(new(userEntityNull)).OpsForInsert().ExecBatch(&user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertNullExecBatch failed: %v\n", err)
	}
	if c := selectUserEntityNullCount(); c != 3 {
		t.Fatalf("TestInsertNullExecBatch failed!")
	}
}

func TestInsertUintPK(t *testing.T) {
	truncateTestUintTable()
	o := New(new(userUintEntity))
	m := userUintEntity{Name: testName}
	if err := o.OpsForInsert().Exec(&m); err != nil {
		t.Fatalf("TestInsertUintPK failed: %v\n", err)
	}
	if m.ID == 0 {
		t.Fatal("TestInsertUintPK failed!")
	}
}

func TestInsertError(t *testing.T) {
	truncateTestTable()
	o := New(new(userEntity))
	m := userEntity{Name: testName}
	txm := NewTxManager()
	if err := o.BeginTx(txm); err != nil {
		t.Fatalf("TestInsertError failed: %v\n", err)
	}
	_ = txm.Rollback()
	if err := o.OpsForInsert().Exec(&m); err == nil {
		t.Fatal("TestInsertError failed")
	}
}
