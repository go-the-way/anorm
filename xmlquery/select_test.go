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
	"errors"
	"github.com/go-the-way/anorm"
	"testing"
)

type testSelectE struct{ T string }

func (t *testSelectE) Configure(*anorm.EC) {}

func TestSelect(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelect" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if se := Select[*testSelectE]("TestSelect", "selectNow"); se == nil {
		t.Fatal("test failed!")
	}
}

func TestSelectList(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectList" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T union all select 'Hugo'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if es, err := Select[*testSelectE]("TestSelectList", "selectNow").List(new(testSelectE)); err != nil {
		t.Fatal("test failed!")
	} else if len(es) != 2 {
		t.Fatal("test failed!")
	}
}

func TestSelectListErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectListErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T union all 'Hugo'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectListErr", "selectNow").List(new(testSelectE)); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSelectOne(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectOne" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if one, err := Select[*testSelectE]("TestSelectOne", "selectNow").One(new(testSelectE)); one == nil || err != nil {
		t.Fatal("test failed!")
	}
}

func TestSelectOneErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectOneErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T union all
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectOneErr", "selectNow").One(new(testSelectE)); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSelectOneErrSelectTooManyResult(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectOneErrSelectTooManyResult" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T union all select 'Hugo'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectOneErrSelectTooManyResult", "selectNow").One(new(testSelectE)); !errors.Is(err, ErrSelectTooManyResult) {
		t.Fatal("test failed!")
	}
}

func TestSelectListTemplate(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectListTemplate" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T union all select 'Hugo'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectListTemplate", "selectNow").ListTemplate(new(testSelectE), nil); err != nil {
		t.Fatal("test failed!")
	}
}

func TestSelectListTemplateParseErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectListTemplateParseErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T {{{}}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectListTemplateParseErr", "selectNow").ListTemplate(new(testSelectE), nil); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSelectListTemplateExecuteErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectListTemplateExecuteErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T {{.AAA}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectListTemplateExecuteErr", "selectNow").ListTemplate(new(testSelectE), 1); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSelectOneTemplate(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectOneTemplate" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectOneTemplate", "selectNow").OneTemplate(new(testSelectE), 1); err != nil {
		t.Fatal("test failed!")
	}
}

func TestSelectOneTemplateErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectOneTemplateErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T {{
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectOneTemplateErr", "selectNow").OneTemplate(new(testSelectE), 1); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSelectOneTemplateErrSelectTooManyResult(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSelectOneTemplateErrSelectTooManyResult" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T union all select 'Hugo'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := Select[*testSelectE]("TestSelectOneTemplateErrSelectTooManyResult", "selectNow").OneTemplate(new(testSelectE), 1); !errors.Is(err, ErrSelectTooManyResult) {
		t.Fatal("test failed!")
	}
}
