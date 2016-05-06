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

import (
	"container/list"
	"time"
	"github.com/pborman/uuid"
	"github.com/russolsen/ohyeah"
	"testing"
)

var Times = []interface{} {
	time.Unix(0, 0)}

func TimeGen(r ohyeah.Int64F) ohyeah.Generator {
	return ohyeah.ElementGen(r, Times)
}

var Uuids = []interface{}{
	uuid.Parse("6E4BA181-E528-4676-84A2-87974DEBBE90"),
	uuid.Parse("6E4BA181-E528-4676-84A2-87974DEBBE91"),
	uuid.Parse("6E4BA181-E528-4676-84A2-87974DEBBE92"),
	uuid.Parse("D5E0D599-83F1-47E9-9A27-73E4131590D8"),
	uuid.Parse("D5E0D599-83F1-47E9-9A27-73E4131590E8")}

func UuidGen(r ohyeah.Int64F) ohyeah.Generator {
	return ohyeah.ElementGen(r, Uuids)
}

var Numbers = []interface{}{
	int8(8), int16(16), int32(32), int64(64),
	uint8(8), uint16(16), uint32(32), uint64(64)}

func NumberGen(r ohyeah.Int64F) ohyeah.Generator {
	return ohyeah.ElementGen(r, Numbers)
}

func KeywordGen(stringGenerator ohyeah.Generator) ohyeah.Generator {
	return func() interface{} {
		s := stringGenerator().(string)
		return Keyword(s)
	}
}

func SymbolGen(stringGenerator ohyeah.Generator) ohyeah.Generator {
	return func() interface{} {
		s := stringGenerator().(string)
		return Symbol(s)
	}
}

func SetGen(r ohyeah.Int64F, elementGenerator ohyeah.Generator, n int) ohyeah.Generator {
	ag := ohyeah.ArrayGen(r, elementGenerator, n)
	return func() interface{} {
		array := ag().([]interface{})
		return NewSet(array)
	}
}

func ListGen(r ohyeah.Int64F, elementGenerator ohyeah.Generator, n int) ohyeah.Generator {
	return func() interface{} {
		n := ohyeah.IntN(r, n)
		lst := list.New()

		for i := 0; i < n; i++ {
			lst.PushBack(elementGenerator())
		}
		return lst
	}
}

func SimpleGen(r ohyeah.Int64F) ohyeah.Generator {
	names := []interface{}{"foo", "bar", "baz", "apple", "organge", "red", "x"}
	strg := ohyeah.ElementGen(r, names)
	symg := ohyeah.RepeatGen(strg, 40)
	keyg := KeywordGen(strg)

	return ohyeah.CycleGen(
		ohyeah.IntGen(r),
		ohyeah.BigRatGen(r),
		ohyeah.BigIntGen(r),
		ohyeah.BigFloatGen(r),
		ohyeah.ConstantGen(1234500),
		ohyeah.RuneGen(r),
		strg,
		symg,
		keyg,
		NumberGen(r),
		UuidGen(r),
		TimeGen(r),
		ohyeah.PatternedStringGen("val"),
		ohyeah.ConstantGen(Keyword("hello")))
}

func TestSimpleValues(t *testing.T) {
	r := ohyeah.RandomFunc(99)
	f := SimpleGen(r)

	for i := 0; i < 50; i++ {
		value := f()
		VerifyRoundTrip(t, value)
	}
}

func TestGeneratedMaps(t *testing.T) {
	r := ohyeah.RandomFunc(99)
	vg := SimpleGen(r)

	g := ohyeah.MapGen(r, KeywordGen(ohyeah.PatternedStringGen("key")), ohyeah.ArrayGen(r, vg, 10), 2000)

	for i := 0; i < 4; i++ {
		value := g()
		VerifyRoundTrip(t, value)
	}

}

func TestLists(t *testing.T) {
	r := ohyeah.RandomFunc(99)

	sg := SimpleGen(r)

	g := ListGen(r, sg, 10)

	for i := 0; i < 40; i++ {
		value := g()
		VerifyRoundTrip(t, value)
	}
}

func TestSets(t *testing.T) {
	r := ohyeah.RandomFunc(99)

	symg := SymbolGen(ohyeah.PatternedStringGen("key"))

	sg := ohyeah.RepeatGen(SetGen(r, symg, 2000), 1)
	g := ohyeah.ArrayGen(r, sg, 100)

	for i := 0; i < 40; i++ {
		value := g()
		VerifyRoundTrip(t, value)
	}
}
