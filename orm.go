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
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/go-the-way/sg"
)

var (
	errModelIsNil              = func(model Model) error { return errors.New(fmt.Sprintf("model is nil")) }
	errModelNotYetRegisterFunc = func(model Model) error {
		return errors.New(fmt.Sprintf("model [%v] not yet register", getModelPkgName(model)))
	}
	errTxNotOpen = errors.New("anorm: Tx not open")
)

type orm struct {
	mux    *sync.Mutex
	model  Model
	tagMap map[string]*tag
	db     *sql.DB
	tx     *sql.Tx
	openTx bool
}

// New defines return a new orm from Model model
func New(model Model) *orm {
	return NewWithDS(model, ifEmpty(model.MetaData().DS, "_"))
}

// NewWithDS defines return a new orm using named DS from Model model
func NewWithDS(model Model, ds string) *orm {
	if model == nil {
		panic(errModelIsNil)
	}

	o := &orm{
		mux:   &sync.Mutex{},
		model: model,
		db:    dsMap.required(ds),
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

// Select defines return *_Select
func (o *orm) Select() *_Select {
	return newSelect(o)
}

// SelectCount defines return *_SelectCount
func (o *orm) SelectCount() *_SelectCount {
	return newSelectCount(o)
}

// Insert defines return *_Insert
func (o *orm) Insert() *_Insert {
	return newInsert(o)
}

// Update defines return *_Update
func (o *orm) Update() *_Update {
	return newUpdate(o)
}

// Delete defines return *_Delete
func (o *orm) Delete() *_Delete {
	return newDelete(o)
}

// Begin a Tx
func (o *orm) Begin() error {
	o.mux.Lock()
	defer o.mux.Unlock()
	if tx, err := o.db.Begin(); err != nil {
		return err
	} else {
		o.openTx = true
		o.tx = tx
		return nil
	}
}

// Commit the Tx
func (o *orm) Commit() error {
	o.mux.Lock()
	defer o.mux.Unlock()
	if !o.openTx {
		return errTxNotOpen
	}
	return o.tx.Commit()
}

// Rollback the Tx
func (o *orm) Rollback() error {
	o.mux.Lock()
	defer o.mux.Unlock()
	if !o.openTx {
		return errTxNotOpen
	}
	return o.tx.Rollback()
}
