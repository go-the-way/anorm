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
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type _unknownEntity struct{}

func (_ *_unknownEntity) Configure(*EC) {
}

func TestNew(t *testing.T) {
	{
		o := New(new(userEntity))
		if o == nil {
			t.Fatal("TestNew failed!")
		}
	}
	{
		o := NewWithDS(new(userEntity), "_")
		if o == nil {
			t.Fatal("TestNewWithDS failed!")
		}
		if o.db == nil {
			t.Fatal("TestNewWithDS failed!")
		}
	}
	func() {
		defer func() { _ = recover() }()
		NewWithDS(new(_unknownEntity), "_")
	}()
	func() {
		defer func() { _ = recover() }()
		var _entity *_unknownEntity
		NewWithDS(_entity, "_")
	}()
}

func TestOrmBegin(t *testing.T) {
	o := New(new(userEntity))
	if err := o.BeginTx(TxManager()); err != nil {
		t.Fatalf("TestOrmBegin failed: %v\n", err)
	}
	if !o.openTx || o.tx == nil {
		t.Fatal("TestOrmBegin failed!")
	}
}

func TestOrmBeginError(t *testing.T) {
	{
		if err := New(new(userEntity)).BeginTx(nil); err != errTxManagerNil {
			t.Error("test failed")
		}
	}
	{
		o := New(new(userEntity))
		_ = o.BeginTx(TxManager())
		if err := o.BeginTx(TxManager()); err != errAlreadyBindTxManager {
			t.Error("test failed")
		}
	}
}

func TestOrm_getRealVal(t *testing.T) {
	o := New(new(_txEntity))
	tt := time.Now()
	for _, tc := range []struct {
		testName string
		val      reflect.Value
		expect   any
	}{
		{"Int", reflect.ValueOf(1), int64(1)},
		{"Uint", reflect.ValueOf(uint(1)), uint64(1)},
		{"Float", reflect.ValueOf(float64(1)), float64(1)},
		{"String", reflect.ValueOf("hello"), "hello"},

		{"NullBool", reflect.ValueOf(NullBool(false)), NullBool(false)},

		{"NullByte", reflect.ValueOf(NullByte(1)), NullByte(1)},

		{"NullInt16", reflect.ValueOf(NullInt16(1)), NullInt16(1)},

		{"NullInt32", reflect.ValueOf(NullInt32(1)), NullInt32(1)},

		{"NullInt64", reflect.ValueOf(NullInt64(1)), NullInt64(1)},

		{"NullFloat64", reflect.ValueOf(NullFloat64(1)), NullFloat64(1)},

		{"NullString", reflect.ValueOf(NullString("1")), NullString("1")},

		{"Time", reflect.ValueOf(tt), tt},

		{"NullTime", reflect.ValueOf(NullTime(tt)), NullTime(tt)},

		{"NullBool", reflect.ValueOf(NullBoolPtr(false)), NullBool(false)},
		// ignore pointer
		//{"NullBoolPtr", reflect.ValueOf(NullBoolPtr(false)), NullBoolPtr(false)},
		//{"NullBytePtr", reflect.ValueOf(NullBytePtr(1)), NullBytePtr(1)},
		//{"NullInt16Ptr", reflect.ValueOf(NullInt16Ptr(1)), NullInt16Ptr(1)},
		//{"NullInt32Ptr", reflect.ValueOf(NullInt32Ptr(1)), NullInt32Ptr(1)},
		//{"NullInt64Ptr", reflect.ValueOf(NullInt64Ptr(1)), NullInt64Ptr(1)},
		//{"NullFloat64Ptr", reflect.ValueOf(NullFloat64Ptr(1)), NullFloat64Ptr(1)},
		//{"NullStringPtr", reflect.ValueOf(NullStringPtr("1")), NullStringPtr("1")},
		//{"TimePtr", reflect.ValueOf(&time.Time{}), &time.Time{}},
		//{"NullTimePtr", reflect.ValueOf(NullTimePtr(time.Time{})), NullTimePtr(time.Time{})},

	} {
		t.Run(tc.testName, func(t *testing.T) {
			rv := fmt.Sprintf("%v", o.getRealVal(tc.val))
			ept := fmt.Sprintf("%v", tc.expect)
			if rv != ept {
				t.Error("test failed")
			}
		})
	}
}

func TestOrmBeginTxError(t *testing.T) {
	{
		o := New(new(_txEntity))
		_ = o.db.Close()
		if err := o.BeginTx(TxManager()); err != nil && err.Error() != "sql: database is closed" {
			t.Error("test failed")
		}
		if err := o.BeginTx(TxManager(), &sql.TxOptions{}); err != nil && err.Error() != "sql: database is closed" {
			t.Error("test failed")
		}
	}
	{
		o := New(new(_txEntity))
		_ = o.db.Close()
		if err := o.BeginTx(TxManager()); err != nil && err.Error() != "sql: database is closed" {
			t.Error("test failed")
		}
	}
}

func setTxDb() {
	db, _ := sql.Open("mysql", os.Getenv("ANORM_TEST_DSN"))
	DataSourcePool.PushDB("tx", db)
}

func init() {
	setTxDb()
	Register(new(_txEntity))
}

type _txEntity struct{}

func (_ *_txEntity) Configure(c *EC) {
	c.DS = "tx"
}
