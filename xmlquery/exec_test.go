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

package xmlquery

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-the-way/anorm"
	"os"
	"testing"
)

func init() {
	dsn := os.Getenv("ANORM_TEST_DSN")
	if tdb, err := sql.Open("mysql", dsn+"?parseTime=true"); err != nil {
		panic(err)
	} else {
		anorm.DataSourcePool.Push(tdb)
	}
}

func TestInsert(t *testing.T) {
	var xmlText = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="insert" datasource=""> 
 	<insert id="insert">
		SQL
 	</insert>
 </xmlquery>
 `
	BindXml(xmlText)
	if i := Insert("insert", "insert"); i == nil {
		t.Fatal("test failed!")
	}
}

func TestDelete(t *testing.T) {
	var xmlText = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="delete" datasource=""> 
 	<delete id="delete">
		SQL
 	</delete>
 </xmlquery>
 `
	BindXml(xmlText)
	if i := Delete("delete", "delete"); i == nil {
		t.Fatal("test failed!")
	}
}

func TestUpdate(t *testing.T) {
	var xmlText = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="update" datasource=""> 
 	<update id="update">
		SQL
 	</update>
 </xmlquery>
 `
	BindXml(xmlText)
	if i := Update("update", "update"); i == nil {
		t.Fatal("test failed!")
	}
}

func TestInsertExec(t *testing.T) {
	var xmlText = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="TestInsertExec" datasource="">
 	<update id="createTable">
		create table TestInsertExec(id int primary key auto_increment,name varchar(20) not null default '')
 	</update>
 	<insert id="insertTable">
		insert into TestInsertExec(name) values (?)
 	</insert>
 	<select id="selectTable">
		select count(0) from TestInsertExec where name = 'Coco'
 	</select>
 	<update id="dropTable">
		drop table TestInsertExec
 	</update>
 </xmlquery>
 `
	BindXml(xmlText)

	dropF := func() {
		if ud := Update("TestInsertExec", "dropTable"); ud == nil {
			t.Fatal("test failed!")
		} else {
			if _, err := ud.Exec(); err != nil {
				t.Fatal("test failed!")
			}
		}
	}
	defer dropF()

	if ud := Update("TestInsertExec", "createTable"); ud == nil {
		t.Fatal("test failed!")
	} else {
		if _, err := ud.Exec(); err != nil {
			t.Fatal("test failed!")
		}
	}

	if is := Insert("TestInsertExec", "insertTable"); is == nil {
		t.Fatal("test failed!")
	} else {
		if _, err := is.Exec("Coco"); err != nil {
			t.Fatal("test failed!")
		}
	}

	if is := SingleSelect[int]("TestInsertExec", "selectTable"); is == nil {
		t.Fatal("test failed!")
	} else {
		if one, err := is.One(); err != nil {
			t.Fatal("test failed!")
		} else if one != 1 {
			t.Fatal("test failed!")
		}
	}
}

func TestInsertExecTemplate(t *testing.T) {
	var xmlText = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="TestInsertExecTemplate" datasource="">
 	<update id="createTable">
		create table TestInsertExecTemplate(id int primary key auto_increment,name varchar(20) not null default '')
 	</update>
 	<insert id="insertTable">
		insert into TestInsertExecTemplate(name) values ('{{.}}')
 	</insert>
 	<select id="selectTable">
		select count(0) from TestInsertExecTemplate where name = 'Coco'
 	</select>
 	<update id="dropTable">
		drop table TestInsertExecTemplate
 	</update>
 </xmlquery>
 `
	BindXml(xmlText)

	dropF := func() {
		if ud := Update("TestInsertExecTemplate", "dropTable"); ud == nil {
			t.Fatal("test failed!")
		} else {
			if _, err := ud.Exec(); err != nil {
				t.Fatal("test failed!")
			}
		}
	}
	defer dropF()

	if ud := Update("TestInsertExecTemplate", "createTable"); ud == nil {
		t.Fatal("test failed!")
	} else {
		if _, err := ud.Exec(); err != nil {
			t.Fatal("test failed!")
		}
	}

	if is := Insert("TestInsertExecTemplate", "insertTable"); is == nil {
		t.Fatal("test failed!")
	} else {
		if _, err := is.ExecTemplate("Coco"); err != nil {
			t.Fatal("test failed!")
		}
	}

	if is := SingleSelect[int]("TestInsertExecTemplate", "selectTable"); is == nil {
		t.Fatal("test failed!")
	} else {
		if one, err := is.One(); err != nil {
			t.Fatal("test failed!")
		} else if one != 1 {
			t.Fatal("test failed!")
		}
	}
}

func TestExecExecErr(t *testing.T) {
	var XML = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="TestExecExecErr" datasource="">

	<insert id="createTable">
		create table TestExecExecErrE(id int)
	</insert>

	<insert id="createTable">
		create table TestExecExecErrE(id int)
	</insert>

	<insert id="insert">
		insert xxx
	</insert>

</xmlquery>
`
	BindXml(XML)
	if _, err := Insert("TestExecExecErr", "insert").Exec(); err == nil {
		t.Fatal("test failed!")
	}
}

func TestExecTemplateParseErr(t *testing.T) {
	var XML = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="TestExecTemplateParseErr" datasource="">

	<insert id="insert">
		insert xxx {{{}}}
	</insert>

</xmlquery>
`
	BindXml(XML)
	if _, err := Insert("TestExecTemplateParseErr", "insert").ExecTemplate(nil); err == nil {
		t.Fatal("test failed!")
	}
}

func TestExecTemplateExecuteErr(t *testing.T) {
	var XML = `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="TestExecTemplateExecuteErr" datasource="">

	<insert id="insert">
		insert xxx {{.ABC}}
	</insert>

</xmlquery>
`
	BindXml(XML)
	if _, err := Insert("TestExecTemplateExecuteErr", "insert").ExecTemplate(111); err == nil {
		t.Fatal("test failed!")
	}
}
