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

	"github.com/go-the-way/sg"
)

func init() {
	testInit()
}

func TestDeleteNullRemove(t *testing.T) {
	o := New(new(userModelNull))
	err := o.Delete().Where(sg.Eq("id", 15)).Remove(&userModelNull{Name: NullString("hugo2")})
	t.Log(err)
}
