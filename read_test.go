// Copyright 2016 Russ Olsen. All Rights Reserved.
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
//
// Much of this file was adopted from the Cognitect Java transit implementation.

package transit

import (
	"encoding/base64"
	"github.com/pborman/uuid"
	"container/list"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		t.Errorf("Expected %v[%v] == %v[%v]", v1, reflect.TypeOf(v1), v2, reflect.TypeOf(v2))
	}
}

func assertTrue(t *testing.T, v interface{}) {
	assertEquals(t, v, true)
}

func assertFalse(t *testing.T, v interface{}) {
	assertEquals(t, v, false)
}

func DecodeTransit(t *testing.T, s string) interface{} {
	value, err := DecodeFromString(s)
	if err != nil {
		t.Errorf("Error decoding transit string: %v: %v", s, err)
		return nil
	}
	return value
}

func TestReadString(t *testing.T) {
	assertEquals(t, "foo", DecodeTransit(t, "\"foo\""))
	assertEquals(t, "~foo", DecodeTransit(t, "\"~~foo\""))
	assertEquals(t, "`foo", DecodeTransit(t, "\"~`foo\""))
	assertEquals(t, "^foo", DecodeTransit(t, "\"~^foo\""))
}

func TestReadBoolean(t *testing.T) {

	assertTrue(t, DecodeTransit(t, "\"~?t\""))
	assertFalse(t, DecodeTransit(t, "\"~?f\""))
}

func TestReadNull(t *testing.T) {

	assertEquals(t, nil, DecodeTransit(t, "\"~_\""))
}

func TestReadKeyword(t *testing.T) {

	v := DecodeTransit(t, "\"~:foo\"")
	assertEquals(t, ":foo", v.(Keyword).String())
}

func TestReadInteger(t *testing.T) {

	i := DecodeTransit(t, "\"~i42\"")
	assertEquals(t, int64(42), i)

	j := DecodeTransit(t, "\"~n1234\"").(*big.Int)
	assertEquals(t, j.Int64(), int64(1234))
}

func TestReadDouble(t *testing.T) {
	assertEquals(t, 42.5, DecodeTransit(t, "\"~d42.5\""))
}

func TestReadSpecialNumbers(t *testing.T) {
	assertTrue(t, math.IsNaN(DecodeTransit(t, "\"~zNaN\"").(float64)))
	assertTrue(t, math.IsInf(DecodeTransit(t, "\"~zINF\"").(float64), 1))
	assertTrue(t, math.IsInf(DecodeTransit(t, "\"~z-INF\"").(float64), -1))
}

func TestReadBigDecimal(t *testing.T) {
	bd := DecodeTransit(t, "\"~f42.5\"").(*big.Float)
	x, _ := bd.Float64()
	assertTrue(t, x-42.5 < 0.001)
}

func TestReadUUID(t *testing.T) {
	u := uuid.Parse("07886363-98EC-4266-BE51-E09539AEE2A0")
	from_transit := DecodeTransit(t, "\"~u"+u.String()+"\"").(uuid.UUID)
	assertEquals(t, u.String(), from_transit.String())
}

func TestReadURI(t *testing.T) {
	from_transit := DecodeTransit(t, "\"~rhttp://www.foo.com\"").(*url.URL)
	assertEquals(t, from_transit.String(), "http://www.foo.com")

}

func TestReadSymbol(t *testing.T) {
	sym := DecodeTransit(t, "\"~$foo\"").(Symbol)
	assertEquals(t, "foo", sym.String())
}

func TestReadCharacter(t *testing.T) {
	assertEquals(t, 'f', DecodeTransit(t, "\"~cf\"").(int32))
}

func TestReadUnknown(t *testing.T) {
	fooThing := DecodeTransit(t, "\"~jfoo\"").(TaggedValue)
	assertEquals(t, fooThing.Tag, TagId("j"))
	assertEquals(t, fooThing.Value, "foo")

	pointThing := DecodeTransit(t, "{\"~#point\":[1,2]}").(TaggedValue)
	assertEquals(t, pointThing.Tag, TagId("point"))

	value := pointThing.Value.([]interface{})
	assertEquals(t, value[0], int64(1))
	assertEquals(t, value[1], int64(2))
}

func TestReadArray(t *testing.T) {
	l := DecodeTransit(t, "[1, 2, 3]").([]interface{})

	assertEquals(t, 3, len(l))
	assertEquals(t, l[0], int64(1))
	assertEquals(t, l[1], int64(2))
	assertEquals(t, l[2], int64(3))
}

func TestReadBinary(t *testing.T) {
	value := []byte("foobarbaz")

	b64 := base64.StdEncoding.EncodeToString(value)

	decoded := DecodeTransit(t, "\"~b"+b64+"\"").([]uint8)

	assertEquals(t, len(value), len(decoded))

	for i := 0; i < len(value); i++ {
		assertEquals(t, value[i], decoded[i])
	}
}

    func TestReadMap(t *testing.T) {

        m := DecodeTransit(t, "{\"a\": 2, \"b\": 4}").(map[interface{}]interface{})

        assertEquals(t, 2, len(m))
        assertEquals(t, int64(2), m["a"].(int64))
        assertEquals(t, int64(4), m["b"].(int64))
    }

    func TestReadSet(t *testing.T) {

        s := DecodeTransit(t, `{"~#set": [1, 2, 3]}`).(*Set)

        assertEquals(t, 3, len(s.Contents))

        assertTrue(t, s.ContainsEq(int64(1)))
        assertTrue(t, s.ContainsEq(int64(2)))
        assertTrue(t, s.ContainsEq(int64(3)))
    }

    func TestReadList(t *testing.T) {

        l := DecodeTransit(t, `{"~#list": [1, 2, 3]}`).(*list.List)

        assertEquals(t, 3, l.Len())

        assertEquals(t, int64(1), l.Front().Value)
        assertEquals(t, int64(2), l.Front().Next().Value)
        assertEquals(t, int64(3), l.Front().Next().Next().Value)
    }

    func TestReadRatio(t *testing.T) {
        r := DecodeTransit(t, "{\"~#ratio\": [\"~n1\",\"~n2\"]}").(big.Rat)
	f64, _ := r.Float64()
        assertTrue(t, math.Abs(f64 - 0.5) < 0.001)
    }

    func TestReadCmap(t *testing.T) {
        //m := DecodeTransit(t, "{\"~#cmap\": [{\"~#ratio\":[\"~n1\",\"~n3\"]},1,{\"~#list\":[1,2,3]},2]}").(*CMap)

        m := DecodeTransit(t, `{"~#cmap": [{"~#ratio":["~n1","~n3"]},1,{"~#list":[1,2,3]},2]}`).(*CMap)

        //m := DecodeTransit(t, "{\"~#cmap\": [\"hello\", 1, \"BYE\", 2]}").(*CMap)

	key := m.Entries[0].Key.(big.Rat)
	value := m.Entries[0].Value.(int64)
	
	f64, _ := key.Float64()
        assertTrue(t, math.Abs(f64 - 0.333333333) < 0.001)
	assertEquals(t, value, int64(1))
    }

