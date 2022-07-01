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

func TestGetStrategyName(t *testing.T) {
	if s := getStrategyName("ABC", Underline); s != "a_b_c" {
		t.Error("TestGetStrategyName failed")
	}
	if s := getStrategyName("ABC_", Underline); s != "a_b_c_" {
		t.Error("TestGetStrategyName failed")
	}
	if s := getStrategyName("ABC", CamelCase); s != "aBC" {
		t.Error("TestGetStrategyName failed")
	}
	if s := getStrategyName("A", CamelCase); s != "a" {
		t.Error("TestGetStrategyName failed")
	}
	if s := getStrategyName("A", -1); s != "A" {
		t.Error("TestGetStrategyName failed")
	}
}

type (
	_scanStruct struct{ Name, Name2 string }
	_tStruct    struct{ T string }
)

func (_ *_scanStruct) Configure(*EC) {}
func (_ *_tStruct) Configure(*EC)    {}

func TestScanStruct(t *testing.T) {
	{
		if rows, err := testDB.Query("select now() as T"); err != nil {
			t.Error("TestScanStruct failed")
		} else {
			if _, err2 := ScanStruct(rows, new(_tStruct), func(entity EntityConfigurator) {}); err2 != nil {
				t.Error("TestScanStruct failed")
			}
		}
	}
	{
		if rows, err := testDB.Query("select now() as t"); err != nil {
			t.Error("TestScanStruct failed")
		} else {
			_ = rows.Close()
			if _, err2 := ScanStruct(rows, new(_scanStruct), func(entity EntityConfigurator) {}); err2 == nil {
				t.Error("TestScanStruct failed")
			}
		}
	}
	{
		if rows, err := testDB.Query("select NULL as Name,111 as Name2"); err != nil {
			t.Error("TestScanStruct failed")
		} else {
			if _, err2 := ScanStruct(rows, new(_scanStruct), func(entity EntityConfigurator) {}); err2 == nil {
				t.Error("TestScanStruct failed")
			}
		}
	}
}
