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

type (
	UpdateOperation[E Entity] interface {
		BeginTx(txm *TxManager, options ...*sql.TxOptions) error
		Ignore(columns ...sg.C) UpdateOperation[E]
		Set(columns ...sg.C) UpdateOperation[E]
		IfWhere(cond bool, wheres ...sg.Ge) UpdateOperation[E]
		IfOnlyWhere(cond bool, wheres ...sg.Ge) UpdateOperation[E]
		Where(wheres ...sg.Ge) UpdateOperation[E]
		OnlyWhere(wheres ...sg.Ge) UpdateOperation[E]
		UpByPK(e E) (c int64, err error)
	}
	updateOperation[E Entity] struct {
		orm                       *Orm[E]
		setColumns, ignoreColumns []sg.C
		wheres, onlyWheres        []sg.Ge
	}
)

func Update[E Entity](e E) UpdateOperation[E] {
	return New(e).OpsForUpdate()
}

func UpdateWithDs[E Entity](e E, ds string) UpdateOperation[E] {
	return NewWithDS(e, ds).OpsForUpdate()
}

func newUpdateOperation[E Entity](o *Orm[E]) *updateOperation[E] {
	return &updateOperation[E]{orm: o, ignoreColumns: make([]sg.C, 0), wheres: make([]sg.Ge, 0), onlyWheres: make([]sg.Ge, 0)}
}

func (o *updateOperation[E]) BeginTx(txm *TxManager, options ...*sql.TxOptions) error {
	return o.orm.BeginTx(txm, options...)
}

// Ignore ignore columns for updates
func (o *updateOperation[E]) Ignore(columns ...sg.C) UpdateOperation[E] {
	o.ignoreColumns = append(o.ignoreColumns, columns...)
	return o
}

// Set if set only set columns
func (o *updateOperation[E]) Set(columns ...sg.C) UpdateOperation[E] {
	o.setColumns = append(o.setColumns, columns...)
	return o
}

func (o *updateOperation[E]) getIgnoreMap() map[string]struct{} {
	ignoreMap := entityUpdateIgnoreMap[getEntityPkgName(o.orm.entity)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *updateOperation[E]) getSetMap() map[string]struct{} {
	setMap := make(map[string]struct{}, 0)
	for _, c := range o.setColumns {
		setMap[string(c)] = struct{}{}
	}
	return setMap
}

// IfWhere if cond is true, append wheres
func (o *updateOperation[E]) IfWhere(cond bool, wheres ...sg.Ge) UpdateOperation[E] {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

// IfOnlyWhere if cond is true, append only wheres
func (o *updateOperation[E]) IfOnlyWhere(cond bool, wheres ...sg.Ge) UpdateOperation[E] {
	if cond {
		return o.OnlyWhere(wheres...)
	}
	return o
}

// Where append wheres
func (o *updateOperation[E]) Where(wheres ...sg.Ge) UpdateOperation[E] {
	o.wheres = append(o.wheres, wheres...)
	return o
}

// OnlyWhere append only wheres
func (o *updateOperation[E]) OnlyWhere(wheres ...sg.Ge) UpdateOperation[E] {
	o.onlyWheres = append(o.onlyWheres, wheres...)
	return o
}

func (o *updateOperation[E]) getUpdateBuilder(entity E) (string, []any) {
	fields := entityFieldMap[getEntityPkgName(o.orm.entity)]
	fieldColumnMap := entityFieldColumnMap[getEntityPkgName(o.orm.entity)]
	pks := entityPKMap[getEntityPkgName(o.orm.entity)]
	pkMap := make(map[string]struct{}, 0)
	for _, pk := range pks {
		pkMap[pk] = struct{}{}
	}
	ignoreMap := o.getIgnoreMap()
	setMap := o.getSetMap()
	builder := sg.UpdateBuilder()
	setGes := make([]sg.Ge, 0)
	whereGes := make([]sg.Ge, 0)
	appendEntityWhere := len(o.onlyWheres) <= 0
	rt := reflect.ValueOf(entity).Elem()
	if appendEntityWhere {
		whereGes = append(whereGes, o.wheres...)
	} else {
		whereGes = append(whereGes, o.onlyWheres...)
	}
	for _, f := range fields {
		column := fieldColumnMap[f]
		val := rt.FieldByName(f).Interface()
		if appendEntityWhere {
			if _, have := setMap[column]; !have {
				if _, have = pkMap[column]; have {
					whereGes = append(whereGes, sg.Eq(sg.C(column), val))
					continue
				}
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

// UpByPK select for page
//
// Params:
//
// - e: the orm wrapper entity
//
// Returns:
//
// - count: RowsAffected count
//
// - err: exec error
//
func (o *updateOperation[E]) UpByPK(e E) (c int64, err error) {
	var result sql.Result
	sqlStr, ps := o.getUpdateBuilder(e)
	queryLog("OpsForUpdate.UpByPK", sqlStr, ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	queryErrorLog(err, "OpsForUpdate.UpByPK", sqlStr, ps)
	if result != nil {
		c, _ = result.RowsAffected()
	}
	return
}
