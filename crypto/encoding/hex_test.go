// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package encoding

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

type testStringPair struct {
	sourceHex          string
	checksumEncodedHex string
}

var encodeStringTestTable = []testStringPair{
	{"5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed", "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"},
	{"19ef9290b8cf5ec5e72f9fde3e044b37736ec0c7", "0x19ef9290B8cf5EC5e72F9fDE3E044b37736Ec0C7"},
	{"dbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB", "0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB"},
	{"D1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", "0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb"},
	{"5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed19ef9290b8cf5ec5e72f9fde3e044b37736ec0c7", "0x5AaEB6053f3E94C9b9A09f33669435e7Ef1BeaEd19eF9290B8Cf5EC5E72F9Fde3E044b37736EC0C7"},
}

func TestHexEncodeWithChecksum(t *testing.T) {
	for _, pair := range encodeStringTestTable {
		data, err := hex.DecodeString(pair.sourceHex)
		require.NoError(t, err, "failed to decode, human error most likely")
		encoded := EncodeHex(data)

		require.Equal(t, pair.checksumEncodedHex, encoded, "expected encoding with a specific result for each input")
	}
}

func TestHexDecodeGoodChecksum(t *testing.T) {
	for _, pair := range encodeStringTestTable {
		rawData, err := hex.DecodeString(pair.sourceHex)
		require.NoError(t, err, "failed to decode, human error most likely")
		decoded, err := DecodeHex(pair.checksumEncodedHex)
		require.NoError(t, err, "checksum should be valid")
		require.Equal(t, rawData, decoded, "data should be decoded correctly")
	}
}

func TestHexDecodeBadChecksum(t *testing.T) {
	sourceHex := "D1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb"
	rawData, err := hex.DecodeString(sourceHex)
	require.NoError(t, err, "failed to decode, human error most likely")
	wrongCheckSum := "0xd1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb"
	decoded, err := DecodeHex(wrongCheckSum)
	require.EqualError(t, err, "invalid checksum", "checksum should be invalid")
	require.Equal(t, rawData, decoded, "data should be decoded correctly even though checksum is invalid")
}

func TestHexDecodeInvalidHex(t *testing.T) {
	decoded, err := DecodeHex("0")
	require.Error(t, err, "should not succeed with invalid hex")
	require.Nil(t, decoded, "result should be nil")
}

func BenchmarkHexEncodeWithChecksum(b *testing.B) {
	rawData, err := hex.DecodeString(encodeStringTestTable[0].sourceHex)
	require.NoError(b, err, "failed to decode, human error most likely")
	for i := 0; i < b.N; i++ {
		EncodeHex(rawData)
	}
}

func BenchmarkHexDecodeWithChecksum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := DecodeHex(encodeStringTestTable[0].checksumEncodedHex)
		if err != nil { // require/testify is very slow
			b.Fail()
		}
	}
}
