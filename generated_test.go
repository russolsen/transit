// Copyright 2016 Russ Olsen. All Rights Reserved.
// 
// This code is a C# port of the Java version created and maintained by Cognitect, therefore:
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
	"github.com/russolsen/ohyeah"
	"testing"
)

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

func TestGeneratedMaps(t *testing.T) {
	r := ohyeah.RandomFunc(99)

	names := []interface{}{"foo", "bar", "baz"}
	symg := ohyeah.RepeatGen(SymbolGen(ohyeah.ElementGen(r, names)), 40)
	keyg := KeywordGen(ohyeah.ElementGen(r, names))

	vg := ohyeah.CycleGen(ohyeah.IntGen(r), ohyeah.ConstantGen(1234500),
		symg,
		keyg,
		ohyeah.PatternedStringGen(r, "val"),
		ohyeah.ConstantGen(Keyword("hello")))

	g := ohyeah.MapGen(r, KeywordGen(ohyeah.PatternedStringGen(r, "key")), ohyeah.ArrayGen(r, vg, 10), 2000)

	for i := 0; i < 4; i++ {
		value := g()
		VerifyRoundTrip(t, value)
	}

}

func TestSets(t *testing.T) {
	r := ohyeah.RandomFunc(99)

	symg := SymbolGen(ohyeah.PatternedStringGen(r, "key"))

	sg := ohyeah.RepeatGen(SetGen(r, symg, 2000), 1)
	g := ohyeah.ArrayGen(r, sg, 100)

	for i := 0; i < 40; i++ {
		value := g()
		VerifyRoundTrip(t, value)
	}
}
