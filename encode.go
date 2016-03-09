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
	"bytes"
	"container/list"
	"io"
	"net/url"
	"reflect"
	"math/big"
	"github.com/pborman/uuid"
	"time"
)

type Encoder struct {
	emitter       DataEmitter
	valueEncoders map[interface{}]ValueEncoder
}

var goListType = reflect.TypeOf(list.New())
var keywordType = reflect.TypeOf(NewKeyword(""))
var symbolType = reflect.TypeOf(NewSymbol(""))
var cmapType = reflect.TypeOf(NewCMap())

var aUrl, _ = url.Parse("http://foo.com")
var urlType = reflect.TypeOf(aUrl)
var setType = reflect.TypeOf(Set{})

var timeType = reflect.TypeOf(time.Now())
var bigIntType = reflect.TypeOf(*big.NewInt(int64(1)))
var bigFloatType = reflect.TypeOf(*big.NewFloat(float64(1.)))
var uuidType = reflect.TypeOf(uuid.NewRandom())
var linkType = reflect.TypeOf(*NewLink())
var taggedValueType = reflect.TypeOf(TaggedValue{TagId("#foo"), 1})

var runeType = reflect.TypeOf('x')
var nilValue = reflect.ValueOf(nil)
var nilEncoder = NewNilEncoder()

// NewEncoder creates a new encoder set to writ to the stream supplied.
func NewEncoder(w io.Writer) *Encoder {
	valueEncoders := make(map[interface{}]ValueEncoder)

	emitter := NewJsonEmitter(w)
	e := Encoder{emitter: emitter, valueEncoders: valueEncoders}

	e.addHandler(reflect.String, NewStringEncoder())

	e.addHandler(reflect.Bool, NewBoolEncoder())
	e.addHandler(reflect.Ptr, NewPointerEncoder())

	floatEncoder := NewFloatEncoder()

	e.addHandler(reflect.Float32, floatEncoder)
	e.addHandler(reflect.Float64, floatEncoder)

	intEncoder := NewIntEncoder()

	e.addHandler(reflect.Int, intEncoder)
	e.addHandler(reflect.Int8, intEncoder)
	e.addHandler(reflect.Int16, intEncoder)
	e.addHandler(reflect.Int32, intEncoder)
	e.addHandler(reflect.Int64, intEncoder)

	uintEncoder := NewUintEncoder()

	e.addHandler(reflect.Uint, uintEncoder)
	e.addHandler(reflect.Uint8, uintEncoder)
	e.addHandler(reflect.Uint16, uintEncoder)
	e.addHandler(reflect.Uint32, uintEncoder)
	e.addHandler(reflect.Uint64, uintEncoder)

	arrayEncoder := NewArrayEncoder()

	e.addHandler(reflect.Array, arrayEncoder)
	e.addHandler(reflect.Slice, arrayEncoder)
	e.addHandler(reflect.Map, NewMapEncoder())

	// Link

	e.addHandler(runeType, NewRuneEncoder())
	e.addHandler(timeType, NewTimeEncoder())
	e.addHandler(uuidType, NewUuidEncoder())
	e.addHandler(bigIntType, NewBigIntEncoder())
	e.addHandler(bigFloatType, NewBigDecimalEncoder())
	e.addHandler(goListType, NewListEncoder())
	e.addHandler(symbolType, NewSymbolEncoder())
	e.addHandler(keywordType, NewKeywordEncoder())
	e.addHandler(cmapType, NewCMapEncoder())
	e.addHandler(setType, NewSetEncoder())
	e.addHandler(urlType, NewUrlEncoder())
	e.addHandler(linkType, NewLinkEncoder())

	e.addHandler(taggedValueType, NewTaggedValueEncoder())

	return &e
}

// AddHandler adds a new handler to the table used by this encoder
// for encoding values. The t value should be an instance
// of reflect.Type and the c value should be an encoder for that type.
func (e Encoder) AddHandler(t reflect.Type, c ValueEncoder) {
	e.addHandler(t, c)
}

// addHandler adds a new handler to the table, but the untyped first
// parameter lets you enter either reflect.Type or reflect.Kind values.
// Used internally.
func (e Encoder) addHandler(t interface{}, c ValueEncoder) {
	e.valueEncoders[t] = c
}

// ValueEncoderFor finds the encoder for the given value.
func (e Encoder) ValueEncoderFor(v reflect.Value) ValueEncoder {
	// Nil is a special case since it doesn't really work
	// very well with the reflect package.

	if v == nilValue {
		return nilEncoder
	}

	// Look for an encoder by the specific type.

	typeEncoder := e.valueEncoders[v.Type()]
	if typeEncoder != nil {
		return typeEncoder
	}

	// If we can't find a type encoder, try finding one
	// by type. This is will catch values of know kinds,
	// say int64 or string which have a different specific
	// type.

	kindEncoder := e.valueEncoders[v.Kind()]
	if kindEncoder != nil {
		return kindEncoder
	}

	// No encoder, for this type, return the error encoder.
	return NewErrorEncoder()
}

// Given a Value, encode it.
func (e Encoder) EncodeValue(v reflect.Value, asKey bool) error {
	valueEncoder := e.ValueEncoderFor(v)
	return valueEncoder.Encode(e, v, asKey)
}

// Given a raw interface, encode it.
func (e Encoder) EncodeInterface(x interface{}, asKey bool) error {
	v := reflect.ValueOf(x)
	return e.EncodeValue(v, asKey)
}

// Encode a value at the top level.
func (e Encoder) Encode(x interface{}) error {
	v := reflect.ValueOf(x)
	valueEncoder := e.ValueEncoderFor(v)

	if valueEncoder.IsStringable(v) {
		x = TaggedValue{TagId("'"), x}
	} 

	return e.EncodeInterface(x, false)
}

// Encode the given value to a string.
func EncodeToString(x interface{}) (string, error) {
	var buf bytes.Buffer
	var encoder = NewEncoder(&buf)
	err := encoder.Encode(x)

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
