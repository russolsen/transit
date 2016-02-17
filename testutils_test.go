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
	"os"
	"reflect"
	"encoding/json"
	"github.com/russolsen/same"
	"io/ioutil"
	"net/url"
	"testing"
)

func ExemplarPath(fileName string) string {
	return "transit-format/examples/0.8/simple/" + fileName
}


func toUrl(s string) *url.URL {
	url, _ := url.Parse(s)
	return url
}


func VerifyRoundTrip(t *testing.T, value interface{}) (interface{}, error) {
	json, err := EncodeToString(value)

	if err != nil {
		t.Errorf("Error encoding to Transit %v: %v", value, err)
		return nil, err
	}

	newValue, err := DecodeFromString(json)

	if err != nil {
		t.Errorf("Error decoding %v: %v", value, err)
		return nil, err
	}

	if !same.IsSame(newValue, value) {
		t.Errorf("Round trip values do not match.\nValue:[%v]\n%v\nNew value:[%v]\n %v",
			value, reflect.TypeOf(value), newValue, reflect.TypeOf(newValue))
		return newValue, err
	}

	return newValue, err
}

func VerifyExemplar(t *testing.T, transitValue interface{}, exemplarPath string) {
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

	//log.Println("***Expected json:", expectedJson)
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

func Verify(t *testing.T, value interface{}, exemplarPath string) {
	transit, _ := VerifyRoundTrip(t, value)
	VerifyExemplar(t, transit, exemplarPath)
}
