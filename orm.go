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
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-the-way/sg"
)

var (
	errModelIsNil              = func(model Model) error { return errors.New(fmt.Sprintf("model is nil")) }
	errModelNotYetRegisterFunc = func(model Model) error {
		return errors.New(fmt.Sprintf("model [%v] not yet register", getModelPkgName(model)))
	}
)

type orm struct {
	model      Model
	tagMap     map[string]*tag
	db         *sql.DB
	tx         *sql.Tx
	openTX     bool
	autoCommit bool
}

// New defines return a new orm from Model model
func New(model Model) *orm {
	return NewWithDS(model, ifEmpty(model.MetaData().DS, "_"), false, false)
}

func NewWithTx(model Model, autoCommit bool) *orm {
	return NewWithDS(model, ifEmpty(model.MetaData().DS, "_"), true, autoCommit)
}

// NewWithDS defines return a new orm using named DS from Model model
func NewWithDS(model Model, ds string, openTx, autoCommit bool) *orm {
	if model == nil {
		panic(errModelIsNil)
	}

	o := &orm{
		model:      model,
		db:         dsMap.required(ds),
		autoCommit: autoCommit,
	}

	if openTx {
		if tx, err := o.db.Begin(); err != nil {
			panic(err)
		} else {
			o.tx = tx
		}
	}

	debug("Model[%v] use DS[%v]", getModelPkgName(model), ds)

	if tagMap, registered := modelTagMap[getModelPkgName(model)]; !registered {
		panic(errModelNotYetRegisterFunc(model))
	} else {
		o.tagMap = tagMap
	}

	return o
}

// table defines return Model's table name
func (o *orm) table() sg.Ge {
	return sg.T(modelTableMap[getModelPkgName(o.model)])
}

// getWhereGes defines return where ges from model
func (o *orm) getWhereGes(model Model) []sg.Ge {
	ges := make([]sg.Ge, 0)
	if model != nil {
		fieldColumnMap := modelFieldColumnMap[getModelPkgName(model)]
		rt := reflect.TypeOf(model).Elem()
		rv := reflect.ValueOf(model).Elem()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			value := rv.Field(i)
			var val interface{}
			switch value.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal := value.Int()
				if intVal != 0 {
					val = intVal
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				intVal := value.Uint()
				if intVal > 0 {
					val = intVal
				}
			case reflect.Float32, reflect.Float64:
				floatVal := value.Float()
				if floatVal > 0 {
					val = floatVal
				}
			case reflect.String:
				strVal := value.String()
				if strVal != "" {
					val = strVal
				}
			case reflect.Struct:
				v := value.Interface()
				if v != nil {
					switch value.Type() {
					case reflect.TypeOf(NullBool(false)):
						if vv := v.(sql.NullBool); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullBoolPtr(false)):
						if vv := v.(*sql.NullBool); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullByte(0)):
						if vv := v.(sql.NullByte); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullBytePtr(0)):
						if vv := v.(*sql.NullByte); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullInt16(0)):
						if vv := v.(sql.NullInt16); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullInt16Ptr(0)):
						if vv := v.(*sql.NullInt16); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullInt32(0)):
						if vv := v.(sql.NullInt32); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullInt32Ptr(0)):
						if vv := v.(*sql.NullInt32); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullInt64(0)):
						if vv := v.(sql.NullInt64); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullInt64Ptr(0)):
						if vv := v.(*sql.NullInt64); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullFloat64(0)):
						if vv := v.(sql.NullFloat64); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullFloat64Ptr(0)):
						if vv := v.(*sql.NullFloat64); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullString("")):
						if vv := v.(sql.NullString); vv.Valid {
							val = vv
						}
					case reflect.TypeOf(NullStringPtr("")):
						if vv := v.(*sql.NullString); vv != nil && vv.Valid {
							val = vv
						}
					case reflect.TypeOf(time.Time{}):
						if vv := v.(time.Time); vv.After(time.Time{}) {
							val = vv
						}
					case reflect.TypeOf(&time.Time{}):
						if vv := v.(*time.Time); vv.After(time.Time{}) {
							val = vv
						}
					case reflect.TypeOf(NullTime(time.Time{})):
						if vv := v.(sql.NullTime); vv.Valid && vv.Time.After(time.Time{}) {
							val = vv
						}
					case reflect.TypeOf(NullTimePtr(time.Time{})):
						if vv := v.(*sql.NullTime); vv.Valid && vv.Time.After(time.Time{}) {
							val = vv
						}
					}
				}
			}
			if val != nil {
				ges = append(ges, sg.Eq(sg.C(fieldColumnMap[field.Name]), val))
			}
		}
	}
	return ges
}

// Select defines return a ormSelect
func (o *orm) Select(columns ...sg.Ge) *ormSelect {
	return &ormSelect{
		orm:      o,
		columns:  columns,
		wheres:   make([]sg.Ge, 0),
		orderBys: make([]sg.Ge, 0),
	}
}

// Insert defines return a ormInsert
func (o *orm) Insert() *ormInsert {
	return &ormInsert{
		orm:           o,
		ignoreColumns: make([]sg.C, 0),
	}
}

// Update defines return a ormUpdate
func (o *orm) Update() *ormUpdate {
	return &ormUpdate{
		orm:           o,
		ignoreColumns: make([]sg.C, 0),
		wheres:        make([]sg.Ge, 0),
	}
}

// Delete defines return a ormDelete
func (o *orm) Delete() *ormDelete {
	return &ormDelete{
		orm:    o,
		wheres: make([]sg.Ge, 0),
	}
}
