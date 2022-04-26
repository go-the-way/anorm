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
	"time"
)

// NullString return sql.NullString
func NullString(str string) sql.NullString {
	return *NullStringPtr(str)
}

// NullStringPtr return *sql.NullString
func NullStringPtr(str string) *sql.NullString {
	return &sql.NullString{String: str, Valid: true}
}

// NullBool return sql.NullBool
func NullBool(b bool) sql.NullBool {
	return *NullBoolPtr(b)
}

// NullBoolPtr return *sql.NullBool
func NullBoolPtr(b bool) *sql.NullBool {
	return &sql.NullBool{Bool: b, Valid: true}
}

// NullByte return sql.NullByte
func NullByte(b byte) sql.NullByte {
	return *NullBytePtr(b)
}

// NullBytePtr return *sql.NullByte
func NullBytePtr(b byte) *sql.NullByte {
	return &sql.NullByte{Byte: b, Valid: true}
}

// NullInt16 return sql.NullInt16
func NullInt16(i int16) sql.NullInt16 {
	return *NullInt16Ptr(i)
}

// NullInt16Ptr return *sql.NullInt16
func NullInt16Ptr(i int16) *sql.NullInt16 {
	return &sql.NullInt16{Int16: i, Valid: true}
}

// NullInt32 return sql.NullInt32
func NullInt32(i int32) sql.NullInt32 {
	return *NullInt32Ptr(i)
}

// NullInt32Ptr return *sql.NullInt32
func NullInt32Ptr(i int32) *sql.NullInt32 {
	return &sql.NullInt32{Int32: i, Valid: true}
}

// NullInt64 return sql.NullInt64
func NullInt64(i int64) sql.NullInt64 {
	return *NullInt64Ptr(i)
}

// NullInt64Ptr return *sql.NullInt64
func NullInt64Ptr(i int64) *sql.NullInt64 {
	return &sql.NullInt64{Int64: i, Valid: true}
}

// NullFloat64 return sql.NullFloat64
func NullFloat64(f float64) sql.NullFloat64 {
	return *NullFloat64Ptr(f)
}

// NullFloat64Ptr return *sql.NullFloat64
func NullFloat64Ptr(f float64) *sql.NullFloat64 {
	return &sql.NullFloat64{Float64: f, Valid: true}
}

// NullTime return sql.NullTime
func NullTime(t time.Time) sql.NullTime {
	return *NullTimePtr(t)
}

// NullTimePtr return *sql.NullTime
func NullTimePtr(t time.Time) *sql.NullTime {
	return &sql.NullTime{Time: t, Valid: true}
}
