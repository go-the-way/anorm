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
)

func TestNilDS(t *testing.T) {
	defer func() {
		if re := recover(); re != nil && reflect.DeepEqual(re, errDSIsNil) {
			t.Log("test ok")
		}
	}()
	DS(nil)
}

func TestDS(t *testing.T) {
	DS(testDB)
	if !reflect.DeepEqual(dsMap["_"], testDB) {
		t.Fatal("call DS, expect have a key `_` DS in dsMap")
	}
	if !reflect.DeepEqual(dsMap["master"], testDB) {
		t.Fatal("call DS, expect have a key `master` DS in dsMap")
	}
}

func TestDSWithName(t *testing.T) {
	DSWithName("secondary", testDB)
	if !reflect.DeepEqual(dsMap["secondary"], testDB) {
		t.Fatal("call DSWithName, expect have a key `secondary` DS in dsMap")
	}
}

func TestDSRequired(t *testing.T) {
	defer func() {
		if re := recover(); re != nil && reflect.DeepEqual(re, errRequiredNamedDSFunc("hello")) {
			t.Logf("test ok")
		}
	}()
	dsMap.required("hello")
}
