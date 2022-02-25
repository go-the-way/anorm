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
	"database/sql"
	"errors"
	"github.com/go-the-way/sg"
)

var (
	errPaginationIsNil = errors.New("the pagination is nil")
)

type ormSelect struct {
	orm      *orm
	columns  []sg.Ge
	wheres   []sg.Ge
	orderBys []sg.Ge
}

func (o *ormSelect) getColumns() []sg.Ge {
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

func (o *ormSelect) getTableName() sg.Ge {
	return sg.T(modelTableMap[getModelPkgName(o.orm.model)])
}

func (o *ormSelect) NotIfWhere(cond bool, wheres ...sg.Ge) *ormSelect {
	return o.IfWhere(!cond, wheres...)
}

func (o *ormSelect) IfWhere(cond bool, wheres ...sg.Ge) *ormSelect {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *ormSelect) Where(wheres ...sg.Ge) *ormSelect {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *ormSelect) OrderBy(orderBys ...sg.Ge) *ormSelect {
	o.orderBys = append(o.orderBys, orderBys...)
	return o
}

func (o *ormSelect) appendWhereGes(model Model) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(model)...)
}

func (o *ormSelect) Count(model Model) (int64, error) {
	o.appendWhereGes(model)
	sqlStr, ps := sg.SelectBuilder().
		Select(sg.Alias(sg.C("count(0)"), "c")).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		Build()
	debug("Count[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(model, &sqlStr, &ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, "_"), o.orm.openTX, o.orm.autoCommit)
	row := etr.queryRow(sqlStr, ps...)
	execSelectHookersAfter(model, sqlStr, ps, row.Err())
	if row.Err() != nil {
		return 0, row.Err()
	}
	count := int64(0)
	err := row.Scan(&count)
	return count, err
}

func (o *ormSelect) QueryRow(model Model) *sql.Row {
	o.appendWhereGes(model)
	sqlStr, ps := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...).
		Build()
	debug("QueryRow[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(model, &sqlStr, &ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, "_"), o.orm.openTX, o.orm.autoCommit)
	row := etr.queryRow(sqlStr, ps...)
	execSelectHookersAfter(model, sqlStr, ps, row.Err())
	return row
}

func (o *ormSelect) Query(model Model) (*sql.Rows, error) {
	o.appendWhereGes(model)
	sqlStr, ps := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...).
		Build()
	execSelectHookersBefore(model, &sqlStr, &ps)
	debug("Query[sql:%v ps:%v]", sqlStr, ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, "_"), o.orm.openTX, o.orm.autoCommit)
	rows, err := etr.query(sqlStr, ps...)
	execSelectHookersAfter(model, sqlStr, ps, err)
	return rows, err
}

func (o *ormSelect) QuerySql(sqlStr string, ps ...interface{}) (*sql.Rows, error) {
	debug("QuerySql[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, "_"), o.orm.openTX, o.orm.autoCommit)
	rows, err := etr.query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	return rows, err
}

func (o *ormSelect) QuerySqlWithDS(ds string, sqlStr string, ps ...interface{}) (*sql.Rows, error) {
	debug("QuerySqlWithDS[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, ds), o.orm.openTX, o.orm.autoCommit)
	rows, err := etr.query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	return rows, err
}

func (o *ormSelect) First(model Model) (Model, error) {
	models, err := o.List(model)
	if err != nil {
		return nil, err
	}
	if len(models) <= 0 {
		return nil, nil
	}
	return models[0], err
}

func (o *ormSelect) List(model Model) ([]Model, error) {
	o.appendWhereGes(model)
	sqlStr, ps := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...).
		Build()
	debug("List[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, "_"), o.orm.openTX, o.orm.autoCommit)
	rows, err := etr.query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	if err != nil {
		return nil, err
	}
	return scanStruct(rows, o.orm.model)
}

func (o *ormSelect) Pagination(model Model, pager Pagination, offset, size int) ([]Model, int64, error) {
	if pager == nil {
		panic(errPaginationIsNil)
	}
	var err error
	count, err := o.Count(model)
	if err != nil {
		return nil, 0, err
	}
	debug("Pagination count [%d]", count)
	sqlStr, ps := sg.SelectBuilder().
		Select(o.getColumns()...).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		OrderBy(o.orderBys...).
		Build()
	debug("Pagination originalSql [%s] originalPs [%v]", sqlStr, ps)
	sqlStr, pps := pager.Page(sqlStr, offset, size)
	ps = append(ps, pps...)
	debug("Pagination currentSql [%s] currentPs [%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	etr := newExecutor(dsMap.dsHold(o.orm.db, "_"), o.orm.openTX, o.orm.autoCommit)
	rows, err := etr.query(sqlStr, ps...)
	execSelectHookersAfter(o.orm.model, sqlStr, ps, err)
	if err != nil {
		return nil, 0, nil
	}
	models, err := scanStruct(rows, o.orm.model)
	return models, count, err
}
