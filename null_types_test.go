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
	"testing"
	"time"
)

func TestNullTypes(t *testing.T) {
	for _, tc := range []struct {
		testName string
		val      func() any
		expect   any
	}{
		{testName: "NullString", val: func() any { return NullString("hello") }, expect: sql.NullString{String: "hello", Valid: true}},
		{testName: "NullString*", val: func() any { return NullStringPtr("hello") }, expect: &sql.NullString{String: "hello", Valid: true}},

		{testName: "NullBool", val: func() any { return NullBool(true) }, expect: sql.NullBool{Bool: true, Valid: true}},
		{testName: "NullBool*", val: func() any { return NullBoolPtr(true) }, expect: &sql.NullBool{Bool: true, Valid: true}},

		{testName: "NullByte", val: func() any { return NullByte(1) }, expect: sql.NullByte{Byte: 1, Valid: true}},
		{testName: "NullByte*", val: func() any { return NullBytePtr(1) }, expect: &sql.NullByte{Byte: 1, Valid: true}},

		{testName: "NullInt16", val: func() any { return NullInt16(1) }, expect: sql.NullInt16{Int16: 1, Valid: true}},
		{testName: "NullInt16*", val: func() any { return NullInt16Ptr(1) }, expect: &sql.NullInt16{Int16: 1, Valid: true}},

		{testName: "NullInt32", val: func() any { return NullInt32(1) }, expect: sql.NullInt32{Int32: 1, Valid: true}},
		{testName: "NullInt32*", val: func() any { return NullInt32Ptr(1) }, expect: &sql.NullInt32{Int32: 1, Valid: true}},

		{testName: "NullInt64", val: func() any { return NullInt64(1) }, expect: sql.NullInt64{Int64: 1, Valid: true}},
		{testName: "NullInt64*", val: func() any { return NullInt64Ptr(1) }, expect: &sql.NullInt64{Int64: 1, Valid: true}},

		{testName: "NullFloat64", val: func() any { return NullFloat64(1) }, expect: sql.NullFloat64{Float64: 1, Valid: true}},
		{testName: "NullFloat64*", val: func() any { return NullFloat64Ptr(1) }, expect: &sql.NullFloat64{Float64: 1, Valid: true}},

		{testName: "NullTime", val: func() any { return NullTime(time.Time{}) }, expect: sql.NullTime{Time: time.Time{}, Valid: true}},
		{testName: "NullTime*", val: func() any { return NullTimePtr(time.Time{}) }, expect: &sql.NullTime{Time: time.Time{}, Valid: true}},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			if val := tc.val(); fmt.Sprintf("%v", val) != fmt.Sprintf("%v", tc.expect) {
				t.Error("test failed")
			}
		})
	}
}
