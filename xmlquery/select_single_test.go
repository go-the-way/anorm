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

import "testing"

func TestSingleSelect(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelect" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if se := SingleSelect[string]("TestSingleSelect", "selectNow"); se == nil {
		t.Fatal("test failed!")
	}
}

func TestSingleSelectOne(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelectOne" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if tt, err := SingleSelect[string]("TestSingleSelectOne", "selectNow").One(); err != nil {
		t.Fatal("test failed!")
	} else if tt != "Haha" {
		t.Fatal("test failed!")
	}
}

func TestSingleSelectOneErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelectOneErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T select
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := SingleSelect[string]("TestSingleSelectOneErr", "selectNow").One(); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSingleSelectOneScanErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelectOneScanErr" datasource="">
 
 	<select id="selectNow">
 		select NULL as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := SingleSelect[string]("TestSingleSelectOneScanErr", "selectNow").One(); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSingleSelectOneTemplate(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelectOneTemplate" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if tt, err := SingleSelect[string]("TestSingleSelectOneTemplate", "selectNow").OneTemplate(nil); err != nil {
		t.Fatal("test failed!")
	} else if tt != "Haha" {
		t.Fatal("test failed!")
	}
}

func TestSingleSelectOneTemplateParseErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelectOneTemplateParseErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T{{{}}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := SingleSelect[string]("TestSingleSelectOneTemplateParseErr", "selectNow").OneTemplate(nil); err == nil {
		t.Fatal("test failed!")
	}
}

func TestSingleSelectOneTemplateExecuteErr(t *testing.T) {
	var XML = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestSingleSelectOneTemplateExecuteErr" datasource="">
 
 	<select id="selectNow">
 		select 'Haha' as T{{.AAA}}
 	</select>
 	
 </xmlquery>
 `
	BindXml(XML)
	if _, err := SingleSelect[string]("TestSingleSelectOneTemplateExecuteErr", "selectNow").OneTemplate(1); err == nil {
		t.Fatal("test failed!")
	}
}
