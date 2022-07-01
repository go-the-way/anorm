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
	InsertOperation[E Entity] interface {
		BeginTx(txm *TxManager, options ...*sql.TxOptions) error
		Ignore(cs ...sg.C) InsertOperation[E]
		One(e E) error
		List(ignoreError bool, es ...E) error
		Batch(es ...E) (count int64, err error)
	}
	insertOperation[E Entity] struct {
		orm           *Orm[E]
		ignoreColumns []sg.C
	}
)

func Insert[E Entity](e E) InsertOperation[E] {
	return New(e).OpsForInsert()
}

func InsertWithDs[E Entity](e E, ds string) InsertOperation[E] {
	return NewWithDS(e, ds).OpsForInsert()
}

func newInsertOperation[E Entity](o *Orm[E]) *insertOperation[E] {
	return &insertOperation[E]{orm: o, ignoreColumns: make([]sg.C, 0)}
}

func (o *insertOperation[E]) BeginTx(txm *TxManager, options ...*sql.TxOptions) error {
	return o.orm.BeginTx(txm, options...)
}

// Ignore add ignore when inserts
func (o *insertOperation[E]) Ignore(cs ...sg.C) InsertOperation[E] {
	o.ignoreColumns = append(o.ignoreColumns, cs...)
	return o
}

func (o *insertOperation[E]) getIgnoreMap() map[string]struct{} {
	ignoreMap := entityInsertIgnoreMap[getEntityPkgName(o.orm.entity)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *insertOperation[E]) getInsertBuilder(es ...E) (string, []any) {
	fields := entityFieldMap[getEntityPkgName(o.orm.entity)]
	fieldColumnMap := entityFieldColumnMap[getEntityPkgName(o.orm.entity)]
	ignoreMap := o.getIgnoreMap()
	builder := sg.InsertBuilder()
	for i, entity := range es {
		rt := reflect.ValueOf(entity).Elem()
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

		if len(es) == 1 {
			builder.Value(sg.NewJoiner(argGes, ", ", "", "", false))
			continue
		}

		if i == 0 {
			builder.Value(sg.NewJoiner(argGes, ", ", "", ")", false))
		} else if i == len(es)-1 {
			builder.Value(sg.NewJoiner(argGes, ", ", "(", "", false))
		} else {
			builder.Value(sg.NewJoiner(argGes, ", ", "(", ")", false))
		}
	}
	return builder.Table(o.orm.table()).Build()
}

// One exec insert one entity
//
// Params:
//
// - e: the orm wrapper entity
//
// Returns:
//
// - err: exec error
//
func (o *insertOperation[E]) One(e E) error {
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getInsertBuilder(e)
	queryLog("OpsForInsert.Count", sqlStr, ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	queryErrorLog(err, "OpsForInsert.Count", sqlStr, ps)
	if err != nil {
		return err
	}
	lastInsertId := int64(0)
	if result != nil {
		lastInsertId, _ = result.LastInsertId()
	}
	if lastInsertId > 0 {
		pks := entityPKMap[getEntityPkgName(e)]
		if pks != nil && len(pks) == 1 {
			pkField := entityColumnFieldMap[getEntityPkgName(e)][pks[0]]
			if pkField != "" {
				value := reflect.ValueOf(e).Elem().FieldByName(pkField)
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

// List exec list entity, each entity call One(entity Entity) error
//
// Params:
//
// - ignoreError: if true, any insert error does not terminate
//
// - entities: the orm wrapper entity list
//
// Returns:
//
// - err: exec error
//
func (o *insertOperation[E]) List(ignoreError bool, es ...E) error {
	var err error
	for _, e := range es {
		err = o.One(e)
		if !ignoreError && err != nil {
			return err
		}
	}
	return err
}

// Batch exec list entity, use batch mode
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
func (o *insertOperation[E]) Batch(entities ...E) (count int64, err error) {
	if len(entities) <= 0 {
		return 0, nil
	}
	var result sql.Result
	sqlStr, ps := o.getInsertBuilder(entities...)
	queryLog("OpsForInsert.Batch", sqlStr, ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	queryErrorLog(err, "OpsForInsert.Batch", sqlStr, ps)
	if err != nil {
		return 0, err
	}
	ra := int64(0)
	if result != nil {
		ra, _ = result.RowsAffected()
	}
	return ra, err
}
