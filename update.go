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
	"github.com/go-the-way/sg"
	"reflect"
)

type ormUpdate struct {
	orm           *orm
	setColumns    []sg.C
	ignoreColumns []sg.C
	wheres        []sg.Ge
}

func (o *ormUpdate) Ignore(columns ...sg.C) *ormUpdate {
	o.ignoreColumns = append(o.ignoreColumns, columns...)
	return o
}

func (o *ormUpdate) Set(columns ...sg.C) *ormUpdate {
	o.setColumns = append(o.setColumns, columns...)
	return o
}

func (o *ormUpdate) getIgnoreMap() map[string]struct{} {
	ignoreMap := modelUpdateIgnoreMap[getModelPkgName(o.orm.model)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *ormUpdate) getSetMap() map[string]struct{} {
	setMap := make(map[string]struct{}, 0)
	for _, c := range o.setColumns {
		setMap[string(c)] = struct{}{}
	}
	return setMap
}

func (o *ormUpdate) NotIfWhere(cond bool, wheres ...sg.Ge) *ormUpdate {
	return o.IfWhere(!cond, wheres...)
}

func (o *ormUpdate) IfWhere(cond bool, wheres ...sg.Ge) *ormUpdate {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *ormUpdate) Where(wheres ...sg.Ge) *ormUpdate {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *ormUpdate) getUpdateBuilder(model Model) (string, []interface{}) {
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

func (o *ormUpdate) Modify(model Model) error {
	sqlStr, ps := o.getUpdateBuilder(model)
	debug("Modify[sql:%v ps:%v]", sqlStr, ps)
	execUpdateHookersBefore(model, &sqlStr, &ps)
	_, err := newExecutorFromOrm(o.orm).exec(sqlStr, ps...)
	execUpdateHookersAfter(model, sqlStr, ps, err)
	if err != nil {
		return err
	}
	return nil
}
