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
	"strings"
)

var (
	errRequiredMasterDS        = errors.New("anorm: required a master DS, please set named `master` or `_` DS")
	errCannotRegisterNilEntity = errors.New("anorm: can not register a not nil entity")
	errEntityMetadataNil       = errors.New("anorm: entity's metadata is nil")
	errDuplicateRegisterEntity = errors.New("anorm: duplicate register entity")

	entityTableMap        = make(map[string]string)              // K<entityPKGName> V<TableName>
	entityTagMap          = make(map[string]map[string]*tag)     // K<entityPKGName> V< K<ColumnName> V<*Tag> >
	entityColumnMap       = make(map[string][]string)            // K<entityPKGName> V<[]Column>
	entityFieldMap        = make(map[string][]string)            // K<entityPKGName> V<[]Column>
	entityFieldColumnMap  = make(map[string]map[string]string)   // K<entityPKGName> V< K<Field> V<Column> >
	entityColumnFieldMap  = make(map[string]map[string]string)   // K<entityPKGName> V< K<Column> V<Field> >
	entityInsertIgnoreMap = make(map[string]map[string]struct{}) // K<entityPKGName> V< K<Column> V<> >
	entityUpdateIgnoreMap = make(map[string]map[string]struct{}) // K<entityPKGName> V< K<Column> V<> >
	entityJoinRefMap      = make(map[string]map[string]*JoinRef) // K<entityPKGName> V< K<Field> V<> >
	entityPKMap           = make(map[string][]string, 0)         // K<entityPKGName> V<[]PKColumn>
	entityDSMap           = make(map[string]string, 0)           // K<entityPKGName> V<DS>
)

type (
	// Entity alias for EntityConfigurator
	Entity = EntityConfigurator
	// EC alias for EntityConfiguration
	EC = EntityConfiguration
	// EntityConfigurator Entity Configurator
	EntityConfigurator interface {
		Configure(c *EC)
	}
	// EntityConfiguration tells orm how to configure this Entity
	EntityConfiguration struct {
		// Migrate defines migrate create or update table in the database
		Migrate bool
		// Table defines table name for the EntityConfigurator
		// If Table is empty, set from Configuration's TableNameStrategy
		Table string
		// IFNotExists only mysql supports
		IFNotExists bool
		// Commented comment options?
		Commented bool
		// Comment defines table comment for the EntityConfigurator
		Comment string
		// Table defines column name Strategy for the EntityConfigurator's Fields
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
		// JoinRefMap defines Table join rel table
		JoinRefMap map[string]*JoinRef
	}

	tag struct {
		PK           bool   `alias:"pk"`
		Column       string `alias:"c"`
		InsertIgnore bool   `alias:"ig"`
		UpdateIgnore bool   `alias:"ug"`
		Definition   string `alias:"def"`
		Join         string `alias:"join"` // join{type,self_column,join_table,join_column}
	}

	JoinRef struct {
		Field      string // SelfID
		Type       string // left,right,inner,...
		SelfColumn string // self_id
		RelTable   string // rel_table
		RelID      string // rel_id
		RelName    string // rel_name
	}
)

func (t *tag) String() string {
	return fmt.Sprintf("PK:%v, Column:%s, InsertIgnore:%v, UpdateIgnore:%v, Definition:%s, Join:%s", t.PK, t.Column, t.InsertIgnore, t.UpdateIgnore, t.Definition, t.Join)
}

// Register defines register a EntityConfigurator struct for anorm
func Register(entity Entity) {
	c := &EC{}
	entity.Configure(c)
	entityPkgName := getEntityPkgName(entity)
	if _, have := entityTableMap[entityPkgName]; have {
		panic(errDuplicateRegisterEntity)
	}

	insertIgnoreMap := make(map[string]struct{}, 0)
	updateIgnoreMap := make(map[string]struct{}, 0)
	joinRefMap := make(map[string]*JoinRef, 0)
	tagMap := make(map[string]*tag, 0)
	fields := make([]string, 0)
	columns := make([]string, 0)
	columnFieldMap := make(map[string]string, 0)
	fieldColumnMap := make(map[string]string, 0)
	pks := make([]string, 0)
	pkGes := make([]sg.Ge, 0)
	columnGes := make([]sg.Ge, 0)

	rt := reflect.TypeOf(entity).Elem()
	numField := rt.NumField()

	structFields, _, tags := reflectx.ParseTagWithRe(entity, new(tag), "alias", "orm", false, "([a-zA-Z0-9]+){([^{}]+)}")
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
			join             string
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
			join = curTag.Join
			Logger.Debug([]*logField{LogField("entity", reflect.TypeOf(entity)), LogField("tag", curTag.String())}, "parsed")
		}
		if pk {
			pks = append(pks, column)
			pkGes = append(pkGes, sg.C(column))
		}
		if column == "" {
			column = getStrategyName(fieldName, c.ColumnNameStrategy)
		}
		if columnDefinition != "" {
			columnGes = append(columnGes, sg.C(columnDefinition))
		}
		if insertIgnore {
			insertIgnoreMap[column] = struct{}{}
		}
		if updateIgnore {
			updateIgnoreMap[column] = struct{}{}
		}
		setJoinMap(entity, fieldName, curTag, join, joinRefMap)

		fields = append(fields, fieldName)
		columns = append(columns, column)
		columnFieldMap[column] = fieldName
		fieldColumnMap[fieldName] = column
	}

	if vs := c.PrimaryKeyColumns; vs != nil {
		for _, v := range vs {
			pks = append(pks, string(v))
			pkGes = append(pkGes, v)
		}
	}

	if vs := c.ColumnDefinitions; vs != nil {
		for _, v := range vs {
			columnGes = append(columnGes, v)
		}
	}

	if vs := c.InsertIgnores; vs != nil {
		for _, v := range vs {
			insertIgnoreMap[string(v)] = struct{}{}
		}
	}

	if vs := c.UpdateIgnores; vs != nil {
		for _, v := range vs {
			updateIgnoreMap[string(v)] = struct{}{}
		}
	}

	if vs := c.JoinRefMap; vs != nil {
		for k, v := range vs {
			joinRefMap[k] = v
		}
	}

	pkgName := reflect.TypeOf(entity).String()
	pss := strings.Split(pkgName, ".")
	tableName := pss[len(pss)-1]
	if t := c.Table; t != "" {
		tableName = t
	}

	entityTableMap[entityPkgName] = tableName
	entityTagMap[entityPkgName] = tagMap
	entityFieldMap[entityPkgName] = fields
	entityColumnMap[entityPkgName] = columns
	entityFieldColumnMap[entityPkgName] = fieldColumnMap
	entityColumnFieldMap[entityPkgName] = columnFieldMap
	entityInsertIgnoreMap[entityPkgName] = insertIgnoreMap
	entityUpdateIgnoreMap[entityPkgName] = updateIgnoreMap
	entityJoinRefMap[entityPkgName] = joinRefMap
	entityPKMap[entityPkgName] = pks
	if c.DS == "" {
		c.DS = "_"
	}
	entityDSMap[entityPkgName] = c.DS

	if !c.Migrate {
		return
	}

	c.Table = tableName
	migrate(entity, c, pkGes, columnGes)
}

func setJoinMap(entity Entity, fieldName string, curTag *tag, join string, joinRefMap map[string]*JoinRef) {
	if join == "" {
		return
	}
	if joinPs := strings.Split(join, ","); len(joinPs) != 5 {
		Logger.Error([]*logField{LogField("entity", reflect.TypeOf(entity)), LogField("tag", curTag.String())}, "parse err: valid like [inner,self_id,rel_table,rel_id,rel_id,rel_name]")
	} else {
		joinType := strings.ToUpper(strings.TrimSpace(joinPs[0]))
		selfColumn := strings.TrimSpace(joinPs[1])
		relTable := strings.TrimSpace(joinPs[2])
		relID := strings.TrimSpace(joinPs[3])
		relName := strings.TrimSpace(joinPs[4])
		jr := JoinRef{
			Field:      fieldName,
			Type:       joinType,
			SelfColumn: selfColumn,
			RelTable:   relTable,
			RelID:      relID,
			RelName:    relName,
		}
		joinRefMap[fieldName] = &jr
	}
}

func migrate(entity Entity, c *EC, pkGes []sg.Ge, columnGes []sg.Ge) {
	db := DataSourcePool.Required(c.DS)

	builder := sg.CreateTableBuilder().
		Table(sg.T(c.Table)).
		Comment(c.Comment).
		PrimaryKey(pkGes...).
		ColumnDefinition(columnGes...)

	if c.IFNotExists {
		builder.IfNotExist()
	}

	if is := c.IndexDefinitions; is != nil && len(is) > 0 {
		builder.Index(is...)
	}

	createSQL, _ := builder.Build()
	Logger.Debug([]*logField{LogField("entity", reflect.TypeOf(entity)), LogField("table", c.Table), LogField("DDL", createSQL)}, "migrating entity ...")

	if _, err := db.Exec(createSQL); err != nil {
		Logger.Fatal([]*logField{LogField("entity", reflect.TypeOf(entity)), LogField("table", c.Table), LogField("DDL", createSQL)}, "migrate entity err: %v", err)
	}

}
