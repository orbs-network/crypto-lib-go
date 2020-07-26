// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package encoding

import (
	"encoding/hex"
	"github.com/orbs-network/crypto-lib-go/crypto/hash"
	"github.com/pkg/errors"
	"strings"
)

func EncodeHex(data []byte) string { // EIP-55 complaint
	result := []byte(hex.EncodeToString(data)) // hex does all lowercase
	hashed := hash.CalcKeccak256(result)
	hashedHex := hex.EncodeToString(hashed)
	hashedHexLen := len(hashedHex)

	for i := 0; i < len(result); i++ {
		if result[i] > '9' && hashedHex[i%hashedHexLen] > '7' { // we rely on 'a' > 'A' > '9'
			result[i] -= 32 // turn lower case to upper case in ascii
		}
	}

	return "0x" + string(result)
}

// on decode error (eg. non hex character in str) returns zero_value, error
// on checksum failure returns decoded_value, error (so users could warn about checksum but still use the decoded)
// if all is lower or upper then the checksum check is ignored (as the checksum was probably not taken into account)
func DecodeHex(str string) ([]byte, error) {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}

	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, errors.Wrap(err, "invalid hex string")
	}

	encoded := EncodeHex(data)
	if encoded[2:] != str {
		// checksum error, we will allow if the source is in uniform case (all lower/upper)
		if strings.ToUpper(str) == str || strings.ToLower(str) == str {
			return data, nil
		} else {
			return data, errors.New("invalid checksum")
		}
	}

	return data, nil
}
