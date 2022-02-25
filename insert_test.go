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

func TestInsertSave(t *testing.T) {
	user := userModel{
		Name:    "hugo",
		Age:     20,
		Address: "wuhan",
		Phone:   "13900110121",
	}
	err := New(new(userModel)).Insert().Save(&user)
	t.Log(user)
	errTest(t, err)
}

func TestInsertSaveAll(t *testing.T) {
	user1 := userModel{Name: "hugo1", Age: 20, Address: "wuhan", Phone: "13900110121"}
	user2 := userModel{Name: "hugo2", Age: 21, Address: "wuhan", Phone: "13900110122"}
	user3 := userModel{Name: "hugo3", Age: 22, Address: "wuhan", Phone: "13900110123"}
	err := New(new(userModel)).Insert().SaveAll(true, &user1, &user2, &user3)
	t.Log(err)
}

func TestInsertSaveBatch(t *testing.T) {
	user1 := userModel{Name: "hugo1", Age: 20, Address: "wuhan", Phone: "13900110122"}
	user2 := userModel{Name: "hugo2", Age: 21, Address: "wuhan", Phone: "13900110122"}
	user3 := userModel{Name: "hugo3", Age: 22, Address: "wuhan", Phone: "13900110123"}
	err := New(new(userModel)).Insert().SaveBatch(&user1, &user2, &user3)
	t.Log(err)
}
