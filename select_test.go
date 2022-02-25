// Copyright 2022 anox Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anox

import (
	"testing"
)

func init() {
	testInit()
}

func TestOrmSelectCount(t *testing.T) {
	o := New(new(userModel))
	count, err := o.Select().Count(&userModel{ID: 31})
	if err != nil {
		t.Error(err)
	}
	t.Logf("count is %d\n", count)
}

func TestOrmSelectPagination(t *testing.T) {
	o := New(new(userModel))
	models, count, err := o.Select().Pagination(&userModel{Name: "hugo2"}, MySQLPagination, 0, 1000)
	t.Log(count)
	t.Log(err)
	for i := range models {
		t.Log(models[i].(*userModel))
	}
}

func TestOrmSelectList(t *testing.T) {
	o := New(new(userModel))
	models, err := o.Select().List(&userModel{Name: "hugo"})
	t.Log(err)
	for i := range models {
		t.Log(models[i].(*userModel))
	}
}
