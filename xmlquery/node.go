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

import "github.com/go-the-way/anorm"

type nodeType uint8

const (
	insertType nodeType = iota
	deleteType
	updateType
	selectType
)

type (
	node interface {
		getType() nodeType
		getDatasource() string
		GetID() string
		GetInnerXml() string
	}
	rootNode struct {
		Namespace   string        `xml:"namespace,attr"`
		Datasource  string        `xml:"datasource,attr"`
		InsertNodes []*insertNode `xml:"insert"`
		DeleteNodes []*deleteNode `xml:"delete"`
		UpdateNodes []*updateNode `xml:"update"`
		SelectNodes []*selectNode `xml:"select"`
	}
	insertNode struct {
		ID         string `xml:"id,attr"`
		Datasource string `xml:"datasource,attr"`
		InnerXml   string `xml:",innerxml"`
	}
	deleteNode struct {
		ID         string `xml:"id,attr"`
		Datasource string `xml:"datasource,attr"`
		InnerXml   string `xml:",innerxml"`
	}
	updateNode struct {
		ID         string `xml:"id,attr"`
		Datasource string `xml:"datasource,attr"`
		InnerXml   string `xml:",innerxml"`
	}
	selectNode struct {
		ID         string `xml:"id,attr"`
		Datasource string `xml:"datasource,attr"`
		InnerXml   string `xml:",innerxml"`
	}
)

func (n *insertNode) getType() nodeType     { return insertType }
func (n *insertNode) getDatasource() string { return n.Datasource }
func (n *insertNode) GetID() string         { return n.ID }
func (n *insertNode) GetInnerXml() string   { return n.InnerXml }

func (n *deleteNode) getType() nodeType     { return deleteType }
func (n *deleteNode) getDatasource() string { return n.Datasource }
func (n *deleteNode) GetID() string         { return n.ID }
func (n *deleteNode) GetInnerXml() string   { return n.InnerXml }

func (n *updateNode) getType() nodeType     { return updateType }
func (n *updateNode) getDatasource() string { return n.Datasource }
func (n *updateNode) GetID() string         { return n.ID }
func (n *updateNode) GetInnerXml() string   { return n.InnerXml }

func (n *selectNode) getType() nodeType     { return selectType }
func (n *selectNode) getDatasource() string { return n.Datasource }
func (n *selectNode) GetID() string         { return n.ID }
func (n *selectNode) GetInnerXml() string   { return n.InnerXml }

var (
	queryLog = func(name, sql string, ps []any) {
		anorm.Logger.Debug([]*anorm.LogF{anorm.LogField("Name", name), anorm.LogField("SQL", sql), anorm.LogField("Parameter", ps)}, "")
	}
	queryErrorLog = func(err error, name, sql string, ps []any) {
		if err != nil {
			anorm.Logger.Error([]*anorm.LogF{anorm.LogField("Name", name), anorm.LogField("SQL", sql), anorm.LogField("Parameter", ps)}, "err: %v", err)
		}
	}
)
