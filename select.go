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
	"fmt"
	"github.com/go-the-way/anorm/pagination"
	"github.com/go-the-way/sg"
)

type selectOperation[E Entity] struct {
	orm      *Orm[E]
	join     bool
	columns  []sg.Ge
	wheres   []sg.Ge
	orderBys []sg.Ge
}

func newSelectOperation[E Entity](o *Orm[E]) *selectOperation[E] {
	return &selectOperation[E]{orm: o, columns: make([]sg.Ge, 0), wheres: make([]sg.Ge, 0), orderBys: make([]sg.Ge, 0)}
}

// Join enable join query
func (o *selectOperation[E]) Join() *selectOperation[E] {
	o.join = true
	return o
}

func (o *selectOperation[E]) getColumns() []sg.Ge {
	columnGes := make([]sg.Ge, 0)
	columns := entityColumnMap[getEntityPkgName(o.orm.entity)]
	joinRefs := entityJoinRefMap[getEntityPkgName(o.orm.entity)]
	for _, c := range columns {
		fieldName := entityColumnFieldMap[getEntityPkgName(o.orm.entity)][c]
		if joinRefs == nil || joinRefs[fieldName] == nil {
			columnGes = append(columnGes, sg.Alias(sg.C("t."+c), fieldName))
		}
	}
	return columnGes
}

func (o *selectOperation[E]) getJoinRef() ([]sg.Ge, []sg.Ge) {
	columnGes := make([]sg.Ge, 0)
	joinGs := make([]sg.Ge, 0)
	refCount := 1
	refTableMap := make(map[string]string, 0)
	if joinRefMap, have := entityJoinRefMap[getEntityPkgName(o.orm.entity)]; have && o.join {
		// append join column
		for k, v := range joinRefMap {
			relAlias, joined := refTableMap[v.RelTable]
			if relAlias == "" {
				relAlias = fmt.Sprintf("rel%d", refCount)
				refTableMap[v.RelTable] = relAlias
				refCount++
			}
			// rel_table.rel_column AS RelColumn
			columnGes = append(columnGes, sg.Alias(sg.C(relAlias+"."+v.RelName), k))
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

// ExecOne select for one return
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
func (o *selectOperation[E]) ExecOne(entity E) (e E, err error) {
	if es, err2 := o.Exec(entity); err2 == nil && len(es) > 0 {
		e = es[0]
	} else {
		err = err2
	}
	return
}

// Exec select for entities
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
func (o *selectOperation[E]) Exec(entity E) (entities []E, err error) {
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
	queryLog("OpsForSelect.Exec", sqlStr, ps)
	var rows *sql.Rows
	if rows, err = o.orm.db.Query(sqlStr, ps...); err != nil {
		queryErrorLog("OpsForSelect.Exec", sqlStr, ps, err)
		return
	}
	return scanStruct(rows, o.orm.entity)
}

// ExecPage select for page
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
func (o *selectOperation[E]) ExecPage(entity E, pager pagination.Pager, offset, size int) (entities []E, total int64, err error) {
	sc := o.orm.OpsForSelectCount()
	sc.wheres = append(sc.wheres, o.wheres...)
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
	refColumns, refJoins := o.getJoinRef()
	if len(refColumns) > 0 {
		selectBuilder.Select(refColumns...)
	}
	if len(refJoins) > 0 {
		selectBuilder.Join(sg.NewJoiner(refJoins, " ", "", "", false))
	}
	sqlStr, ps := selectBuilder.Build()
	sqlStr, pps := pager.Page(sqlStr, offset, size)
	ps = append(ps, pps...)
	queryLog("OpsForSelect.ExecPage", sqlStr, ps)
	var rows *sql.Rows
	if rows, err = o.orm.db.Query(sqlStr, ps...); err != nil {
		queryErrorLog("OpsForSelect.ExecPage", sqlStr, ps, err)
		return
	}
	entities, err = scanStruct(rows, o.orm.entity)
	return
}
