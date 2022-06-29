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
	"errors"
	"fmt"
	"github.com/go-the-way/anorm/pagination"
	"github.com/go-the-way/sg"
)

type selectOperation[E Entity] struct {
	orm       *Orm[E]
	countJoin bool
	join      bool
	columns   []sg.Ge
	wheres    []sg.Ge
	orderBys  []sg.Ge
}

func newSelectOperation[E Entity](o *Orm[E]) *selectOperation[E] {
	return &selectOperation[E]{orm: o, columns: make([]sg.Ge, 0), wheres: make([]sg.Ge, 0), orderBys: make([]sg.Ge, 0)}
}

// CountJoin enable join query
func (o *selectOperation[E]) CountJoin() *selectOperation[E] {
	o.countJoin = true
	return o
}

// Join enable join query
func (o *selectOperation[E]) Join() *selectOperation[E] {
	o.join = true
	return o
}

func (o *selectOperation[E]) getColumns() []sg.Ge {
	columnGes := make([]sg.Ge, 0)
	columns, cHave := entityColumnMap[getEntityPkgName(o.orm.entity)]
	joinRefMap, jHave := entityJoinRefMap[getEntityPkgName(o.orm.entity)]
	var (
		nullFieldMap map[string]*NullField
		nullHave     bool
	)
	if o.join {
		nullFieldMap, nullHave = entityJoinNullFieldMap[getEntityPkgName(o.orm.entity)]
	} else {
		nullFieldMap, nullHave = entityNullFieldMap[getEntityPkgName(o.orm.entity)]
	}
	if cHave {
		for _, c := range columns {
			fieldName, fieldHave := entityColumnFieldMap[getEntityPkgName(o.orm.entity)][c]
			if fieldHave && (!jHave || joinRefMap[fieldName] == nil) {
				// add: null function
				if nullHave && nullFieldMap[fieldName] != nil {
					jn := nullFieldMap[fieldName]
					// IFNULL(rel1.name, 'defaultVal') AS alias
					if jn.DefaultArg {
						columnGes = append(columnGes, sg.Alias(newFuncGe(fmt.Sprintf("%s(%s, ?)", jn.FuncName, "t."+c), jn.DefaultVal), fieldName))
					} else {
						columnGes = append(columnGes, sg.Alias(newFuncGe(fmt.Sprintf("%s(%s, %v)", jn.FuncName, "t."+c, jn.DefaultVal)), fieldName))
					}
				} else {
					// rel_table.rel_column AS RelColumn
					columnGes = append(columnGes, sg.Alias(sg.C("t."+c), fieldName))
				}
			}
		}
	}
	return columnGes
}

type funcGe struct {
	define string
	args   []any
}

func newFuncGe(define string, args ...any) *funcGe {
	return &funcGe{define: define, args: args}
}

func (f *funcGe) SQL() (string, []interface{}) {
	// IFNULL(t.name, ?), []any{""}
	return f.define, f.args
}

func (o *selectOperation[E]) getJoinRef() ([]sg.Ge, []sg.Ge) {
	columnGes := make([]sg.Ge, 0)
	joinGs := make([]sg.Ge, 0)
	refCount := 1
	refTableMap := make(map[string]string, 0)
	if joinRefMap, have := entityJoinRefMap[getEntityPkgName(o.orm.entity)]; have && o.join {
		nullFieldMap, nullHave := entityJoinNullFieldMap[getEntityPkgName(o.orm.entity)]
		// append join column
		for k, v := range joinRefMap {
			relAlias, joined := refTableMap[v.RelTable]
			if relAlias == "" {
				relAlias = fmt.Sprintf("rel%d", refCount)
				refTableMap[v.RelTable] = relAlias
				refCount++
			}

			// add: null function
			if nullHave && nullFieldMap[k] != nil {
				jn := nullFieldMap[k]
				// IFNULL(rel1.name, 'defaultVal') AS alias
				if jn.DefaultArg {
					columnGes = append(columnGes, sg.Alias(newFuncGe(fmt.Sprintf("%s(%s, ?)", jn.FuncName, relAlias+"."+v.RelName), jn.DefaultVal), k))
				} else {
					columnGes = append(columnGes, sg.Alias(newFuncGe(fmt.Sprintf("%s(%s, %v)", jn.FuncName, relAlias+"."+v.RelName, jn.DefaultVal)), k))
				}
			} else {
				// rel_table.rel_column AS RelColumn
				columnGes = append(columnGes, sg.Alias(sg.C(relAlias+"."+v.RelName), k))
			}

			// LEFT JOIN rel_table ON rel_table.rel_id = t.self_id
			if !joined {
				joinGs = append(joinGs, sg.NewJoiner(
					[]sg.Ge{sg.C(v.Type),
						sg.C("JOIN"),
						sg.Alias(sg.T(v.RelTable), relAlias),
						sg.C("ON"),
						sg.C(relAlias + "." + v.RelID),
						sg.C("="),
						sg.C("t." + v.SelfColumn)},
					" ", "", "", false),
				)
			}
		}
	}
	return columnGes, joinGs
}

func (o *selectOperation[E]) getTableName() sg.Ge {
	return sg.T(entityTableMap[getEntityPkgName(o.orm.entity)])
}

// IfWhere if cond is true append wheres
func (o *selectOperation[E]) IfWhere(cond bool, wheres ...sg.Ge) *selectOperation[E] {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

// Where append wheres
func (o *selectOperation[E]) Where(wheres ...sg.Ge) *selectOperation[E] {
	o.wheres = append(o.wheres, wheres...)
	return o
}

// OrderBy append OrderBys
func (o *selectOperation[E]) OrderBy(orderBys ...sg.Ge) *selectOperation[E] {
	o.orderBys = append(o.orderBys, orderBys...)
	return o
}

func (o *selectOperation[E]) appendWhereGes(entity E) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(entity)...)
}

var (
	ErrQueryTooManyResult = errors.New("query one return too many result")
)

// QueryOne select for one return
//
// Params:
//
// - entity: the orm wrapper entity
//
// Returns:
//
// - e: return one entity
//
// - err: exec error
//
func (o *selectOperation[E]) QueryOne(entity E) (e E, err error) {
	if es, err2 := o.Query(entity); err2 != nil {
		err = err2
	} else if len(es) > 1 {
		err = ErrQueryTooManyResult
	} else if len(es) == 1 {
		e = es[0]
	}
	return
}

// Query select for entities
//
// Params:
//
// - entity: the orm wrapper entity
//
// Returns:
//
// - entities: entities
//
// - err: exec error
//
func (o *selectOperation[E]) Query(entity E) (entities []E, err error) {
	o.appendWhereGes(entity)
	selectBuilder := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(sg.Alias(o.getTableName(), "t")).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...)
	refColumns, refJoins := o.getJoinRef()
	if len(refColumns) > 0 {
		selectBuilder.Select(refColumns...)
	}
	if len(refJoins) > 0 {
		selectBuilder.Join(sg.NewJoiner(refJoins, " ", "", "", false))
	}
	sqlStr, ps := selectBuilder.Build()
	queryLog("OpsForSelect.Query", sqlStr, ps)
	var rows *sql.Rows
	if rows, err = o.orm.db.Query(sqlStr, ps...); err != nil {
		queryErrorLog(err, "OpsForSelect.Query", sqlStr, ps)
		return
	}
	return ScanStruct(rows, o.orm.entity, entityComplete[getEntityPkgName(entity)])
}

// QueryPage select for page
//
// Params:
//
// - entity: the orm wrapper entity
//
// - pager: the pager see pkg pagination
//
// - offset: start index
//
// - size: select size
//
// Returns:
//
// - entities: entities
//
// - total: total rows size
//
// - err: exec error
//
func (o *selectOperation[E]) QueryPage(entity E, pager pagination.Pager, offset, size int) (entities []E, total int64, err error) {
	refColumns, refJoins := o.getJoinRef()
	sc := o.orm.OpsForSelectCount()
	sc.Where(o.wheres...)
	if o.countJoin && len(refJoins) > 0 {
		sc.join(sg.NewJoiner(refJoins, " ", "", "", false))
	}
	total, err = sc.Exec(entity)
	if err != nil {
		return
	}
	if total <= 0 {
		return make([]E, 0), 0, nil
	}
	selectBuilder := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(sg.Alias(o.getTableName(), "t")).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...)
	if len(refColumns) > 0 {
		selectBuilder.Select(refColumns...)
	}
	if len(refJoins) > 0 {
		selectBuilder.Join(sg.NewJoiner(refJoins, " ", "", "", false))
	}
	sqlStr, ps := selectBuilder.Build()
	sqlStr, pps := pager.Page(sqlStr, offset, size)
	ps = append(ps, pps...)
	queryLog("OpsForSelect.QueryPage", sqlStr, ps)
	var rows *sql.Rows
	if rows, err = o.orm.db.Query(sqlStr, ps...); err != nil {
		queryErrorLog(err, "OpsForSelect.QueryPage", sqlStr, ps)
		return
	}
	entities, err = ScanStruct(rows, o.orm.entity, entityComplete[getEntityPkgName(entity)])
	return
}
