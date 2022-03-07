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

type _SelectCount struct {
	orm    *orm
	wheres []sg.Ge
}

func newSelectCount(o *orm) *_SelectCount {
	return &_SelectCount{orm: o, wheres: make([]sg.Ge, 0)}
}

func (o *_SelectCount) IfWhere(cond bool, wheres ...sg.Ge) *_SelectCount {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *_SelectCount) Where(wheres ...sg.Ge) *_SelectCount {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *_SelectCount) getTableName() sg.Ge {
	return sg.T(modelTableMap[getModelPkgName(o.orm.model)])
}

func (o *_SelectCount) appendWhereGes(model Model) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(model)...)
}

func (o *_SelectCount) Exec(model Model) (int64, error) {
	o.appendWhereGes(model)
	sqlStr, ps := sg.SelectBuilder().
		Select(sg.Alias(sg.C("count(0)"), "c")).
		From(o.getTableName()).
		Where(sg.AndGroup(o.wheres...)).
		Build()
	debug("selectCount.Exec[sql:%v ps:%v]", sqlStr, ps)
	execSelectHookersBefore(o.orm.model, &sqlStr, &ps)
	row := o.orm.db.QueryRow(sqlStr, ps...)
	execSelectHookersAfter(model, sqlStr, ps, row.Err())
	if err := row.Err(); err != nil {
		return 0, err
	}
	count := int64(0)
	err := row.Scan(&count)
	return count, err
}
