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
	"fmt"
	"log"
	"net/url"
	"reflect"
)

// ValueEncoder is the interface for objects that know how to
// transit encode a single value.
type ValueEncoder interface {
	IsStringable(reflect.Value) bool
	Encode(e Encoder, value reflect.Value, asString bool) error
}

type NilEncoder struct{}

func NewNilEncoder() *NilEncoder {
	return &NilEncoder{}
}

func (ie NilEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie NilEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	return e.emitter.EmitNil(asKey)
}

type PointerEncoder struct{}

func NewPointerEncoder() *PointerEncoder {
	return &PointerEncoder{}
}

func (ie PointerEncoder) IsStringable(v reflect.Value) bool {
	return false
}

func (ie PointerEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	//log.Println("*** Defer pointer to:", v.Elem())
	return e.EncodeInterface(v.Elem().Interface(), asKey)
}

type BoolEncoder struct{}

func NewBoolEncoder() *BoolEncoder {
	return &BoolEncoder{}
}

func (ie BoolEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie BoolEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	b := v.Bool()
	return e.emitter.EmitBool(b, asKey)
}

type FloatEncoder struct{}

func NewFloatEncoder() *FloatEncoder {
	return &FloatEncoder{}
}

func (ie FloatEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie FloatEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	f := v.Float()
	return e.emitter.EmitFloat(f, asKey)
}

type IntEncoder struct{}

func NewIntEncoder() *IntEncoder {
	return &IntEncoder{}
}

func (ie IntEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie IntEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	return e.emitter.EmitInt(v.Int(), asKey)
}

type UintEncoder struct{}

func NewUintEncoder() *UintEncoder {
	return &UintEncoder{}
}

func (ie UintEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie UintEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	return e.emitter.EmitInt(int64(v.Uint()), asKey)
}

type KeywordEncoder struct{}

func NewKeywordEncoder() *KeywordEncoder {
	return &KeywordEncoder{}
}

func (ie KeywordEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie KeywordEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	s := v.String()
	//log.Println("Encoding keyword:", s)
	return e.emitter.EmitString("~:"+s, true)
}

type SymbolEncoder struct{}

func NewSymbolEncoder() *SymbolEncoder {
	return &SymbolEncoder{}
}

func (ie SymbolEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie SymbolEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	s := v.String()
	//log.Println("Encoding symbol:", s)
	return e.emitter.EmitString("~$"+s, true)
}

type StringEncoder struct{}

func NewStringEncoder() *StringEncoder {
	return &StringEncoder{}
}

func (ie StringEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie StringEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	s := v.String()
	if len(s) > 0 && s[0:1] == "~" {
		s = "~" + s
	}
	return e.emitter.EmitString(s, asKey)
}

type UrlEncoder struct{}

func NewUrlEncoder() *UrlEncoder {
	return &UrlEncoder{}
}

func (ie UrlEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie UrlEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	log.Println("Encode url::::", v)

	u := v.Interface().(*url.URL)

	var us string

	us = fmt.Sprintf("%s://%s/%s?%s#%s",
		u.Scheme,
		u.Host,
		u.Path,
		u.Query,
		u.Fragment)

	log.Println("Encoded:", us)
	return e.emitter.EmitString(fmt.Sprintf("~r%s", us), asKey)
}

type ErrorEncoder struct{}

func NewErrorEncoder() *ErrorEncoder {
	return &ErrorEncoder{}
}

func (ie ErrorEncoder) IsStringable(v reflect.Value) bool {
	return true
}

func (ie ErrorEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	return NewTransitError("Dont know how to encode value", v)
}

type ArrayEncoder struct{}

func NewArrayEncoder() *ArrayEncoder {
	return &ArrayEncoder{}
}

func (ie ArrayEncoder) IsStringable(v reflect.Value) bool {
	return false
}

func (ie ArrayEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	e.emitter.EmitStartArray()

	l := v.Len()
	for i := 0; i < l; i++ {
		if i > 0 {
			e.emitter.EmitArraySeparator()
		}
		element := v.Index(i)
		err := e.EncodeInterface(element.Interface(), asKey)
		if err != nil {
			return err
		}
	}

	return e.emitter.EmitEndArray()
}

type MapEncoder struct{}

func NewMapEncoder() *MapEncoder {
	return &MapEncoder{}
}

func (me MapEncoder) IsStringable(v reflect.Value) bool {
	return false
}

func (me MapEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	//log.Println("Map Encoder for", v)

	keys := KeyValues(v)
	if me.allStringable(e, keys) {
		return me.encodeNormalMap(e, v)
	} else {
		return me.encodeCompositeMap(e, v)
	}
}

func (me MapEncoder) allStringable(e Encoder, keys []reflect.Value) bool {
	for _, key := range keys {
		valueEncoder := e.ValueEncoderFor(reflect.ValueOf(key.Interface()))
		if !valueEncoder.IsStringable(key) {
			return false
		}
	}
	return true
}

func (me MapEncoder) encodeCompositeMap(e Encoder, v reflect.Value) error {
	e.emitter.EmitStartArray()

	e.emitter.EmitTag("cmap")

	keys := KeyValues(v)

	for _, key := range keys {
		e.emitter.EmitArraySeparator()

		err := e.EncodeValue(key, true)

		if err != nil {
			return err
		}

		e.emitter.EmitArraySeparator()

		value := GetMapElement(v, key)
		err = e.EncodeValue(value, false)

		if err != nil {
			return err
		}
	}

	return e.emitter.EmitEndArray()
}

func (me MapEncoder) encodeNormalMap(e Encoder, v reflect.Value) error {
	//l := v.Len()
	e.emitter.EmitStartArray()

	e.emitter.EmitString("^ ", false)

	keys := KeyValues(v)

	for _, key := range keys {
		e.emitter.EmitArraySeparator()

		err := e.EncodeValue(key, true)

		if err != nil {
			return err
		}

		e.emitter.EmitArraySeparator()

		value := GetMapElement(v, key)

		err = e.EncodeValue(value, false)

		if err != nil {
			return err
		}
	}

	return e.emitter.EmitEndArray()
}

type SetEncoder struct{}

func NewSetEncoder() *SetEncoder {
	return &SetEncoder{}
}

func (ie SetEncoder) IsStringable(v reflect.Value) bool {
	return false
}

func (ie SetEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	s := v.Interface().(Set)

	//log.Println("*** Encode set:", v)

	//l := v.Len()
	e.emitter.EmitStartArray()
	e.emitter.EmitTag("set")
	e.emitter.EmitArraySeparator()

	e.emitter.EmitStartArray()

	for i, element := range s.Contents {
		if i != 0 {
			e.emitter.EmitArraySeparator()
		}
		err := e.EncodeInterface(element, asKey)
		if err != nil {
			return err
		}
	}

	e.emitter.EmitEndArray()

	return e.emitter.EmitEndArray()
}

type ListEncoder struct{}

func NewListEncoder() *ListEncoder {
	return &ListEncoder{}
}

func (ie ListEncoder) IsStringable(v reflect.Value) bool {
	return false
}

func (ie ListEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	lst := v.Interface().(*list.List)

	//log.Println("*** Encode go list")

	//l := v.Len()
	e.emitter.EmitStartArray()
	e.emitter.EmitTag("list")

	for element := lst.Front(); element != nil; element = element.Next() {
		e.emitter.EmitArraySeparator()
		err := e.EncodeInterface(element.Value, asKey)
		if err != nil {
			return err
		}
	}

	return e.emitter.EmitEndArray()
}

type CMapEncoder struct{}

func NewCMapEncoder() *CMapEncoder {
	return &CMapEncoder{}
}

func (ie CMapEncoder) IsStringable(v reflect.Value) bool {
	return false
}

func (ie CMapEncoder) Encode(e Encoder, v reflect.Value, asKey bool) error {
	cmap := v.Interface().(*CMap)

	//l := v.Len()
	e.emitter.EmitStartArray()
	e.emitter.EmitTag("cmap")

	for _, entry := range cmap.Entries {
		e.emitter.EmitArraySeparator()

		err := e.EncodeInterface(entry.Key, true)
		if err != nil {
			return err
		}

		e.emitter.EmitArraySeparator()

		err = e.EncodeInterface(entry.Value, false)
		if err != nil {
			return err
		}
	}

	return e.emitter.EmitEndArray()
}
