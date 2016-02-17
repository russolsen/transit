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
	"fmt"
	"log"
	"testing"
)

var exemplars map[string]interface{}

func init() {
	exemplars = make(map[string]interface{})

	// /*
	exemplars["small_strings.json"] = []string{
		"", "a", "ab", "abc", "abcd", "abcde", "abcdef"}

	exemplars["ints.json"] = []int64{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127}

	exemplars["keywords.json"] = []Keyword{
		Keyword("a"), Keyword("ab"), Keyword("abc"),
		Keyword("abcd"), Keyword("abcde"), Keyword("a1"),
		Keyword("b2"), Keyword("c3"), Keyword("a_b")}

	exemplars["symbols.json"] = []Symbol{
		Symbol("a"), Symbol("ab"), Symbol("abc"),
		Symbol("abcd"), Symbol("abcde"), Symbol("a1"),
		Symbol("b2"), Symbol("c3"), Symbol("a_b")}

	exemplars["doubles_interesting.json"] = []float64{
		-3.14159, 3.14159, 4.0E11, 2.998E8, 6.626E-34}

	exemplars["vector_empty.json"] = []interface{}{}

	exemplars["vector_simple.json"] = []int64{1, 2, 3}

	mixed := []interface{}{0, 1, 2.0, true, false,
		"five", Keyword("six"), Symbol("seven"), "~eight", nil}

	exemplars["vector_nested.json"] = []interface{}{
		[]int64{1, 2, 3},
		mixed}

	exemplars["map_numeric_keys.json"] = map[int]string{1: "one", 2: "two"}

	exemplars["vector_1935_keywords_repeated_twice.json"] = makeBigArray(1935)
	exemplars["vector_1936_keywords_repeated_twice.json"] = makeBigArray(1936)
	exemplars["vector_1937_keywords_repeated_twice.json"] = makeBigArray(1937)

	exemplars["set_nested.json"] =
		MakeSet(MakeSet(1, 3, 2), MakeSet(nil, 0, 2.0, "~eight", 1, true, "five", false, Symbol("seven"), Keyword("six")))
	// */

	/*
		exemplars["uris.json"] = []*url.URL{
			toUrl("http://example.com"),
			toUrl("ftp://example.com"),
			toUrl("file:///path/to/file.txt"),
			toUrl("http://www.詹姆斯.com/"),
		}

	*/

	exemplars["maps_four_char_keyword_keys.json"] = []interface{}{
		map[interface{}]int{Keyword("bbbb"): 2, Keyword("aaaa"): 1},
		map[Keyword]int{Keyword("bbbb"): 4, Keyword("aaaa"): 3},
		map[interface{}]int{Keyword("bbbb"): 6, Keyword("aaaa"): 5},
	}
}

func makeBigArray(size int) []Keyword {
	result := make([]Keyword, 2*size)

	for i := 0; i < 2*size; i++ {
		j := i
		if j >= size {
			j = j - size
		}
		sKey := fmt.Sprintf("key%04d", j)
		result[i] = Keyword(sKey)
	}

	return result
}

func makeBigMap(size int) *map[Keyword]int {
	result := map[Keyword]int{}

	for i := 0; i < size; i++ {
		sKey := fmt.Sprintf("key%04d", i)
		result[Keyword(sKey)] = i
	}

	return &result
}

func makeBigNestedMap(size int) *map[Keyword]interface{} {
	f := makeBigMap(size)
	s := makeBigMap(size)

	return &map[Keyword]interface{}{Keyword("f"): *f, Keyword("s"): *s}
}
func TestValues(t *testing.T) {
	for exemplar, value := range exemplars {
		Verify(t, value, ExemplarPath(exemplar))
	}
}

func XTestCaching(t *testing.T) {
	value := []interface{}{
		Symbol("abcdefg"), Symbol("abcdefg"), Symbol("abc"),
		//Keyword("abcdefg"), Keyword("abcdefg"), Keyword("abc"),
	}

	s, _ := EncodeToString(value)
	log.Println("=======> ", s)
	newV, _ := DecodeFromString(s)
	log.Println("======> ", newV)
}

