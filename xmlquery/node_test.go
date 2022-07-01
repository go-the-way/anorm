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

func TestGetDS(t *testing.T) {
	if ds := getDS(&rootNode{Datasource: "haha"}, &insertNode{"", "haha2", ""}); ds != "haha2" {
		t.Fatal("test failed!")
	}
	if ds := getDS(&rootNode{Datasource: "haha"}, &insertNode{"", "", ""}); ds != "haha" {
		t.Fatal("test failed!")
	}
}
