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
	"fmt"
	"github.com/billcoding/reflectx"
	"github.com/go-the-way/sg"
	"reflect"
)

var (
	errRequiredDS             = errors.New("anorm: required a DS, please call anorm.DS or anorm.DSWithName to set")
	errRequiredMasterDS       = errors.New("anorm: required a master DS, please set named `master` or `_` DS")
	errCannotRegisterNilModel = errors.New("anorm: can not register a not nil model")
	errModelMetadataNil       = errors.New("anorm: model's metadata is nil")
	errDuplicateRegisterModel = errors.New("anorm: duplicate register model")

	// K<ModelPKGName> V<TableName>
	modelTableMap = make(map[string]string)
	// K<ModelPKGName> V< K<ColumnName> V<*Tag> >
	modelTagMap = make(map[string]map[string]*tag)
	// K<ModelPKGName> V<[]Column>
	modelColumnMap = make(map[string][]string)
	// K<ModelPKGName> V<[]Column>
	modelFieldMap = make(map[string][]string)
	// K<ModelPKGName> V< K<Field> V<Column> >
	modelFieldColumnMap = make(map[string]map[string]string)
	// K<ModelPKGName> V< K<Column> V<Field> >
	modelColumnFieldMap = make(map[string]map[string]string)
	// K<ModelPKGName> V< K<Column> V<> >
	modelInsertIgnoreMap = make(map[string]map[string]struct{})
	// K<ModelPKGName> V< K<Column> V<> >
	modelUpdateIgnoreMap = make(map[string]map[string]struct{})
	// K<ModelPKGName> V<[]PKColumn>
	modelPKMap = make(map[string][]string, 0)
)

type (
	// Model defines standard Model interface and implements to marked a Model struct
	Model interface {
		MetaData() *ModelMeta
	}
	// ModelMeta defines Model struct meta info
	ModelMeta struct {
		// Migrate defines migrate create or update table in the database
		Migrate bool
		// Table defines table name for the Model
		// If Table is empty, set from Configuration's TableNameStrategy
		Table string
		// Comment defines table comment for the Model
		Comment string
		// Table defines column name Strategy for the Model's Fields
		// see strategy.Default, strategy.Underline, strategy.CamelCase,
		ColumnNameStrategy Strategy
		// DB defines DS name from dsMap, default: `_` or `master`
		DS string
		// PrimaryKeyColumns defines Table primary key columns
		PrimaryKeyColumns []sg.C
		// ColumnDefinitions defines Table column definitions
		// Use sg.ColumnDefinition to build
		ColumnDefinitions []sg.Ge
		// Indexes defines Table index definitions
		// Use sg.IndexDefinition to build
		IndexDefinitions []sg.Ge
		// InsertIgnores defines Table insert ignore columns
		InsertIgnores []sg.C
		// UpdateIgnores defines Table update ignore columns
		UpdateIgnores []sg.C
	}

	tag struct {
		PK           bool   `alias:"pk"`
		Column       string `alias:"c"`
		InsertIgnore bool   `alias:"ig"`
		UpdateIgnore bool   `alias:"ug"`
		Definition   string `alias:"def"`
	}
)

func (t *tag) String() string {
	return fmt.Sprintf("{PK:%v, Column:%s, InsertIgnore:%v, UpdateIgnore:%v, Definition:%s}", t.PK, t.Column, t.InsertIgnore, t.UpdateIgnore, t.Definition)
}

// Register defines register a Model struct for anorm
func Register(model Model) {
	if model == nil {
		panic(errCannotRegisterNilModel)
	}

	metaData := model.MetaData()
	if metaData == nil {
		panic(errModelMetadataNil)
	}

	modelPkgName := getModelPkgName(model)
	if _, have := modelTableMap[modelPkgName]; have {
		panic(errDuplicateRegisterModel)
	}

	insertIgnoreMap := make(map[string]struct{}, 0)
	updateIgnoreMap := make(map[string]struct{}, 0)
	tagMap := make(map[string]*tag, 0)
	fields := make([]string, 0)
	columns := make([]string, 0)
	columnFieldMap := make(map[string]string, 0)
	fieldColumnMap := make(map[string]string, 0)
	pks := make([]string, 0)
	pkGes := make([]sg.Ge, 0)
	columnGes := make([]sg.Ge, 0)

	rt := reflect.TypeOf(model).Elem()
	numField := rt.NumField()

	structFields, _, tags := reflectx.ParseTagWithRe(model, new(tag), "alias", "orm", false, "([a-zA-Z0-9]+){([^{}]+)}")
	indexMap := make(map[string]int, 0)
	for i := range structFields {
		indexMap[structFields[i].Name] = i
	}

	for i := 0; i < numField; i++ {
		structField := rt.Field(i)
		var (
			curTag           *tag
			column           string
			columnDefinition string
			insertIgnore     bool
			updateIgnore     bool
			pk               bool
		)
		fieldName := structField.Name
		if len(indexMap) > 0 {
			if curIndex, have := indexMap[fieldName]; have {
				curTag = tags[curIndex].(*tag)
			}
		}
		if curTag != nil {
			column = curTag.Column
			columnDefinition = curTag.Definition
			insertIgnore = curTag.InsertIgnore
			updateIgnore = curTag.UpdateIgnore
			pk = curTag.PK
			debug("parse model [%v] tag [%s]", reflect.TypeOf(model), curTag.String())
		}
		if pk {
			pks = append(pks, column)
			pkGes = append(pkGes, sg.C(column))
		}
		column = ifEmpty(column, getStrategyName(fieldName, Configuration.ColumnNameStrategy|metaData.ColumnNameStrategy))
		if columnDefinition != "" {
			columnGes = append(columnGes, sg.C(columnDefinition))
		}
		if insertIgnore {
			insertIgnoreMap[column] = struct{}{}
		}
		if updateIgnore {
			updateIgnoreMap[column] = struct{}{}
		}
		fields = append(fields, fieldName)
		columns = append(columns, column)
		columnFieldMap[column] = fieldName
		fieldColumnMap[fieldName] = column
	}

	if cs := metaData.PrimaryKeyColumns; cs != nil {
		for _, c := range cs {
			pks = append(pks, string(c))
			pkGes = append(pkGes, c)
		}
	}

	if cs := metaData.ColumnDefinitions; cs != nil {
		for _, c := range cs {
			columnGes = append(columnGes, c)
		}
	}

	if cs := metaData.InsertIgnores; cs != nil {
		for _, c := range cs {
			insertIgnoreMap[string(c)] = struct{}{}
		}
	}

	if cs := metaData.UpdateIgnores; cs != nil {
		for _, c := range cs {
			updateIgnoreMap[string(c)] = struct{}{}
		}
	}

	tableName := getModelTableName(model, Configuration.TableNameStrategy)
	modelTableMap[modelPkgName] = tableName
	modelTagMap[modelPkgName] = tagMap
	modelFieldMap[modelPkgName] = fields
	modelColumnMap[modelPkgName] = columns
	modelFieldColumnMap[modelPkgName] = fieldColumnMap
	modelColumnFieldMap[modelPkgName] = columnFieldMap
	modelInsertIgnoreMap[modelPkgName] = insertIgnoreMap
	modelUpdateIgnoreMap[modelPkgName] = updateIgnoreMap
	modelPKMap[modelPkgName] = pks

	if !(Configuration.Migrate || metaData.Migrate) {
		return
	}

	if len(dsMap) <= 0 {
		panic(errRequiredDS)
	}

	masterDB := dsMap["master"]
	if masterDB == nil {
		masterDB = dsMap["_"]
		if masterDB == nil {
			panic(errRequiredMasterDS)
		}
	}

	builder := sg.CreateTableBuilder().
		Table(sg.T(tableName)).
		IfNotExist().
		Comment(metaData.Comment).
		PrimaryKey(pkGes...).
		ColumnDefinition(columnGes...)

	if is := metaData.IndexDefinitions; is != nil && len(is) > 0 {
		builder.Index(is...)
	}

	createSQL, _ := builder.Build()
	debug("migrate model [%v] table [%v] create DDL [%v]", reflect.TypeOf(model), tableName, createSQL)

	_, err := masterDB.Exec(createSQL)
	if err != nil {
		handleErr(errors.New(fmt.Sprintf("migrate model [%v] table [%v] error [%v]", reflect.TypeOf(model), tableName, err)))
	}
}
