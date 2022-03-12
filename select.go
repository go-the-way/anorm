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
	"fmt"
	"github.com/go-the-way/sg"
)

type _Select struct {
	orm      *orm
	join     bool
	columns  []sg.Ge
	wheres   []sg.Ge
	orderBys []sg.Ge
}

func newSelect(o *orm) *_Select {
	return &_Select{orm: o, columns: make([]sg.Ge, 0), wheres: make([]sg.Ge, 0), orderBys: make([]sg.Ge, 0)}
}

func (o *_Select) Join() *_Select {
	o.join = true
	return o
}

func (o *_Select) getColumns() []sg.Ge {
	// fixed: have no alias for field
	cm := make(map[sg.Ge]struct{}, 0)
	if o.columns != nil {
		for _, c := range o.columns {
			cm[c] = struct{}{}
		}
	}
	columnGes := make([]sg.Ge, 0)
	columns := modelColumnMap[getModelPkgName(o.orm.model)]
	joinRefs := modelJoinRefMap[getModelPkgName(o.orm.model)]
	for _, c := range columns {
		if len(cm) > 0 {
			if _, have := cm[sg.C(c)]; !have {
				continue
			}
		}
		fieldName := modelColumnFieldMap[getModelPkgName(o.orm.model)][c]
		if joinRefs == nil || joinRefs[fieldName] == nil {
			columnGes = append(columnGes, sg.Alias(sg.C("t."+c), fieldName))
		}
	}
	return columnGes
}

func (o *_Select) getJoinRef() ([]sg.Ge, []sg.Ge) {
	columnGes := make([]sg.Ge, 0)
	joinGs := make([]sg.Ge, 0)
	refCount := 1
	if joinRefMap, have := modelJoinRefMap[getModelPkgName(o.orm.model)]; have && o.join {
		// append join column
		for k, v := range joinRefMap {
			relAlias := fmt.Sprintf("rel%d", refCount)
			// rel_table.rel_column AS RelColumn
			columnGes = append(columnGes, sg.Alias(sg.C(relAlias+"."+v.RelName), k))
			// LEFT JOIN rel_table ON rel_table.rel_id = t.self_id
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
			refCount++
		}
	}
	return columnGes, joinGs
}

func (o *_Select) getTableName() sg.Ge {
	return sg.T(modelTableMap[getModelPkgName(o.orm.model)])
}

func (o *_Select) IfWhere(cond bool, wheres ...sg.Ge) *_Select {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *_Select) Where(wheres ...sg.Ge) *_Select {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *_Select) OrderBy(orderBys ...sg.Ge) *_Select {
	o.orderBys = append(o.orderBys, orderBys...)
	return o
}

func (o *_Select) appendWhereGes(model Model) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(model)...)
}

func (o *_Select) Exec(model Model) ([]Model, error) {
	o.appendWhereGes(model)
	selectBuilder := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(sg.Alias(o.getTableName(), "t")).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...)
	refColumns, refJoins := o.getJoinRef()
	if len(refColumns) > 0 && len(refColumns) == len(refJoins) {
		selectBuilder.Select(refColumns...)
		selectBuilder.Join(sg.NewJoiner(refJoins, " ", "", "", false))
	}
	sqlStr, ps := selectBuilder.Build()
	if debug() {
		(&execLog{"Select.Exec", sqlStr, ps}).Log()
	}
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	rows, err := o.orm.db.Query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	if err != nil {
		return nil, err
	}
	return scanStruct(rows, o.orm.model)
}

func (o *_Select) ExecPage(model Model, pager Pagination, offset, size int) ([]Model, int64, error) {
	var err error
	sc := o.orm.SelectCount()
	sc.wheres = append(sc.wheres, o.wheres...)
	c, err := sc.Exec(model)
	if err != nil {
		return nil, 0, err
	}
	if c <= 0 {
		return make([]Model, 0), 0, nil
	}
	selectBuilder := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(sg.Alias(o.getTableName(), "t")).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...)
	refColumns, refJoins := o.getJoinRef()
	if len(refColumns) > 0 && len(refColumns) == len(refJoins) {
		selectBuilder.Select(refColumns...)
		selectBuilder.Join(sg.NewJoiner(refJoins, " ", "", "", false))
	}
	sqlStr, ps := selectBuilder.Build()
	sqlStr, pps := pager.Page(sqlStr, offset, size)
	ps = append(ps, pps...)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	if debug() {
		(&execLog{"Select.ExecPage", sqlStr, ps}).Log()
	}
	rows, err := o.orm.db.Query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	if err != nil {
		return nil, 0, err
	}
	models, err := scanStruct(rows, o.orm.model)
	return models, c, err
}
