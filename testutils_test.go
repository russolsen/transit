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
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/russolsen/same"
	"github.com/shamaton/msgpack"
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

func ExemplarPath(fileName string) string {
	return "transit-format/examples/0.8/simple/" + fileName
}

func toUrl(s string) *url.URL {
	url, _ := url.Parse(s)
	return url
}

func VerifyRoundTrip(t testing.TB, value interface{}) interface{} {

	var newValue interface{}

	for _, verbose := range []bool{true, false} {
		json, err := EncodeToString(value, verbose)

		if err != nil {
			t.Errorf("Error encoding to Transit %v: %v", value, err)
		}

		newValue, err = DecodeFromString(json)

		if err != nil {
			t.Errorf("Error decoding %v: %v.\nJson:\n%v", value, err, json)
		}

		if !same.IsSame(value, newValue) {
			t.Errorf("Round trip values do not match.\nValue:[%v]\n%v\nNew value:[%v]\n %v\nJson:\n%v",
				value, reflect.TypeOf(value), newValue, reflect.TypeOf(newValue), json)
		}
	}

	return newValue
}

func VerifyExemplar(t testing.TB, transitValue interface{}, exemplarPath string) {
	t.Helper()
	f, err := os.Open(exemplarPath)

	if err != nil {
		t.Errorf("%v", err)
		return
	}

	exemplarValue, err := NewDecoder(f).Decode()

	if err != nil {
		t.Errorf("Error decoding exemplar JSON [%v]: %v", exemplarPath, err)
		return
	}

	// Use the more stringent DeepEqual here because both values were read with
	// Transit and should therefore have the same type.

	if !reflect.DeepEqual(exemplarValue, transitValue) {
		t.Errorf("Value read from exemplar file [%v]:\n%v\nDoes not match round trip value:\n%v",
			exemplarPath, exemplarValue, transitValue)
	}

	mpkg, err := os.Open(strings.Replace(exemplarPath, ".json", ".mp", 1))
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	msg, err := ioutil.ReadAll(mpkg)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	var buf bytes.Buffer
	if err := NewMsgPackEncoder(&buf).Encode(exemplarValue); err != nil {
		t.Errorf("%v", err)
		return
	}

	if buf.String() != string(msg) {
		t.Errorf("(%s) %v, \n%v \n!= \n%v", exemplarPath, exemplarValue, buf.Bytes(), msg)
		return
	}

	var x interface{}
	if err = msgpack.Decode(buf.Bytes(), &x); err != nil {
		t.Errorf("%v", err)
		return
	}

	if !reflect.DeepEqual(x, transitValue) {
		t.Errorf("Value read from exemplar file [%v]:\n%#v\nDoes not match round trip value:\n%#v",
			exemplarPath, x, transitValue)
	}
}

func VerifyJson(t *testing.T, transit string, path string) error {
	var actual interface{}
	var err error

	if err = json.Unmarshal([]byte(transit), &actual); err != nil {
		t.Errorf("Error reading generated transit %v: %v", transit, err)
		return err
	}

	expectedJson, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Error reading exemplare [%v]: %v", path, err)
		return err
	}

	var expected interface{}
	if err = json.Unmarshal([]byte(expectedJson), &expected); err != nil {
		t.Errorf("Error decoding expected [%v]: %v", path, err)
		return err
	}

	if !same.IsSame(actual, expected) {
		t.Errorf("Actual does not match exemplar.\nExpected: %v\nActual: %v", expected, actual)
		return NewTransitError("No match", actual)
	}
	return nil
}

func Verify(t testing.TB, value interface{}, exemplarPath string) {
	t.Helper()

	transit := VerifyRoundTrip(t, value)
	VerifyExemplar(t, transit, exemplarPath)
}

func VerifyReadError(t *testing.T, badTransitJson string) {
	_, err := DecodeFromString(badTransitJson)
	if err == nil {
		t.Errorf("Expected that decoding [%v] would generate an error, but it did not.", badTransitJson)
	}
}
