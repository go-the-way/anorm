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
	"encoding/xml"
	"errors"
	"github.com/go-the-way/anorm"
)

var (
	nodeMap     = make(map[string]*rootNode, 0)
	nodeTypeMap = map[nodeType]string{insertType: "insert", deleteType: "delete", updateType: "update", selectType: "select"}

	ErrReadFS      = errors.New("xmlquery: read fs err")
	ErrUnmarshal   = errors.New("xmlquery: unmarshal err")
	ErrParse       = errors.New("xmlquery: parse err")
	ErrNSEmpty     = errors.New("xmlquery: namespace empty")
	ErrNSDuplicate = errors.New("xmlquery: namespace duplicate")
	ErrIDEmpty     = errors.New("xmlquery: id empty")
	ErrSQLEmpty    = errors.New("xmlquery: sql empty")

	ErrUnbindNamespace = errors.New("xmlquery: unbind namespace")
	ErrUnbindNode      = errors.New("xmlquery: unbind node")
)

func getNode(namespace, id string, nt nodeType) (rn *rootNode, nd node) {
	if n, have := nodeMap[namespace]; !have {
		panic(ErrUnbindNamespace)
	} else {
		rn = n
		switch nt {
		case insertType:
			if n.InsertNodes != nil {
				for _, nn := range n.InsertNodes {
					if nn.GetID() == id {
						nd = nn
					}
				}
			}
		case deleteType:
			if n.DeleteNodes != nil {
				for _, nn := range n.DeleteNodes {
					if nn.GetID() == id {
						nd = nn
					}
				}
			}
		case updateType:
			if n.UpdateNodes != nil {
				for _, nn := range n.UpdateNodes {
					if nn.GetID() == id {
						nd = nn
					}
				}
			}
		case selectType:
			if n.SelectNodes != nil {
				for _, nn := range n.SelectNodes {
					if nn.GetID() == id {
						nd = nn
					}
				}
			}
		}
	}

	if nd == nil {
		panic(ErrUnbindNode)
	}

	return
}

// Bind bind named of fs to rootNode
func Bind(fs *embed.FS, name string) {
	if bytes, err := fs.ReadFile(name); err != nil {
		panic(ErrReadFS)
	} else {
		BindRaw(bytes)
	}
}

// BindXml bind xml text to rootNode
func BindXml(xmlText string) {
	BindRaw([]byte(xmlText))
}

func checkNode(rn *rootNode) {
	nodes := make([]node, 0)
	if rn.InsertNodes != nil {
		for _, nn := range rn.InsertNodes {
			nodes = append(nodes, nn)
		}
	}
	if rn.DeleteNodes != nil {
		for _, nn := range rn.DeleteNodes {
			nodes = append(nodes, nn)
		}
	}
	if rn.UpdateNodes != nil {
		for _, nn := range rn.UpdateNodes {
			nodes = append(nodes, nn)
		}
	}
	if rn.SelectNodes != nil {
		for _, nn := range rn.SelectNodes {
			nodes = append(nodes, nn)
		}
	}
	for _, nn := range nodes {
		if id := nn.GetID(); id == "" {
			panic(ErrIDEmpty)
		} else if nn.GetInnerXml() == "" {
			panic(ErrSQLEmpty)
		} else {
			anorm.Logger.Debug(nil, "xmlquery: namespace[%s] datasource[%s] node[type:%s, id:%s, datasource:%s] check pass", rn.Namespace, rn.Datasource, nodeTypeMap[nn.getType()], id, nn.getDatasource())
		}
	}
}

// BindRaw bind raw bytes text to rootNode
func BindRaw(bytes []byte) {
	rn := rootNode{}
	if err := xml.Unmarshal(bytes, &rn); err != nil {
		panic(ErrUnmarshal)
	} else if rn.Namespace == "" {
		panic(ErrNSEmpty)
	} else if _, have := nodeMap[rn.Namespace]; have {
		panic(ErrNSDuplicate)
	} else {
		nodeMap[rn.Namespace] = &rn
		checkNode(&rn)
		anorm.Logger.Debug(nil, "xmlquery: namespace[%s] datasource[%s] parse success", rn.Namespace, rn.Datasource)
	}
}
