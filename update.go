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
	"github.com/go-the-way/sg"
	"reflect"
)

type _Update struct {
	orm           *orm
	setColumns    []sg.C
	ignoreColumns []sg.C
	wheres        []sg.Ge
}

func newUpdate(o *orm) *_Update {
	return &_Update{orm: o, ignoreColumns: make([]sg.C, 0), wheres: make([]sg.Ge, 0)}
}

func (o *_Update) Ignore(columns ...sg.C) *_Update {
	o.ignoreColumns = append(o.ignoreColumns, columns...)
	return o
}

func (o *_Update) Set(columns ...sg.C) *_Update {
	o.setColumns = append(o.setColumns, columns...)
	return o
}

func (o *_Update) getIgnoreMap() map[string]struct{} {
	ignoreMap := modelUpdateIgnoreMap[getModelPkgName(o.orm.model)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *_Update) getSetMap() map[string]struct{} {
	setMap := make(map[string]struct{}, 0)
	for _, c := range o.setColumns {
		setMap[string(c)] = struct{}{}
	}
	return setMap
}

func (o *_Update) IfWhere(cond bool, wheres ...sg.Ge) *_Update {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *_Update) Where(wheres ...sg.Ge) *_Update {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *_Update) getUpdateBuilder(model Model) (string, []interface{}) {
	fields := modelFieldMap[getModelPkgName(o.orm.model)]
	fieldColumnMap := modelFieldColumnMap[getModelPkgName(o.orm.model)]
	pks := modelPKMap[getModelPkgName(o.orm.model)]
	pkMap := make(map[string]struct{}, 0)
	for _, pk := range pks {
		pkMap[pk] = struct{}{}
	}
	ignoreMap := o.getIgnoreMap()
	setMap := o.getSetMap()
	builder := sg.UpdateBuilder()
	setGes := make([]sg.Ge, 0)
	whereGes := make([]sg.Ge, 0)
	whereGes = append(whereGes, o.wheres...)
	rt := reflect.ValueOf(model).Elem()
	for _, f := range fields {
		column := fieldColumnMap[f]
		val := rt.FieldByName(f).Interface()
		if _, have := setMap[column]; !have {
			if _, have = pkMap[column]; have {
				whereGes = append(whereGes, sg.Eq(sg.C(column), val))
				continue
			}
		}
		if len(setMap) > 0 {
			if _, have := setMap[column]; have {
				setGes = append(setGes, sg.SetEq(sg.C(column), val))
			}
		} else {
			if _, have := ignoreMap[column]; !have {
				setGes = append(setGes, sg.SetEq(sg.C(column), val))
			}
		}
	}
	return builder.Set(setGes...).Where(sg.AndGroup(whereGes...)).Update(o.orm.table()).Build()
}

func (o *_Update) Exec(model Model) (int64, error) {
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getUpdateBuilder(model)
	execUpdateHookersBefore(model, &sqlStr, &ps)
	if debug() {
		(&execLog{"Update.Exec", sqlStr, ps}).Log()
	}
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	execUpdateHookersAfter(model, sqlStr, ps, err)
	ra := int64(0)
	if a, aErr := result.RowsAffected(); aErr != nil {
		ra = a
	}
	return ra, err
}
