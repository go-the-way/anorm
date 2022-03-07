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
	"github.com/go-the-way/sg"
)

type _Select struct {
	orm      *orm
	columns  []sg.Ge
	wheres   []sg.Ge
	orderBys []sg.Ge
}

func newSelect(o *orm) *_Select {
	return &_Select{orm: o, columns: make([]sg.Ge, 0), wheres: make([]sg.Ge, 0), orderBys: make([]sg.Ge, 0)}
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
	for _, c := range columns {
		if len(cm) > 0 {
			if _, have := cm[sg.C(c)]; !have {
				continue
			}
		}
		fieldName := modelColumnFieldMap[getModelPkgName(o.orm.model)][c]
		columnGes = append(columnGes, sg.Alias(sg.C(c), fieldName))
	}
	return columnGes
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
	sqlStr, ps := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...).
		Build()
	debug("List[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	rows, err := o.orm.db.Query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	if err != nil {
		return nil, err
	}
	return scanStruct(rows, o.orm.model)
}

func (o *_Select) ExecPage(model Model, pager Pagination, offset, size int) ([]Model, int64, error) {
	if pager == nil {
		panic("anorm: the pager is nil")
	}
	var err error
	sc := o.orm.SelectCount()
	sc.wheres = append(sc.wheres, o.wheres...)
	c, err := sc.Exec(model)
	debug("[ExecPage]count [%d]", c)
	if err != nil {
		return nil, 0, err
	}
	if c <= 0 {
		return make([]Model, 0), 0, nil
	}
	sqlStr, ps := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...).
		Build()
	debug("[ExecPage]originalSql [%s] originalPs [%v]", sqlStr, ps)
	sqlStr, pps := pager.Page(sqlStr, offset, size)
	ps = append(ps, pps...)
	debug("[ExecPage]currentSql [%s] currentPs [%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	rows, err := o.orm.db.Query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	if err != nil {
		return nil, 0, nil
	}
	models, err := scanStruct(rows, o.orm.model)
	return models, c, err
}
