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
	"reflect"

	"github.com/go-the-way/sg"
)

type (
	_Insert struct {
		orm           *orm
		ignoreColumns []sg.C
	}
)

func newInsert(o *orm) *_Insert {
	return &_Insert{orm: o, ignoreColumns: make([]sg.C, 0)}
}

func (o *_Insert) Ignore(columns ...sg.C) *_Insert {
	o.ignoreColumns = append(o.ignoreColumns, columns...)
	return o
}

func (o *_Insert) getIgnoreMap() map[string]struct{} {
	ignoreMap := modelInsertIgnoreMap[getModelPkgName(o.orm.model)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *_Insert) getInsertBuilder(models ...Model) (string, []interface{}) {
	fields := modelFieldMap[getModelPkgName(o.orm.model)]
	fieldColumnMap := modelFieldColumnMap[getModelPkgName(o.orm.model)]
	ignoreMap := o.getIgnoreMap()
	builder := sg.InsertBuilder()
	for i, model := range models {
		rt := reflect.ValueOf(model).Elem()
		argGes := make([]sg.Ge, 0)
		for _, f := range fields {
			if _, have := ignoreMap[fieldColumnMap[f]]; have {
				continue
			}
			if i == 0 {
				builder.Column(sg.C(fieldColumnMap[f]))
			}
			val := rt.FieldByName(f).Interface()
			argGes = append(argGes, sg.Arg(val))
		}

		if len(models) == 1 {
			builder.Value(sg.NewJoiner(argGes, ", ", "", "", false))
			continue
		}

		if i == 0 {
			builder.Value(sg.NewJoiner(argGes, ", ", "", ")", false))
		} else if i == len(models)-1 {
			builder.Value(sg.NewJoiner(argGes, ", ", "(", "", false))
		} else {
			builder.Value(sg.NewJoiner(argGes, ", ", "(", ")", false))
		}
	}
	return builder.Table(o.orm.table()).Build()
}

func (o *_Insert) Exec(model Model) error {
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getInsertBuilder(model)
	execInsertHookersBefore(model, &sqlStr, &ps)
	if debug() {
		(&execLog{"Insert.Exec", sqlStr, ps}).Log()
	}
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	execInsertHookersAfter(model, sqlStr, ps, err)
	if err != nil {
		return err
	}
	lastInsertId := int64(0)
	if result != nil {
		lastInsertId, _ = result.LastInsertId()
	}
	if lastInsertId > 0 {
		pks := modelPKMap[getModelPkgName(model)]
		if pks != nil && len(pks) == 1 {
			pkField := modelColumnFieldMap[getModelPkgName(model)][pks[0]]
			if pkField != "" {
				value := reflect.ValueOf(model).Elem().FieldByName(pkField)
				if value.CanSet() {
					if value.Kind() >= reflect.Int && value.Kind() <= reflect.Int64 {
						value.SetInt(lastInsertId)
					} else if value.Kind() >= reflect.Uint && value.Kind() <= reflect.Uint64 {
						value.SetUint(uint64(lastInsertId))
					}
				}
			}
		}
	}
	return nil
}

func (o *_Insert) ExecList(ignoreError bool, models ...Model) error {
	var err error
	for _, m := range models {
		err = o.Exec(m)
		if !ignoreError && err != nil {
			return err
		}
	}
	return err
}

func (o *_Insert) ExecBatch(models ...Model) (int64, error) {
	if len(models) <= 0 {
		return 0, nil
	}
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getInsertBuilder(models...)
	if debug() {
		(&execLog{"Insert.ExecBatch", sqlStr, ps}).Log()
	}
	execInsertHookersBefore(models[0], &sqlStr, &ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	execInsertHookersAfter(models[0], sqlStr, ps, err)
	if err != nil {
		return 0, err
	}
	ra := int64(0)
	if result != nil {
		if a, aErr := result.RowsAffected(); aErr != nil {
			ra = a
		}
	}
	return ra, err
}
