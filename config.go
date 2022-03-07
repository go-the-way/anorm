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
	"errors"
	"log"
	"os"
)

var (
	errInvalidSetConfig = errors.New(`anorm: invalid config set of nil`)

	Configuration = &Config{
		Migrate:            false,
		TableNameStrategy:  Underline,
		ColumnNameStrategy: Underline,
		Logger:             log.New(os.Stdout, "[anorm] ", log.LstdFlags),
		InsertHookers:      make([]ExecHooker, 0),
		DeleteHookers:      make([]ExecHooker, 0),
		UpdateHookers:      make([]ExecHooker, 0),
		SelectHookers:      make([]ExecHooker, 0),
	}
)

// Config defines configuration struct
type Config struct {
	// Migrate defines migrate create or update table in the database
	Migrate bool
	// TableNameStrategy defines binding table name for Model struct
	TableNameStrategy Strategy
	// TableNameStrategy defines binding table's column name for Model's Fields
	ColumnNameStrategy Strategy
	// Logger defines logger for anorm
	Logger *log.Logger
	// Logger Debug
	Debug         bool
	InsertHookers []ExecHooker
	DeleteHookers []ExecHooker
	UpdateHookers []ExecHooker
	SelectHookers []ExecHooker
}

// SetConfig define Set configuration for anorm
func SetConfig(config *Config) {
	if config == nil {
		panic(errInvalidSetConfig)
	}
	Configuration = config
}

func execInsertHookersBefore(model Model, sqlStr *string, ps *[]interface{}) {
	if hks := Configuration.InsertHookers; hks != nil {
		for _, hk := range hks {
			hk.BeforeExec(model, sqlStr, ps)
		}
	}
}

func execInsertHookersAfter(model Model, sqlStr string, ps []interface{}, err error) {
	if hks := Configuration.InsertHookers; hks != nil {
		for _, hk := range hks {
			hk.AfterExec(model, sqlStr, ps, err)
		}
	}
}

func execDeleteHookersBefore(model Model, sqlStr *string, ps *[]interface{}) {
	if hks := Configuration.DeleteHookers; hks != nil {
		for _, hk := range hks {
			hk.BeforeExec(model, sqlStr, ps)
		}
	}
}

func execDeleteHookersAfter(model Model, sqlStr string, ps []interface{}, err error) {
	if hks := Configuration.DeleteHookers; hks != nil {
		for _, hk := range hks {
			hk.AfterExec(model, sqlStr, ps, err)
		}
	}
}

func execUpdateHookersBefore(model Model, sqlStr *string, ps *[]interface{}) {
	if hks := Configuration.UpdateHookers; hks != nil {
		for _, hk := range hks {
			hk.BeforeExec(model, sqlStr, ps)
		}
	}
}

func execUpdateHookersAfter(model Model, sqlStr string, ps []interface{}, err error) {
	if hks := Configuration.UpdateHookers; hks != nil {
		for _, hk := range hks {
			hk.AfterExec(model, sqlStr, ps, err)
		}
	}
}

func execSelectHookersBefore(model Model, sqlStr *string, ps *[]interface{}) {
	if hks := Configuration.SelectHookers; hks != nil {
		for _, hk := range hks {
			hk.BeforeExec(model, sqlStr, ps)
		}
	}
}

func execSelectHookersAfter(model Model, sqlStr string, ps []interface{}, err error) {
	if hks := Configuration.SelectHookers; hks != nil {
		for _, hk := range hks {
			hk.AfterExec(model, sqlStr, ps, err)
		}
	}
}
