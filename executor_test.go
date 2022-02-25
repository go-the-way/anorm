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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestETRCommit(t *testing.T) {
	{
		etr := newExecutor(testDB, false, false)
		require.Nil(t, etr.Commit())
	}
	{
		etr := newExecutor(testDB, true, false)
		_, _ = etr.exec("create table tmp(id int)")
		_, _ = etr.exec("insert into tmp(id) values (?)", 100)
		_ = etr.Commit()
		r := etr.queryRow("select count(*) from tmp where id = ?", 100)
		etr.Begin()
		id := 0
		_ = r.Scan(&id)
		etr.Begin()
		_, _ = etr.exec("drop table tmp")
		_ = etr.Commit()
		require.Equal(t, 1, id)
	}
}
