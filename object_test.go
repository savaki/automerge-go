// Copyright 2020 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package automerge

import (
	"fmt"
	"github.com/savaki/automerge/encoding"
	"github.com/tj/assert"
	"io"
	"testing"
)

func TestObject_Insert(t *testing.T) {
	const n = 1e4
	actor := []byte("me")
	node := NewObject(encoding.RawTypeVarInt)
	for i := int64(0); i < n; i++ {
		refActor := actor
		if i == 0 {
			refActor = nil
		}
		op := Op{
			ID:    NewID(i+1, actor),
			Ref:   NewID(i, refActor),
			Type:  0,
			Value: encoding.RuneValue('a'),
		}
		err := node.Insert(op)
		assert.Nil(t, err)
	}

	fmt.Println(node.Size())
}

func TestObject_NextValue(t *testing.T) {
	var (
		actor = []byte("me")
		want  = "hello world"
		obj   = NewObject(encoding.RawTypeVarInt)
	)

	for i, r := range want {
		refCounter := int64(i)
		refActor := actor
		if i == 0 {
			refActor = nil
		}

		op := Op{
			ID:    NewID(refCounter+1, actor),
			Ref:   NewID(refCounter, refActor),
			Type:  0,
			Value: encoding.RuneValue(r),
		}
		err := obj.Insert(op)
		assert.Nil(t, err)
	}

	got := string(readAllRunes(t, obj))
	assert.Equal(t, want, got)
}

func readAllRunes(t *testing.T, obj *Object) []rune {
	var runes []rune
	var token ValueToken
	var err error
	for {
		token, err = obj.NextValue(token)
		if err == io.EOF {
			break
		}
		assert.Nil(t, err)
		runes = append(runes, rune(token.Value.Int))
	}
	return runes
}
