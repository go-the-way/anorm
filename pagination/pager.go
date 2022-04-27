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

package pagination

// Pager define pagination interface
type Pager interface {
	Page(sql string, offset, size int) (sqlStr string, args []any)
}

var (
	// MySql define MySQL pager
	//
	// SELECT t.* FROM (...) LIMIT ?, ?
	//
	// $1: offset
	//
	// $2: size
	//
	MySql = &mysql{}
	// Pg define pgsql pager
	//
	// SELECT t.* FROM (...) LIMIT ? OFFSET ?
	//
	// $1: size
	//
	// $2: offset
	Pg = &pg{}
	// SqlServer tells orm how to use ROW_NUMBER()
	//
	// SELECT t.* FROM (
	// 		SELECT t.*, ROW_NUMBER() over(order by id asc) as rn FROM (...) as t
	// ) as t WHERE t.rn BETWEEN ? AND ?
	//
	//
	// $1: offset + 1
	//
	// $2: offset + size
	SqlServer = func(orderBy string) *sqlServer { return &sqlServer{orderBy} }
)
