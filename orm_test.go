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

func TestNew(t *testing.T) {
	o := New(new(userModel))
	if o == nil {
		t.Fatal("TestNew failed!")
	}
}

func TestNewWithDS(t *testing.T) {
	o := NewWithDS(new(userModel), "_")
	if o == nil {
		t.Fatal("TestNewWithDS failed!")
	}
	if o.db == nil {
		t.Fatal("TestNewWithDS failed!")
	}
}

func TestOrmBegin(t *testing.T) {
	o := New(new(userModel))
	if err := o.Begin(); err != nil {
		t.Fatalf("TestOrmBegin failed: %v\n", err)
	}
	if !o.openTx || o.tx == nil {
		t.Fatal("TestOrmBegin failed!")
	}
}

func TestOrmRollback(t *testing.T) {
	o := New(new(userModel))
	if err := o.Begin(); err != nil {
		t.Fatalf("TestOrmRollback failed: %v\n", err)
	}
	truncateTestTable()
	if err := o.Insert().Exec(getTest()); err != nil {
		t.Fatalf("TestOrmRollback failed: %v\n", err)
	}
	if err := o.Commit(); err != nil {
		t.Fatalf("TestOrmRollback failed: %v\n", err)
	}
	if c := selectUserModelCount(); c == 0 {
		t.Fatal("TestOrmRollback failed!")
	}
}
