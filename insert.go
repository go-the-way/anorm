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
	"errors"
	"fmt"
	"reflect"

	"github.com/go-the-way/sg"
)

var (
	errModelsIsNilOrEmpty = errors.New("models is nil or empty")
)

type ormInsert struct {
	orm           *orm
	ignoreColumns []sg.C
}

type argAppend struct {
	prepend string
	append  string
	p       interface{}
}

func (a *argAppend) SQL() (string, []interface{}) {
	return fmt.Sprintf("%s?%s", a.prepend, a.append), []interface{}{a.p}
}

func (o *ormInsert) Ignore(columns ...sg.C) *ormInsert {
	o.ignoreColumns = append(o.ignoreColumns, columns...)
	return o
}

func (o *ormInsert) getIgnoreMap() map[string]struct{} {
	ignoreMap := modelInsertIgnoreMap[getModelPkgName(o.orm.model)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *ormInsert) getInsertBuilder(models ...Model) (string, []interface{}) {
	fields := modelFieldMap[getModelPkgName(o.orm.model)]
	fieldColumnMap := modelFieldColumnMap[getModelPkgName(o.orm.model)]
	ignoreMap := o.getIgnoreMap()
	builder := sg.InsertBuilder()
	for i, model := range models {
		rt := reflect.ValueOf(model).Elem()
		argGes := make([]sg.Ge, 0)
		idx := 0
		for j, f := range fields {
			if _, have := ignoreMap[fieldColumnMap[f]]; !have {
				idx++
				if i == 0 {
					builder.Column(sg.C(fieldColumnMap[f]))
				}
				val := rt.FieldByName(f).Interface()
				if j == len(fields)-1 {
					if i == len(models)-1 {
						argGes = append(argGes, sg.Arg(val))
					} else {
						argGes = append(argGes, &argAppend{"", ")", val})
					}
				} else if idx == 1 && i > 0 {
					argGes = append(argGes, &argAppend{"(", "", val})
				} else {
					argGes = append(argGes, sg.Arg(val))
				}
			}
		}
		builder.Value(sg.NewJoiner(argGes, ", ", "", "", false))
	}
	return builder.Table(o.orm.table()).Build()
}

func (o *ormInsert) save(etr *executor, models ...Model) error {
	model := o.orm.model
	var (
		sqlStr string
		ps     []interface{}
		result sql.Result
		err    error
	)

	sqlStr, ps = o.getInsertBuilder(models...)
	curModel := models[0]

	debug("Save[sql:%v ps:%v]", sqlStr, ps)
	execInsertHookersBefore(model, &sqlStr, &ps)
	if etr == nil {
		etr = newExecutorFromOrm(o.orm)
	}
	result, err = etr.exec(sqlStr, ps...)
	execInsertHookersAfter(model, sqlStr, ps, err)

	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if lastInsertId == 0 {
		return err
	}

	pks := modelPKMap[getModelPkgName(model)]
	if pks != nil && len(pks) == 1 {
		pkField := modelColumnFieldMap[getModelPkgName(model)][pks[0]]
		if pkField != "" {
			value := reflect.ValueOf(curModel).Elem().FieldByName(pkField)
			if value.CanSet() {
				switch {
				default:
					debug("model[%v] set lastInsertId ignore of pk field type[%v]", getModelPkgName(model), value.Kind().String())
				case value.Kind() >= reflect.Int && value.Kind() <= reflect.Int64:
					value.SetInt(lastInsertId)
				case value.Kind() >= reflect.Uint && value.Kind() <= reflect.Uint64:
					value.SetUint(uint64(lastInsertId))
				}
			}
		}
	}
	return nil
}

func (o *ormInsert) Save(model Model) error {
	return o.save(nil, model)
}

func (o *ormInsert) saveAll(etr *executor, ignoreError bool, models ...Model) error {
	if models == nil || len(models) <= 0 {
		panic(errModelsIsNilOrEmpty)
	}
	var (
		tx   *sql.Tx
		err  error
		isTx bool
	)
	if etr == nil {
		isTx = true
		if o.orm.tx != nil {
			tx = o.orm.tx
		} else {
			tx, err = dsMap.dsHold(o.orm.db, "_").Begin()
			if err != nil {
				return err
			}
		}
		etr = newExecutorFromOrm(o.orm)
	}
	for i := range models {
		err1 := o.save(etr, models[i])
		if err1 != nil {
			err = err1
			if ignoreError {
				handleErr(err)
			} else {
				if isTx && o.orm.tx == nil {
					_ = tx.Rollback()
				}
				return err
			}
		}
	}
	if isTx && o.orm.tx == nil {
		_ = tx.Commit()
	}
	return err
}

func (o *ormInsert) SaveAll(ignoreError bool, models ...Model) error {
	return o.saveAll(nil, ignoreError, models...)
}

func (o *ormInsert) SaveBatch(models ...Model) error {
	return o.save(nil, models...)
}
