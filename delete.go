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
)

type (
	DeleteOperation[E Entity] interface {
		BeginTx(txm *TxManager, options ...*sql.TxOptions) error
		IfWhere(cond bool, wheres ...sg.Ge) DeleteOperation[E]
		IfOnlyWhere(cond bool, wheres ...sg.Ge) DeleteOperation[E]
		Where(wheres ...sg.Ge) DeleteOperation[E]
		OnlyWhere(wheres ...sg.Ge) DeleteOperation[E]
		Del(e E) (count int64, err error)
	}
	deleteOperation[E Entity] struct {
		orm                *Orm[E]
		wheres, onlyWheres []sg.Ge
	}
)

func Delete[E Entity](e E) DeleteOperation[E] {
	return New(e).OpsForDelete()
}

func DeleteWithDs[E Entity](e E, ds string) DeleteOperation[E] {
	return NewWithDS(e, ds).OpsForDelete()
}

func newsDeleteOperation[E Entity](o *Orm[E]) *deleteOperation[E] {
	return &deleteOperation[E]{orm: o, wheres: make([]sg.Ge, 0), onlyWheres: make([]sg.Ge, 0)}
}

func (o *deleteOperation[E]) BeginTx(txm *TxManager, options ...*sql.TxOptions) error {
	return o.orm.BeginTx(txm, options...)
}

// IfWhere if cond is true, append wheres
func (o *deleteOperation[E]) IfWhere(cond bool, wheres ...sg.Ge) DeleteOperation[E] {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

// IfOnlyWhere if cond is true, append only wheres
func (o *deleteOperation[E]) IfOnlyWhere(cond bool, wheres ...sg.Ge) DeleteOperation[E] {
	if cond {
		return o.OnlyWhere(wheres...)
	}
	return o
}

// Where append wheres
func (o *deleteOperation[E]) Where(wheres ...sg.Ge) DeleteOperation[E] {
	o.wheres = append(o.wheres, wheres...)
	return o
}

// OnlyWhere append only wheres
func (o *deleteOperation[E]) OnlyWhere(wheres ...sg.Ge) DeleteOperation[E] {
	o.onlyWheres = append(o.onlyWheres, wheres...)
	return o
}

func (o *deleteOperation[E]) appendWhereGes(entity E) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(entity)...)
}

func (o *deleteOperation[E]) getDeleteBuilder(entity E) (string, []any) {
	builder := sg.DeleteBuilder().From(o.orm.table())
	if len(o.onlyWheres) > 0 {
		builder.Where(sg.AndGroup(o.onlyWheres...))
	} else {
		o.appendWhereGes(entity)
		builder.Where(sg.AndGroup(o.wheres...))
	}
	return builder.Build()
}

// Del delete entities
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
func (o *deleteOperation[E]) Del(e E) (int64, error) {
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getDeleteBuilder(e)
	queryLog("OpsForDelete.Del", sqlStr, ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	queryErrorLog(err, "OpsForDelete.Del", sqlStr, ps)
	ra := int64(0)
	if result != nil {
		ra, err = result.RowsAffected()
		queryErrorLog(err, "OpsForDelete.Del", sqlStr, ps)
	}
	return ra, err
}
