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

import (
	"reflect"
	"testing"
)

func TestPager(t *testing.T) {
	{
		sqlStr, args := MySql.Page("haha", 0, 10)
		if sqlStr != "haha LIMIT ?, ?" || !reflect.DeepEqual([]any{0, 10}, args) {
			t.Error("test failed")
		}
	}

	{
		sqlStr, args := SqlServer("id asc").Page("haha", 0, 10)
		if sqlStr != "SELECT t.* FROM (SELECT _t.*, row_number() over (order by id asc) as rn FROM (haha) as _t) WHERE t.rn between ? and ?" || !reflect.DeepEqual([]any{1, 10}, args) {
			t.Error("test failed")
		}
	}

	{
		sqlStr, args := Pg.Page("haha", 0, 10)
		if sqlStr != "haha LIMIT ? OFFSET ?" || !reflect.DeepEqual([]any{10, 0}, args) {
			t.Error("test failed")
		}
	}
}
