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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/shamaton/msgpack"
)

type DataEmitter interface {
	Emit(s string) error
	EmitString(s string, cacheable bool) error
	EmitTag(s string) error
	EmitInt(i int64, asKey bool) error
	EmitFloat(f float64, asKey bool) error
	EmitNil(asKey bool) error
	EmitBool(bool, asKey bool) error
	EmitStartArray(size int64) error
	EmitArraySeparator() error
	EmitEndArray() error

	EmitStartMap(size int64) error
	EmitMapSeparator() error
	EmitKeySeparator() error
	EmitEndMap() error
}

type JsonEmitter struct {
	writer io.Writer
	cache  Cache
}

func NewJsonEmitter(w io.Writer, cache Cache) *JsonEmitter {
	return &JsonEmitter{writer: w, cache: cache}
}

// Emit the string unaltered and without quotes. This is the lowest level emitter.

func (je JsonEmitter) Emit(s string) error {
	_, err := je.writer.Write([]byte(s))
	return err
}

// EmitBase emits the basic value supplied, encoding it as JSON.

func (je JsonEmitter) EmitBase(x interface{}) error {
	bytes, err := json.Marshal(x)
	if err == nil {
		_, err = je.writer.Write(bytes)
	}
	return err
}

// EmitsTag emits a transit #tag. The string supplied should not include the '#'.

func (je JsonEmitter) EmitTag(s string) error {
	return je.EmitString("~#"+s, true)
}

func (je JsonEmitter) EmitString(s string, cacheable bool) error {
	if je.cache.IsCacheable(s, cacheable) {
		s = je.cache.Write(s)
	}
	return je.EmitBase(s)
}

const MaxJsonInt = 1<<53 - 1

func (je JsonEmitter) EmitInt(i int64, asKey bool) error {
	if asKey || (i > MaxJsonInt) {
		return je.EmitString(fmt.Sprintf("~i%d", i), asKey)
	}
	return je.EmitBase(i)
}

func (je JsonEmitter) EmitNil(asKey bool) error {
	if asKey {
		return je.EmitString("~_", false)
	} else {
		return je.EmitBase(nil)
	}
}

func (je JsonEmitter) EmitFloat(f float64, asKey bool) error {
	if asKey {
		return je.EmitString(fmt.Sprintf("~d%g", f), asKey)
	} else {
		s := fmt.Sprintf("%g", f)
		if !strings.ContainsAny(s, ".eE") {
			s = s + ".0" // Horible hack!
		}
		return je.Emit(s)
	}
}

func (je JsonEmitter) EmitStartArray(size int64) error {
	return je.Emit("[")
}

func (je JsonEmitter) EmitEndArray() error {
	return je.Emit("]")
}

func (je JsonEmitter) EmitArraySeparator() error {
	return je.Emit(",")
}

func (je JsonEmitter) EmitStartMap(size int64) error {
	return je.Emit("{")
}

func (je JsonEmitter) EmitEndMap() error {
	return je.Emit("}")
}

func (je JsonEmitter) EmitMapSeparator() error {
	return je.Emit(",")
}

func (je JsonEmitter) EmitKeySeparator() error {
	return je.Emit(":")
}

func (je JsonEmitter) EmitBool(x bool, asKey bool) error {
	if asKey {
		if x {
			return je.EmitString("~?t", false)
		} else {
			return je.EmitString("~?f", false)
		}
	} else {
		return je.EmitBase(x)
	}
}

type MsgPackEmitter struct {
	writer io.Writer
	cache  Cache
}

func (m MsgPackEmitter) Emit(s string) error {
	_, err := m.writer.Write([]byte(s))
	return err
}

func (m MsgPackEmitter) EmitString(s string, cacheable bool) error {
	if m.cache.IsCacheable(s, cacheable) {
		s = m.cache.Write(s)
	}

	return m.EmitBase(s)
}

func (m MsgPackEmitter) EmitTag(s string) error {
	return m.EmitString(fmt.Sprintf("~#%s", s), true)
}

func (m MsgPackEmitter) EmitInt(i int64, asKey bool) error {
	return m.EmitBase(i)
}

func (m MsgPackEmitter) EmitFloat(f float64, asKey bool) error {
	return m.EmitBase(f)
}

func (m MsgPackEmitter) EmitNil(asKey bool) error {
	return m.EmitBase([]byte(nil))
}

func (m MsgPackEmitter) EmitBool(bool, asKey bool) error {
	return m.EmitBase(bool)
}

func (m MsgPackEmitter) EmitStartArray(size int64) error {
	if size < 0 {
		panic("size is 0")
	} else if size < 16 {
		return m.EmitByte(byte(0x90 | size))
	} else if size < 65536 {
		err := m.EmitByte(byte(0xdc))
		if err != nil {
			return err
		}

		return m.EmitShort(uint16(size))
	} else {
		err := m.EmitByte(byte(0xdd))
		if err != nil {
			return err
		}

		return m.EmitShort(uint16(size))
	}
}

func (m MsgPackEmitter) EmitArraySeparator() error {
	return nil
}

func (m MsgPackEmitter) EmitEndArray() error {
	return nil
}

func (m MsgPackEmitter) EmitStartMap(size int64) error {
	if size < 0 {
		panic("size is 0")
	} else if size < 16 {
		return m.EmitByte(byte(0x80 | size))
	} else if size < 65536 {
		err := m.EmitByte(byte(0xde))
		if err != nil {
			return err
		}

		return m.EmitShort(uint16(size))
	} else {
		err := m.EmitByte(byte(0xdf))
		if err != nil {
			return err
		}

		return m.EmitShort(uint16(size))
	}
}

func (m MsgPackEmitter) EmitMapSeparator() error {
	return nil
}

func (m MsgPackEmitter) EmitKeySeparator() error {
	panic("implement me")
}

func (m MsgPackEmitter) EmitEndMap() error {
	return nil
}

func (m MsgPackEmitter) EmitByte(v byte) error {
	_, err := m.writer.Write([]byte{v})
	return err
}

func (m MsgPackEmitter) EmitShort(v uint16) error {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], v)
	_, err := m.writer.Write(b[:])
	return err
}

func (m MsgPackEmitter) EmitBase(v interface{}) error {
	b, err := msgpack.Encode(v)
	if err != nil {
		return err
	}

	_, err = m.writer.Write(b)
	return err
}

func NewMsgPackEmitter(w io.Writer, cache Cache) DataEmitter {
	return &MsgPackEmitter{writer: w, cache: cache}
}
