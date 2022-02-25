// Copyright 2022 anox Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anox

import (
	"database/sql"
	"time"
)

func NullString(str string) sql.NullString {
	return *NullStringPtr(str)
}

func NullStringPtr(str string) *sql.NullString {
	return &sql.NullString{String: str, Valid: true}
}

func NullBool(b bool) sql.NullBool {
	return *NullBoolPtr(b)
}

func NullBoolPtr(b bool) *sql.NullBool {
	return &sql.NullBool{Bool: b, Valid: true}
}

func NullByte(b byte) sql.NullByte {
	return *NullBytePtr(b)
}

func NullBytePtr(b byte) *sql.NullByte {
	return &sql.NullByte{Byte: b, Valid: true}
}

func NullInt16(i int16) sql.NullInt16 {
	return *NullInt16Ptr(i)
}

func NullInt16Ptr(i int16) *sql.NullInt16 {
	return &sql.NullInt16{Int16: i, Valid: true}
}

func NullInt32(i int32) sql.NullInt32 {
	return *NullInt32Ptr(i)
}

func NullInt32Ptr(i int32) *sql.NullInt32 {
	return &sql.NullInt32{Int32: i, Valid: true}
}

func NullInt64(i int64) sql.NullInt64 {
	return *NullInt64Ptr(i)
}

func NullInt64Ptr(i int64) *sql.NullInt64 {
	return &sql.NullInt64{Int64: i, Valid: true}
}

func NullFloat64(f float64) sql.NullFloat64 {
	return *NullFloat64Ptr(f)
}

func NullFloat64Ptr(f float64) *sql.NullFloat64 {
	return &sql.NullFloat64{Float64: f, Valid: true}
}

func NullTime(t time.Time) sql.NullTime {
	return *NullTimePtr(t)
}

func NullTimePtr(t time.Time) *sql.NullTime {
	return &sql.NullTime{Time: t, Valid: true}
}
