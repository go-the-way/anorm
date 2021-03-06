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

import "testing"

func TestTxManager(t *testing.T) {
	txm := NewTxManager()
	if tx, err := testDB.Begin(); err != nil {
		t.Error("TestTxManager failed")
	} else {
		txm.Join(tx)
		_ = tx.Rollback()
	}

	if err := txm.Commit(); err == nil {
		t.Error("TestTxManager failed")
	}

	if err := txm.Rollback(); err == nil {
		t.Error("TestTxManager failed")
	}
}
