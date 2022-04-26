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
	"sync"
)

type txManager struct {
	mu  *sync.Mutex
	txs []*sql.Tx
}

// TxManager simplify tx manager
func TxManager() *txManager {
	return &txManager{mu: &sync.Mutex{}, txs: make([]*sql.Tx, 0)}
}

// Join other *sql.Tx
func (txm *txManager) Join(tx *sql.Tx) {
	txm.mu.Lock()
	defer txm.mu.Unlock()
	txm.txs = append(txm.txs, tx)
}

// Commit all txs
func (txm *txManager) Commit() error {
	txm.mu.Lock()
	defer txm.mu.Unlock()
	for _, tx := range txm.txs {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// Rollback all txs
func (txm *txManager) Rollback() error {
	txm.mu.Lock()
	defer txm.mu.Unlock()
	for _, tx := range txm.txs {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}
	return nil
}
