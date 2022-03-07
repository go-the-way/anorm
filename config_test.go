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
	"reflect"
	"testing"
	"time"
)

func TestSetNilConfig(t *testing.T) {
	defer func() {
		if re := recover(); re != nil {
			if !reflect.DeepEqual(re, errInvalidSetConfig) {
				t.Fatal(`when set nil to SetConfig, expect get errInvalidSetConfig`)
			}
		}
	}()
	SetConfig(nil)
}

func TestSetDefaultConfig(t *testing.T) {
	c := &Config{}
	SetConfig(c)
	if !reflect.DeepEqual(c, Configuration) {
		t.Fatal(`when set c to SetConfig, expect get c`)
	}
}

func TestDefaultConfig(t *testing.T) {
	cc := &Config{Migrate: false}
	SetConfig(cc)
	if !reflect.DeepEqual(cc, Configuration) {
		t.Fatal(`when set cc to SetConfig, expect get cc`)
	}
}

type _ExecHooker struct{ int }

func (a *_ExecHooker) BeforeExec(_ Model, _ *string, _ *[]interface{}) {
	a.int += 2
}

func (a *_ExecHooker) AfterExec(_ Model, _ string, _ []interface{}, _ error) {
	a.int--
}

func TestExecHook(t *testing.T) {
	h1 := _ExecHooker{}
	h2 := _ExecHooker{}
	h3 := _ExecHooker{}
	h4 := _ExecHooker{}
	Configuration.InsertHookers = append(Configuration.InsertHookers, &h1)
	Configuration.UpdateHookers = append(Configuration.UpdateHookers, &h2)
	Configuration.DeleteHookers = append(Configuration.DeleteHookers, &h3)
	Configuration.SelectHookers = append(Configuration.SelectHookers, &h4)
	_m := userModel{0, "", 0, "", "", time.Time{}}
	_ = New(new(userModel)).Insert().Exec(&_m)
	_, _ = New(new(userModel)).Update().Exec(&_m)
	_, _ = New(new(userModel)).Delete().Exec(&_m)
	_, _ = New(new(userModel)).Select().Exec(&_m)
	if !(reflect.DeepEqual(1, h1.int) &&
		reflect.DeepEqual(1, h2.int) &&
		reflect.DeepEqual(1, h3.int) &&
		reflect.DeepEqual(1, h4.int)) {
		t.Fatal("TestExecHooker failed!")
	}
}
