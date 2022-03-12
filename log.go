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

package anorm

import "sync"

var logMu = &sync.Mutex{}

type execLog struct {
	callName string
	sql      string
	ps       []interface{}
}

func (el *execLog) Log() {
	logMu.Lock()
	defer logMu.Unlock()
	Configuration.Logger.Printf("%s================================================>\n", el.callName)
	Configuration.Logger.Printf("sql{%s}\n", el.sql)
	Configuration.Logger.Printf("parameters{%v}\n", el.ps)
	Configuration.Logger.Printf("<================================================%s\n", el.callName)
}

func debug() bool {
	return Configuration.Logger != nil && Configuration.Debug
}

func debugLog(str string, ps ...interface{}) {
	if debug() {
		Configuration.Logger.Printf(str, ps...)
	}
}

func handleErr(err error) {
	if Configuration.Logger != nil && Configuration.Debug {
		if err != nil {
			Configuration.Logger.Fatalf("%v", err)
		}
	}
}
