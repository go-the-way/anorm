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

type insertOperation[E Entity] struct {
	orm           *Orm[E]
	ignoreColumns []sg.C
}

func newInsertOperation[E Entity](o *Orm[E]) *insertOperation[E] {
	return &insertOperation[E]{orm: o, ignoreColumns: make([]sg.C, 0)}
}

// Ignore add ignore when inserts
func (o *insertOperation[E]) Ignore(columns ...sg.C) *insertOperation[E] {
	o.ignoreColumns = append(o.ignoreColumns, columns...)
	return o
}

func (o *insertOperation[E]) getIgnoreMap() map[string]struct{} {
	ignoreMap := entityInsertIgnoreMap[getEntityPkgName(o.orm.entity)]
	for _, c := range o.ignoreColumns {
		ignoreMap[string(c)] = struct{}{}
	}
	return ignoreMap
}

func (o *insertOperation[E]) getInsertBuilder(entities ...E) (string, []any) {
	fields := entityFieldMap[getEntityPkgName(o.orm.entity)]
	fieldColumnMap := entityFieldColumnMap[getEntityPkgName(o.orm.entity)]
	ignoreMap := o.getIgnoreMap()
	builder := sg.InsertBuilder()
	for i, entity := range entities {
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

		if len(entities) == 1 {
			builder.Value(sg.NewJoiner(argGes, ", ", "", "", false))
			continue
		}

		if i == 0 {
			builder.Value(sg.NewJoiner(argGes, ", ", "", ")", false))
		} else if i == len(entities)-1 {
			builder.Value(sg.NewJoiner(argGes, ", ", "(", "", false))
		} else {
			builder.Value(sg.NewJoiner(argGes, ", ", "(", ")", false))
		}
	}
	return builder.Table(o.orm.table()).Build()
}

func (o *insertOperation[E]) Exec(entity E) error {
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getInsertBuilder(entity)
	queryLog("OpsForInsert.Exec", sqlStr, ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	queryErrorLog("OpsForInsert.Exec", sqlStr, ps, err)
	if err != nil {
		return err
	}
	lastInsertId := int64(0)
	if result != nil {
		lastInsertId, _ = result.LastInsertId()
	}
	if lastInsertId > 0 {
		pks := entityPKMap[getEntityPkgName(entity)]
		if pks != nil && len(pks) == 1 {
			pkField := entityColumnFieldMap[getEntityPkgName(entity)][pks[0]]
			if pkField != "" {
				value := reflect.ValueOf(entity).Elem().FieldByName(pkField)
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

func (o *insertOperation[E]) ExecList(ignoreError bool, entities ...E) error {
	var err error
	for _, m := range entities {
		err = o.Exec(m)
		if !ignoreError && err != nil {
			return err
		}
	}
	return err
}

func (o *insertOperation[E]) ExecBatch(entities ...E) (int64, error) {
	if len(entities) <= 0 {
		return 0, nil
	}
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getInsertBuilder(entities...)
	queryLog("OpsForInsert.ExecBatch", sqlStr, ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	queryErrorLog("OpsForInsert.ExecBatch", sqlStr, ps, err)
	if err != nil {
		return 0, err
	}
	ra := int64(0)
	if result != nil {
		ra, _ = result.RowsAffected()
	}
	return ra, err
}
