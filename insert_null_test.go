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

func TestInsertNullSave(t *testing.T) {
	user := userModelNull{
		Name:    NullString("hugo"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110121"),
	}
	err := New(new(userModelNull)).Insert().Save(&user)
	t.Log(user)
	errTest(t, err)
}

func TestInsertNullSaveAll(t *testing.T) {
	user1 := userModelNull{
		Name:    NullString("hugo"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110121"),
	}
	user2 := userModelNull{
		Name:    NullString("hugo2"),
		Age:     NullInt32(21),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user3 := userModelNull{
		Name:    NullString("hugo3"),
		Age:     NullInt32(22),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110123"),
	}
	err := New(new(userModelNull)).Insert().SaveAll(true, &user1, &user2, &user3)
	t.Log(err)
}

func TestInsertNullSaveBatch(t *testing.T) {
	user1 := userModelNull{
		Name:    NullString("hugo1"),
		Age:     NullInt32(20),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user2 := userModelNull{
		Name:    NullString("hugo2"),
		Age:     NullInt32(21),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110122"),
	}
	user3 := userModelNull{
		Name:    NullString("hugo1"),
		Age:     NullInt32(22),
		Address: NullString("wuhan"),
		Phone:   NullString("13900110123"),
	}
	err := New(new(userModelNull)).Insert().SaveBatch(&user1, &user2, &user3)
	t.Log(err)
}
