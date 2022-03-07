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

func TestInsertExec(t *testing.T) {
	truncateTestTable()
	user := userModel{
		Name:    "hugo",
		Age:     20,
		Address: "wuhan",
		Phone:   "13900110121",
	}
	if err := New(new(userModel)).Insert().Exec(&user); err != nil {
		t.Fatalf("TestInsertExec failed: %v\n", err)
	}
	if c := selectUserModelCount(); c != 1 {
		t.Fatal("TestInsertExec failed!")
	}
}

func TestInsertExecList(t *testing.T) {
	truncateTestTable()
	user1 := userModel{Name: "hugo1", Age: 20, Address: "wuhan", Phone: "13900110121"}
	user2 := userModel{Name: "hugo2", Age: 21, Address: "wuhan", Phone: "13900110122"}
	user3 := userModel{Name: "hugo3", Age: 22, Address: "wuhan", Phone: "13900110123"}
	if err := New(new(userModel)).Insert().ExecList(true, &user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertExecList failed: %v\n", err)
	}
	if c := selectUserModelCount(); c != 3 {
		t.Fatal("TestInsertExecList failed!")
	}
}

func TestInsertExecBatch(t *testing.T) {
	truncateTestTable()
	user1 := userModel{Name: "hugo1", Age: 20, Address: "wuhan", Phone: "13900110122"}
	user2 := userModel{Name: "hugo2", Age: 21, Address: "wuhan", Phone: "13900110122"}
	user3 := userModel{Name: "hugo3", Age: 22, Address: "wuhan", Phone: "13900110123"}
	if c, err := New(new(userModel)).Insert().ExecBatch(); !(err == nil && c == 0) {
		t.Fatal("TestInsertExecBatch failed!")
	}
	if _, err := New(new(userModel)).Insert().ExecBatch(&user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertExecBatch failed: %v\n", err)
	}
	if c := selectUserModelCount(); c != 3 {
		t.Fatal("TestInsertExecBatch failed!")
	}
}

func TestInsertNullExec(t *testing.T) {
	truncateTestNullTable()
	user := userModelNull{
		Name:    NullString("hugo"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110121"),
	}
	if err := New(new(userModelNull)).Insert().Exec(&user); err != nil {
		t.Fatalf("TestInsertNullExec failed: %v\n", err)
	}
	if c := selectUserModelNullCount(); c != 1 {
		t.Fatal("TestInsertNullExec failed!")
	}
}

func TestInsertNullExecList(t *testing.T) {
	truncateTestNullTable()
	user1 := userModelNull{
		Name:    NullString("hugo"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110121"),
	}
	user2 := userModelNull{
		Name:    NullString("hugo2"),
		Age:     NullInt32(21),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user3 := userModelNull{
		Name:    NullString("hugo3"),
		Age:     NullInt32(22),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110123"),
	}
	if err := New(new(userModelNull)).Insert().ExecList(true, &user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertNullExecList failed: %v\n", err)
	}
	if c := selectUserModelNullCount(); c != 3 {
		t.Fatalf("TestInsertNullExecList failed!")
	}
}

func TestInsertNullExecBatch(t *testing.T) {
	truncateTestNullTable()
	user1 := userModelNull{
		Name:    NullString("hugo1"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user2 := userModelNull{
		Name:    NullString("hugo2"),
		Age:     NullInt32(21),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user3 := userModelNull{
		Name:    NullString("hugo1"),
		Age:     NullInt32(22),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110123"),
	}
	if _, err := New(new(userModelNull)).Insert().ExecBatch(&user1, &user2, &user3); err != nil {
		t.Fatalf("TestInsertNullExecBatch failed: %v\n", err)
	}
	if c := selectUserModelNullCount(); c != 3 {
		t.Fatalf("TestInsertNullExecBatch failed!")
	}
}

func TestInsertUintPK(t *testing.T) {
	truncateTestUintTable()
	o := New(new(userUintModel))
	m := userUintModel{Name: testName}
	if err := o.Insert().Exec(&m); err != nil {
		t.Fatalf("TestInsertUintPK failed: %v\n", err)
	}
	if m.ID == 0 {
		t.Fatal("TestInsertUintPK failed!")
	}
}
