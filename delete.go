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
	"github.com/go-the-way/sg"
)

type ormDelete struct {
	orm    *orm
	wheres []sg.Ge
}

func (o *ormDelete) NotIfWhere(cond bool, wheres ...sg.Ge) *ormDelete {
	return o.IfWhere(!cond, wheres...)
}

func (o *ormDelete) IfWhere(cond bool, wheres ...sg.Ge) *ormDelete {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *ormDelete) Where(wheres ...sg.Ge) *ormDelete {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *ormDelete) appendWhereGes(model Model) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(model)...)
}

func (o *ormDelete) getDeleteBuilder(model Model) (string, []interface{}) {
	o.appendWhereGes(model)
	builder := sg.DeleteBuilder()
	return builder.Where(sg.AndGroup(o.wheres...)).From(o.orm.table()).Build()
}

func (o *ormDelete) Remove(model Model) error {
	sqlStr, ps := o.getDeleteBuilder(model)
	debug("Remove[sql:%v ps:%v]", sqlStr, ps)
	execDeleteHookersBefore(model, &sqlStr, &ps)
	etr := newExecutorFromOrm(o.orm)
	_, err := etr.exec(sqlStr, ps...)
	execDeleteHookersAfter(model, sqlStr, ps, err)
	return err
}
