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
	"testing"
)

func TestRowsSelect(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelect" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if se := RowsSelect("TestRowsSelect", "selectNow"); se == nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRow(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRow" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if r := RowsSelect("TestRowsSelectRow", "selectNow").Row(); r == nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRows(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRows" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if r, err := RowsSelect("TestRowsSelectRows", "selectNow").Rows(); r == nil {
		t.Fatal("test failed!")
	} else if err != nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRowTemplate(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRowTemplate" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if r, err := RowsSelect("TestRowsSelectRowTemplate", "selectNow").RowTemplate(nil); r == nil {
		t.Fatal("test failed!")
	} else if err != nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRowTemplateParseErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRowTemplateParseErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T{{{}}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := RowsSelect("TestRowsSelectRowTemplateParseErr", "selectNow").RowTemplate(nil); err == nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRowTemplateExecuteErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRowTemplateExecuteErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T{{.AAA}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := RowsSelect("TestRowsSelectRowTemplateExecuteErr", "selectNow").RowTemplate(1); err == nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRowsTemplate(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRowsTemplate" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if r, err := RowsSelect("TestRowsSelectRowsTemplate", "selectNow").RowsTemplate(nil); r == nil {
		t.Fatal("test failed!")
	} else if err != nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRowsTemplateParseErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRowsTemplateParseErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T{{{}}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := RowsSelect("TestRowsSelectRowsTemplateParseErr", "selectNow").RowsTemplate(nil); err == nil {
		t.Fatal("test failed!")
	}
}

func TestRowsSelectRowsTemplateExecuteErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestRowsSelectRowsTemplateExecuteErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T{{.AAA}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := RowsSelect("TestRowsSelectRowsTemplateExecuteErr", "selectNow").RowsTemplate(1); err == nil {
		t.Fatal("test failed!")
	}
}
