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
	"container/list"
	"encoding/base64"
	"github.com/pborman/uuid"
	"math"
	"math/big"
	"net/url"
	"strconv"
	"time"
)

// DecodeKeyword decodes ~: style keywords.
func DecodeKeyword(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	var result = Keyword(s)
	return result, nil
}

// DecodeKeyword decodes ~$ style symbols.
func DecodeSymbol(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	var result = Symbol(s)
	return result, nil
}

// DecodeKeyword decodes ` reserved values. Actually it just nil.
func DecodeReserved(d Decoder, x interface{}) (interface{}, error) {
	return nil, nil
}

// DecodeTag decodes #tag style tags. Note that this function just handles
// the tag itself.
func DecodeTag(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	return TagId(s), nil
}

// DecodeTaggedValue handles an entire tagged value by returning it as more or less unknown.
func DecodeTaggedValue(d Decoder, x interface{}) (interface{}, error) {
	return x, nil
}


func DecodeCMap(d Decoder, x interface{}) (interface{}, error) {
	array := x.(TaggedValue).Value.([]interface{})

	var result = NewCMap()

	l := len(array)

	for i := 0; i < l; i += 2 {
		key := array[i]
		value := array[i+1]
		result.Append(key, value)
	}

	return result, nil
}

// DecodeSet decodes a transit set into a transit.Set instance.
func DecodeSet(d Decoder, x interface{}) (interface{}, error) {
	tagged := x.(TaggedValue)
	values := (tagged.Value).([]interface{})
	result := NewSet(values)
	return result, nil
}

// DecodeList decodes a transit list into a Go list.
func DecodeList(d Decoder, x interface{}) (interface{}, error) {
	tagged := x.(TaggedValue)
	values := (tagged.Value).([]interface{})
	result := list.New()
	for _, item := range values {
		result.PushBack(item)
	}
	return result, nil
}

// DecodeQuote decodes a transit quoted value by simply returning the value.
func DecodeQuote(d Decoder, x interface{}) (interface{}, error) {
	tagged := x.(TaggedValue)
	return tagged.Value, nil
}

// DecodeTilde decodes an escaped string by stripping off the leading ~.
func DecodeTilde(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	return s[1:], nil
}

// DecodeRFC3339 decodes a time value into a Go time instance.
// TBD not 100% this covers all possible values.
func DecodeRFC3339(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	var result, err = time.Parse(time.RFC3339Nano, s)
	return result, err
}

// DecodeTime decodes a time value represended as millis since 1970.
func DecodeTime(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	var millis, _ = strconv.ParseInt(s[2:], 10, 64)
	var seconds = millis / 1000
	var nanos = (millis % 1000) * 1000
	result := time.Unix(seconds, nanos)
	return result, nil
}

// DecodeBoolean decodes a transit boolean into a Go bool.
func DecodeBoolean(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	if s == "t" {
		return true, nil
	} else if s == "f" {
		return false, nil
	} else {
		return nil, &TransitError{Message: "Unknow boolean value"}
	}
}

// DecodeBigInteger decodes a transit big integer into a Go big.Int.
func DecodeBigInteger(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	result := new(big.Int)
	result.SetString(s, 10)
	return result, nil
}

// DecodeBigInteger decodes a transit integer into a plain Go int64
func DecodeInteger(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	result, err := strconv.ParseInt(s, 10, 64)
	return result, err
}

func newRational(a, b *big.Int) *big.Rat {
	var r = big.NewRat(1, 1)
	r.SetFrac(a, b)
	return r
}

func toBigInt(x interface{}) *big.Int {
	switch v := x.(type) {
	default:
		return big.NewInt(0)
	case *big.Int:
		return v
	case int64:
		return big.NewInt(v)
	}
}

// DecodeRatio decodes a transit ratio into a Go big.Rat.
func DecodeRatio(d Decoder, x interface{}) (interface{}, error) {
	tagged := x.(TaggedValue)
	values := (tagged.Value).([]interface{})
	a := toBigInt(values[0])
	b := toBigInt(values[1])
	result := newRational(a, b)
	return *result, nil
}

// DecodeRatio decodes a transit char.
func DecodeChar(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	return rune(s[0]), nil
}

// DecodeRatio decodes a transit decimal into an int64.
func DecodeDecimal(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	return strconv.ParseFloat(s, 64)
}

// DecodeRatio decodes a transit big decimal into an float64.
func DecodeBigDecimal(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	result, _, _ := big.ParseFloat(s, 10, 20, big.ToZero)
	return result, nil
}

// DecodeRatio decodes a transit null/nil.
func DecodeNil(d Decoder, x interface{}) (interface{}, error) {
	return nil, nil
}

// DecodeRatio decodes a transit base64 encoded byte array into a
// Go byte array.
func DecodeByte(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	return base64.StdEncoding.DecodeString(s)
}

// DecodeURI decodes a transit URI into an instance of net/Url.
// Despite the name, Go Urls are almost URIs.
func DecodeURI(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	return url.Parse(s)
}

// DecodeUUID decodes a transit UUID into an instance of net/UUID
func DecodeUUID(d Decoder, x interface{}) (interface{}, error) {
	s := x.(string)
	var u = uuid.Parse(s)
	return u, nil
}

// DecodeSpecialNumber decodes NaN, INF and -INF into their Go equivalents.
func DecodeSpecialNumber(d Decoder, x interface{}) (interface{}, error) {
	tag := x.(string)
	if tag == "NaN" {
		return math.NaN(), nil
	} else if tag == "INF" {
		return math.Inf(1), nil
	} else if tag == "-INF" {
		return math.Inf(-1), nil
	} else {
		return nil, &TransitError{Message: "Bad special number:"}
	}
}

// DecodeUnknown decodes a tag that we don't otherwise recognize.
func DecodeUnknown(d Decoder, x interface{}) (interface{}, error) {
	return x, nil
}

