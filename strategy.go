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

type Strategy int

const (
	// Default defines default not changed
	//
	// UserName  => UserName
	Default Strategy = iota
	// Underline defines to transform a underline style name
	//
	// UserName => user_name
	Underline
	// CamelCase defines to transform a camelcase style name
	//
	// UserName => userName
	CamelCase
)
