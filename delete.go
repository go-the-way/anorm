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

type _Delete struct {
	orm    *orm
	wheres []sg.Ge
}

func newDelete(o *orm) *_Delete {
	return &_Delete{orm: o, wheres: make([]sg.Ge, 0)}
}

func (o *_Delete) IfWhere(cond bool, wheres ...sg.Ge) *_Delete {
	if cond {
		return o.Where(wheres...)
	}
	return o
}

func (o *_Delete) Where(wheres ...sg.Ge) *_Delete {
	o.wheres = append(o.wheres, wheres...)
	return o
}

func (o *_Delete) appendWhereGes(model Model) {
	o.wheres = append(o.wheres, o.orm.getWhereGes(model)...)
}

func (o *_Delete) getDeleteBuilder(model Model) (string, []interface{}) {
	o.appendWhereGes(model)
	builder := sg.DeleteBuilder()
	return builder.Where(sg.AndGroup(o.wheres...)).From(o.orm.table()).Build()
}

func (o *_Delete) Exec(model Model) (int64, error) {
	var (
		result sql.Result
		err    error
	)
	sqlStr, ps := o.getDeleteBuilder(model)
	debug("delete.Exec[sql:%v ps:%v]", sqlStr, ps)
	execDeleteHookersBefore(model, &sqlStr, &ps)
	if o.orm.openTx {
		result, err = o.orm.tx.Exec(sqlStr, ps...)
	} else {
		result, err = o.orm.db.Exec(sqlStr, ps...)
	}
	execDeleteHookersAfter(model, sqlStr, ps, err)
	ra := int64(0)
	if result != nil {
		ra, _ = result.RowsAffected()
	}
	return ra, err
}
