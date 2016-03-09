// Copyright 2016 Russ Olsen. All Rights Reserved.
// 
// This code is a Go port of the Java version created and maintained by Cognitect, therefore:
//
// Copyright 2014 Cognitect. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS-IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transit

import("fmt")

// MatchF is an equality function protocol used by
// sets and cmaps.
type MatchF func(a, b interface{}) bool


// Matches keys with a simple == test. Satisfies the
// MatchF protocol.
func Equals(a, b interface{}) bool {
	return a == b
}


// A tag id represents a #tag in the transit protocol.
type TagId string


func (t TagId) String() string {
	return fmt.Sprintf("[Tag: %s]", string(t))
}

// TaggedValue is a simple struct to hold the data from
// a transit #tag.

type TaggedValue struct {
	Tag   TagId
	Value interface{}
}

// A Keyword is a transit keyword, really just a string by another type.
type Keyword string

func NewKeyword(s string) Keyword {
	return Keyword(s)
}

func (k Keyword) String() string {
	return fmt.Sprintf(":%s", string(k))
}


// A Symbol is a transit symbol, really just a string by another type.
type Symbol string

func NewSymbol(s string) Symbol {
	return Symbol(s)
}

func (s Symbol) String() string {
	return string(s)
}

