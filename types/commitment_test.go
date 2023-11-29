// Copyright (c) 2022 The illium developers
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package types

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnlockingScript_Hash(t *testing.T) {
	sh, err := hex.DecodeString("13e0143cceae5e7e44d8025c57f4759cfb6384e4a2d3d1106e6c098603845900")
	assert.NoError(t, err)
	param2, err := hex.DecodeString("0cef7dd85c04c505d55c063824a5bad62170db0d37e2068fc6c749ada2cb8293")
	assert.NoError(t, err)
	param1 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	ul := UnlockingScript{
		ScriptCommitment: sh,
		ScriptParams:     [][]byte{param1, param2},
	}

	expr, err := ul.lurkExpression()
	assert.NoError(t, err)
	assert.Equal(t, "(cons 0x13e0143cceae5e7e44d8025c57f4759cfb6384e4a2d3d1106e6c098603845900 (cons 1 (cons 0x0cef7dd85c04c505d55c063824a5bad62170db0d37e2068fc6c749ada2cb8293 nil)))", expr)
	start := time.Now()
	h, err := ul.Hash()
	fmt.Println(time.Since(start))
	assert.NoError(t, err)
	assert.Equal(t, "0e259200938dd2eb040d998ebbbbac8c14dc631125d8105cd996d2f1d0d24301", h.String())
}

func TestSpendNote_Deserialize(t *testing.T) {
	sh, err := hex.DecodeString("13e0143cceae5e7e44d8025c57f4759cfb6384e4a2d3d1106e6c098603845900")
	assert.NoError(t, err)
	param2, err := hex.DecodeString("0cef7dd85c04c505d55c063824a5bad62170db0d37e2068fc6c749ada2cb8293")
	assert.NoError(t, err)
	param1 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	ul := UnlockingScript{
		ScriptCommitment: sh,
		ScriptParams:     [][]byte{param1, param2},
	}

	ser := ul.Serialize()
	assert.Len(t, ser, 32+33+9)

	ul2 := new(UnlockingScript)
	err = ul2.Deserialize(ser)
	assert.NoError(t, err)

	assert.Equal(t, ul.ScriptCommitment, ul2.ScriptCommitment)
	assert.Equal(t, ul.ScriptParams, ul2.ScriptParams)
}
