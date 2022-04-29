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
	"reflect"
	"strings"
)

func getEntityPkgName[E Entity](entity E) string {
	return reflect.TypeOf(entity).String()
}

func getStrategyName(str string, sty Strategy) string {
	// fix ID
	str = strings.Replace(str, "ID", "Id", -1)
	switch sty {
	case Default:
		return str
	case Underline:
		// ABC => a_b_c
		name := ""
		for i, v := range str {
			if i == 0 {
				name += strings.ToLower(string(v))
			} else if v >= 'A' && v <= 'Z' {
				name += "_" + strings.ToLower(string(v))
			} else {
				name += string(v)
			}
		}
		return name
	case CamelCase:
		// ABC => aBC
		name := ""
		if len(str) == 1 {
			name = strings.ToLower(str)
		} else {
			name = strings.ToLower(string(str[0])) + str[1:]
		}
		return name
	}
	return str
}

// ScanStruct scan rows return []E
func ScanStruct[E Entity](rows *sql.Rows, entity E, complete func(entity EntityConfigurator)) ([]E, error) {
	defer func() { _ = rows.Close() }()
	ptr := reflect.TypeOf(entity)
	result := make([]E, 0)
	if columns, err := rows.Columns(); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			structVal := reflect.New(ptr.Elem())
			if err2 := rows.Scan(NewColumnPtr(structVal, columns)...); err2 != nil {
				return nil, err2
			}
			e := structVal.Interface().(E)
			if complete != nil {
				complete(e)
			}
			result = append(result, e)
		}
	}
	return result, nil
}

// NewColumnPtr return column ptr array
func NewColumnPtr(structVal reflect.Value, columns []string) []any {
	pts := make([]any, len(columns))
	for i, c := range columns {
		pts[i] = structVal.Elem().FieldByName(c).Addr().Interface()
	}
	return pts
}

func entityNotNil(entity Entity) bool {
	s := fmt.Sprintf("%v", entity)
	return s != "<nil>"
}
