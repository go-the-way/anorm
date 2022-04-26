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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-the-way/sg"
	"reflect"
	"sync"
	"time"
)

var (
	errEntityNil     = errors.New(fmt.Sprintf("anorm: entity is nil"))
	errUnknownEntity = func(entity EntityConfigurator) error {
		return errors.New(fmt.Sprintf("anorm: unknown entity [%v]", getEntityPkgName(entity)))
	}
	errTxNotOpen            = errors.New("anorm: tx not open")
	errTxManagerNil         = errors.New("anorm: tx manager is nil")
	errAlreadyBindTxManager = errors.New("anorm: already bind tx manager")
)

type Orm[E Entity] struct {
	mu *sync.Mutex

	entity E

	tagMap map[string]*tag

	db     *sql.DB
	tx     *sql.Tx
	openTx bool

	txm *txManager
}

// New defines return a new Orm from EntityConfigurator entity
func New[E Entity](entity E) *Orm[E] {
	return NewWithDS[E](entity, "")
}

// NewWithDS defines return a new Orm using named DS from EntityConfigurator entity
func NewWithDS[E Entity](entity E, ds string) *Orm[E] {
	if !entityNotNil(entity) {
		panic(errEntityNil)
	}

	if ds == "" {
		ds = entityDSMap[getEntityPkgName(entity)]
	}

	o := &Orm[E]{
		mu:     &sync.Mutex{},
		entity: entity,
		db:     DataSourcePool.Required(ds),
	}

	Logger.Debug([]*logField{LogField("entity", getEntityPkgName(entity)), LogField("DS", ds)}, "created")

	if tagMap, registered := entityTagMap[getEntityPkgName(entity)]; !registered {
		panic(errUnknownEntity(entity))
	} else {
		o.tagMap = tagMap
	}

	return o
}

// table defines return EntityConfigurator's table name
func (o *Orm[E]) table() sg.Ge {
	return sg.T(entityTableMap[getEntityPkgName(o.entity)])
}

func (o *Orm[E]) getWhereGes(entity E) []sg.Ge {
	ges := make([]sg.Ge, 0)
	if entityNotNil(entity) {
		fieldColumnMap := entityFieldColumnMap[getEntityPkgName(entity)]
		rt := reflect.TypeOf(entity).Elem()
		rv := reflect.ValueOf(entity).Elem()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			value := rv.Field(i)
			if val := o.getRealVal(value); val != nil {
				ges = append(ges, sg.Eq(sg.C(fieldColumnMap[field.Name]), val))
			}
		}
	}
	return ges
}

func (o *Orm[E]) getRealVal(value reflect.Value) (val any) {
	if value.Kind() == reflect.Ptr {
		return o.getRealVal(value.Elem())
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v := value.Int(); v > 0 {
			return v
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v := value.Uint(); v > 0 {
			return v
		}
	case reflect.Float32, reflect.Float64:
		if v := value.Float(); v > 0 {
			return v
		}
	case reflect.String:
		if v := value.String(); v != "" {
			return v
		}
	case reflect.Struct:
		if v := value.Interface(); v != nil {
			switch value.Type() {
			case reflect.TypeOf(NullBool(false)):
				if vv, ok := v.(sql.NullBool); ok {
					return vv
				}
			case reflect.TypeOf(NullByte(0)):
				if vv, ok := v.(sql.NullByte); ok {
					return vv
				}
			case reflect.TypeOf(NullInt16(0)):
				if vv, ok := v.(sql.NullInt16); ok {
					return vv
				}
			case reflect.TypeOf(NullInt32(0)):
				if vv, ok := v.(sql.NullInt32); ok {
					return vv
				}
			case reflect.TypeOf(NullInt64(0)):
				if vv, ok := v.(sql.NullInt64); ok {
					return vv
				}
			case reflect.TypeOf(NullFloat64(0)):
				if vv, ok := v.(sql.NullFloat64); ok {
					return vv
				}
			case reflect.TypeOf(NullString("")):
				if vv, ok := v.(sql.NullString); ok {
					return vv
				}
			case reflect.TypeOf(time.Time{}):
				if vv, ok := v.(time.Time); ok && vv.After(time.Time{}) {
					return vv
				}
			case reflect.TypeOf(NullTime(time.Time{})):
				if vv, ok := v.(sql.NullTime); ok && vv.Valid && vv.Time.After(time.Time{}) {
					return vv
				}

				/* ignore Pointer
				case reflect.TypeOf(NullBoolPtr(false)):
					if vv, ok := v.(*sql.NullBool); ok {
						return vv
					}
				case reflect.TypeOf(NullBytePtr(0)):
					if vv, ok := v.(*sql.NullByte); ok {
						return vv
					}
				case reflect.TypeOf(NullInt16Ptr(0)):
					if vv, ok := v.(*sql.NullInt16); ok {
						return vv
					}
				case reflect.TypeOf(NullInt32Ptr(0)):
					if vv, ok := v.(*sql.NullInt32); ok {
						return vv
					}
				case reflect.TypeOf(NullInt64Ptr(0)):
					if vv, ok := v.(*sql.NullInt64); ok {
						return vv
					}
				case reflect.TypeOf(NullFloat64Ptr(0)):
					if vv, ok := v.(*sql.NullFloat64); ok {
						return vv
					}
				case reflect.TypeOf(NullStringPtr("")):
					if vv, ok := v.(*sql.NullString); ok {
						return vv
					}
				case reflect.TypeOf(&time.Time{}):
					if vv, ok := v.(*time.Time); ok {
						return vv
					}
				case reflect.TypeOf(NullTimePtr(time.Time{})):
					if vv, ok := v.(*sql.NullTime); ok {
						return vv
					}
				*/

			}
		}
	}
	return
}

// OpsForSelect defines return *selectOperation
func (o *Orm[E]) OpsForSelect() *selectOperation[E] {
	return newSelectOperation(o)
}

// OpsForSelectCount defines return *selectCountOperation
func (o *Orm[E]) OpsForSelectCount() *selectCountOperation[E] {
	return newSelectCountOperation(o)
}

// OpsForInsert defines return *insertOperation
func (o *Orm[E]) OpsForInsert() *insertOperation[E] {
	return newInsertOperation(o)
}

// OpsForUpdate defines return *updateOperation
func (o *Orm[E]) OpsForUpdate() *updateOperation[E] {
	return newUpdateOperation(o)
}

// OpsForDelete defines return *deleteOperation
func (o *Orm[E]) OpsForDelete() *deleteOperation[E] {
	return newsDeleteOperation(o)
}

func beginTx(db *sql.DB, options ...*sql.TxOptions) (tx *sql.Tx, err error) {
	if options != nil && len(options) > 0 {
		if tx, err = db.BeginTx(context.Background(), options[0]); err != nil {
			return nil, err
		}
	} else {
		if tx, err = db.Begin(); err != nil {
			return nil, err
		}
	}
	return
}

// BeginTx begin a tx with tx manager
func (o *Orm[E]) BeginTx(txm *txManager, options ...*sql.TxOptions) error {
	if txm == nil {
		return errTxManagerNil
	}
	if o.txm != nil {
		return errAlreadyBindTxManager
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if tx, err := beginTx(o.db, options...); err != nil {
		return err
	} else {
		o.openTx = true
		o.tx = tx
	}
	txm.Join(o.tx)
	o.txm = txm
	return nil
}
