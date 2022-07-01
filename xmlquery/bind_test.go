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
	"embed"
	"errors"
	"testing"
)

//go:embed testdata/bind.xml
var bindFs embed.FS

func TestBind(t *testing.T) {
	Bind(&bindFs, "testdata/bind.xml")
	if _, have := nodeMap["bind"]; !have {
		t.Fatal("test failed!")
	}
}

func TestBindErrXmlSelectReadFS(t *testing.T) {
	defer func() {
		if re := recover(); re != ErrReadFS {
			t.Fatal("test failed!")
		}
	}()
	Bind(&bindFs, "testdata/bind2.xml")
}

func TestBindXml(t *testing.T) {
	xmlText := `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="bindXml">

	<select id="selectNow">
		select now() as t
	</select>

</xmlquery>
 `
	BindXml(xmlText)
	if _, have := nodeMap["bindXml"]; !have {
		t.Fatal("test failed!")
	}
}

func TestBindRaw(t *testing.T) {
	xmlText := `
<?xml version="1.0" encoding="utf-8"?>
<xmlquery namespace="bindRaw">

	<select id="selectNow">
		select now() as t
	</select>

</xmlquery>
 `
	BindRaw([]byte(xmlText))
	if _, have := nodeMap["bindRaw"]; !have {
		t.Fatal("test failed!")
	}
}

func TestBingRawErrUnmarshal(t *testing.T) {
	var xmlText = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery_aa namespace="" datasource="">
 
 	<select id="selectNow">
 		select now() as t
 	</select>
 	
 </xmlquery>
 `
	defer func() {
		if re := recover(); !errors.Is(re.(error), ErrUnmarshal) {
			t.Fatal("test failed!")
		}
	}()
	BindRaw([]byte(xmlText))
}

func TestCheckNodeErrXmlSelectIDEmpty(t *testing.T) {
	defer func() {
		if re := recover(); re != ErrIDEmpty {
			t.Fatal("test failed!")
		}
	}()
	checkNode(&rootNode{"ns", "", []*insertNode{{InnerXml: "SQL"}}, []*deleteNode{{InnerXml: "SQL"}}, []*updateNode{{InnerXml: "SQL"}}, []*selectNode{{InnerXml: "SQL"}}})
}

func TestCheckNodeErrXmlSelectSQLEmpty(t *testing.T) {
	defer func() {
		if re := recover(); re != ErrSQLEmpty {
			t.Fatal("test failed!")
		}
	}()
	checkNode(&rootNode{"ns", "", []*insertNode{{ID: "SQL"}}, []*deleteNode{{ID: "SQL"}}, []*updateNode{{ID: "SQL"}}, []*selectNode{{ID: "SQL"}}})
}

func TestGetNodeErrUnbindNamespace(t *testing.T) {
	var xmlText = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestGetNodeErrUnbindNamespace" datasource="">
 
 	<select id="selectNow">
 		select now() as t
 	</select>
 	
 </xmlquery>
 `
	BindXml(xmlText)
	defer func() {
		if re := recover(); !errors.Is(re.(error), ErrUnbindNamespace) {
			t.Fatal("test failed!")
		}
	}()
	getNode("TestGetNodeErrUnbindNamespace_2", "sss", selectType)
}

func TestGetNodeErrNSEmpty(t *testing.T) {
	var xmlText = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="" datasource="">
 
 	<select id="selectNow">
 		select now() as t
 	</select>
 	
 </xmlquery>
 `
	defer func() {
		if re := recover(); !errors.Is(re.(error), ErrNSEmpty) {
			t.Fatal("test failed!")
		}
	}()
	BindXml(xmlText)
}

func TestGetNodeErrUnbindNode(t *testing.T) {
	var xmlText = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestGetNodeErrUnbindNode" datasource="">
 
 	<select id="selectNow">
 		select now() as t
 	</select>
 	
 </xmlquery>
 `
	BindXml(xmlText)
	defer func() {
		if re := recover(); !errors.Is(re.(error), ErrUnbindNode) {
			t.Fatal("test failed!")
		}
	}()
	getNode("TestGetNodeErrUnbindNode", "sss", selectType)
}

func TestGetNodeErrNSDuplicate(t *testing.T) {
	var xmlText = `
 <?xml version="1.0" encoding="utf-8"?>
 <xmlquery namespace="TestGetNodeErrUnbindNode" datasource="">
 
 	<select id="selectNow">
 		select now() as t
 	</select>
 	
 </xmlquery>
 `
	defer func() {
		if re := recover(); !errors.Is(re.(error), ErrNSDuplicate) {
			t.Fatal("test failed!")
		}
	}()
	BindXml(xmlText)
	BindXml(xmlText)
}
