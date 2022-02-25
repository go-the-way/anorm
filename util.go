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
	"reflect"
	"strings"
)

func getModelPkgName(model Model) string {
	return reflect.TypeOf(model).String()
}

func getModelTableName(model Model, sty Strategy) string {
	pkgName := reflect.TypeOf(model).String()
	pss := strings.Split(pkgName, ".")
	tableName := pss[len(pss)-1]
	metadata := model.MetaData()
	if t := metadata.Table; t != "" {
		tableName = t
	} else {
		tableName = getStrategyName(tableName, sty)
	}
	return tableName
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

func scanStruct(rows *sql.Rows, model Model) ([]Model, error) {
	ptr := reflect.TypeOf(model)
	defer func() {
		if rows != nil {
			err := rows.Close()
			handleErr(err)
		}
	}()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := make([]Model, 0)
	for rows.Next() {
		structVal := reflect.New(ptr.Elem())
		handleErr(rows.Scan(newColumnPtr(structVal, columns)...))
		result = append(result, structVal.Interface().(Model))
	}
	return result, nil
}

func newColumnPtr(structVal reflect.Value, columns []string) []interface{} {
	pts := make([]interface{}, len(columns))
	for i, c := range columns {
		pts[i] = structVal.Elem().FieldByName(c).Addr().Interface()
	}
	return pts
}

func ifEmpty(str, defaultVal string) string {
	if str == "" {
		return defaultVal
	}
	return str
}
