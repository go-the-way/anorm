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
	"github.com/go-the-way/anorm"
	"github.com/go-the-way/anorm/pagination"
	"testing"
)

type testPageSelectE struct{ T string }

func (t *testPageSelectE) Configure(*anorm.EC) {}

func TestPageSelect(t *testing.T) {
	type eE = testPageSelectE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelect" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if se := PageSelect[*eE]("TestPageSelect", "selectNow"); se == nil {
		t.Fatal("test failed!")
	}
}

type testPageSelectListE struct{ T string }

func (t *testPageSelectListE) Configure(*anorm.EC) {}

func TestPageSelectList(t *testing.T) {
	type eE = testPageSelectListE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectList" datasource="">
 
 	<select id="selectPage">
 		select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all select '3'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if es, c, err := PageSelect[*eE]("TestPageSelectList", "selectPage").List(new(eE), pagination.MySql, 0, 10); err != nil {
		t.Fatal("test failed!")
	} else if len(es) != 4 {
		t.Fatal("test failed!")
	} else if c != 4 {
		t.Fatal("test failed!")
	}
}

type testPageSelectListQueryCountErrE struct{ T string }

func (t *testPageSelectListQueryCountErrE) Configure(*anorm.EC) {}

func TestPageSelectListQueryCountErr(t *testing.T) {
	type eE = testPageSelectListQueryCountErrE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectListQueryCountErr" datasource="">
 
 	<select id="selectPage">
 		select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all '3' 
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, _, err := PageSelect[*eE]("TestPageSelectListQueryCountErr", "selectPage").List(new(eE), pagination.MySql, 0, 10); err == nil {
		t.Fatal("test failed!")
	}
}

type testPageSelectListQueryCount0E struct{ T string }

func (t *testPageSelectListQueryCount0E) Configure(*anorm.EC) {}

func TestPageSelectListQueryCount0(t *testing.T) {
	type eE = testPageSelectListQueryCount0E
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectListQueryCount0" datasource="">
 
 	<select id="selectPage">
 		select t.* from (select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all select '3')as t where t.T = 'hea' 
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if es, c, err := PageSelect[*eE]("TestPageSelectListQueryCount0", "selectPage").List(new(eE), pagination.MySql, 0, 10); err != nil {
		t.Fatal("test failed!")
	} else if len(es) != 0 {
		t.Fatal("test failed!")
	} else if c != 0 {
		t.Fatal("test failed!")
	}
}

type testPageSelectListErrE struct{ T string }

func (t *testPageSelectListErrE) Configure(*anorm.EC) {}

func TestPageSelectListErr(t *testing.T) {
	type eE = testPageSelectListErrE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectListErr" datasource="">
 
 	<select id="selectPage">
 		select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all select '3' 
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, _, err := PageSelect[*eE]("TestPageSelectListErr", "selectPage").List(new(eE), pagination.SqlServer("hallo"), 0, 10); err == nil {
		t.Fatal("test failed!")
	}
}

type testPageSelectListTemplateE struct{ T string }

func (t *testPageSelectListTemplateE) Configure(*anorm.EC) {}

func TestPageSelectListTemplate(t *testing.T) {
	type eE = testPageSelectListTemplateE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectListTemplate" datasource="">
 
 	<select id="selectPage">
 		select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all select '3' 
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if es, c, err := PageSelect[*eE]("TestPageSelectListTemplate", "selectPage").ListTemplate(new(eE), pagination.MySql, 0, 10, nil); err != nil {
		t.Fatal("test failed!")
	} else if len(es) != 4 {
		t.Fatal("test failed!")
	} else if c != 4 {
		t.Fatal("test failed!")
	}
}

type testPageSelectListTemplateParseErrE struct{ T string }

func (t *testPageSelectListTemplateParseErrE) Configure(*anorm.EC) {}

func TestPageSelectListTemplateParseErr(t *testing.T) {
	type eE = testPageSelectListTemplateParseErrE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectListTemplateParseErr" datasource="">
 
 	<select id="selectPage">
 		select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all select '3' {{{}}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, _, err := PageSelect[*eE]("TestPageSelectListTemplateParseErr", "selectPage").ListTemplate(new(eE), pagination.MySql, 0, 10, nil); err == nil {
		t.Fatal("test failed!")
	}
}

type testPageSelectListTemplateExecuteErrE struct{ T string }

func (t *testPageSelectListTemplateExecuteErrE) Configure(*anorm.EC) {}

func TestPageSelectListTemplateExecuteErr(t *testing.T) {
	type eE = testPageSelectListTemplateExecuteErrE
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestPageSelectListTemplateExecuteErr" datasource="">
 
 	<select id="selectPage">
 		select 'Haha' as T 
		union all select '1'
		union all select '2'
		union all select '3{{.AAA}}'
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, _, err := PageSelect[*eE]("TestPageSelectListTemplateExecuteErr", "selectPage").ListTemplate(new(eE), pagination.MySql, 0, 10, 1); err == nil {
		t.Fatal("test failed!")
	}
}
