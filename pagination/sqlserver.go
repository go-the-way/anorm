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

import "fmt"

type sqlServer struct{ string }

// Page sqlserver implementation
func (s *sqlServer) Page(sql string, offset, size int) (sqlStr string, args []any) {
	return fmt.Sprintf("SELECT t.* FROM (SELECT _t.*, row_number() over (order by %s) as rn FROM (%s) as _t) WHERE t.rn between ? and ?", s.string, sql), []any{offset + 1, offset + size}
}
