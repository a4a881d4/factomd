// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package primitives_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/FactomProject/factomd/common/constants"
	. "github.com/FactomProject/factomd/common/primitives"
)

func TestHashIsEqual(t *testing.T) {
	// A hash
	var hash = [constants.ADDRESS_LENGTH]byte{
		0x61, 0xe3, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72, 0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
		0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c, 0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	h1 := new(Hash)
	h2 := new(Hash)

	if h1.IsEqual(h2) != nil { // Out of the box, hashes should be equal
		t.Errorf("Hashes are not equal")
	}

	h1.SetBytes(hash[:])

	if h1.IsEqual(h2) == nil { // Now they should not be equal
		t.Errorf("Hashes are equal")
	}

	h2.SetBytes(hash[:])

	if h1.IsEqual(h2) != nil { // Back to equality!
		t.Errorf("Hashes are not equal")
	}

	hash2 := h1.Fixed()
	for i := range hash {
		if hash[i] != hash2[i] {
			t.Errorf("Hashes are not equal")
		}
	}
}

//Test vectors: http://www.di-mgt.com.au/sha_testvectors.html

func TestHash(t *testing.T) {
	h := new(Hash)
	err := h.SetBytes(constants.EC_CHAINID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	bytes1, err := h.MarshalBinary()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("bytes1: %v\n", bytes1)

	h2 := new(Hash)
	err = h2.UnmarshalBinary(bytes1)
	t.Logf("h2.bytes: %v\n", h2.Bytes)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	bytes2, err := h2.MarshalBinary()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("bytes2: %v\n", bytes2)

	if bytes.Compare(bytes1, bytes2) != 0 {
		t.Errorf("Invalid output")
	}

	if h2.GetHash() != nil {
		t.Errorf("Hash GetHashed returned something other than nil")
	}
}

func TestSha(t *testing.T) {
	testVector := map[string]string{}
	testVector["abc"] = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	testVector[""] = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	testVector["abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"] = "248d6a61d20638b8e5c026930c3e6039a33ce45964ff2167f6ecedd419db06c1"
	testVector["abcdefghbcdefghicdefghijdefghijkefghijklfghijklmghijklmnhijklmnoijklmnopjklmnopqklmnopqrlmnopqrsmnopqrstnopqrstu"] = "cf5b16a778af8380036ce59e7b0492370b249b11e8f07a51afac45037afee9d1"

	for k, v := range testVector {
		answer, err := DecodeBinary(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		hash := Sha([]byte(k))

		if bytes.Compare(hash.Bytes(), answer) != 0 {
			t.Errorf("Wrong SHA hash for %v", k)
		}
		if hash.String() != v {
			t.Errorf("Wrong SHA hash string for %v", k)
		}
	}
}

func TestSha512Half(t *testing.T) {
	testVector := map[string]string{}
	testVector["abc"] = "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	testVector[""] = "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce"
	testVector["abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"] = "204a8fc6dda82f0a0ced7beb8e08a41657c16ef468b228a8279be331a703c335"
	testVector["abcdefghbcdefghicdefghijdefghijkefghijklfghijklmghijklmnhijklmnoijklmnopjklmnopqklmnopqrlmnopqrsmnopqrstnopqrstu"] = "8e959b75dae313da8cf4f72814fc143f8f7779c6eb9f7fa17299aeadb6889018"

	for k, v := range testVector {
		answer, err := DecodeBinary(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		hash := Sha512Half([]byte(k))

		if bytes.Compare(hash.Bytes(), answer) != 0 {
			t.Errorf("Wrong SHA512Half hash for %v", k)
		}
		if hash.String() != v {
			t.Errorf("Wrong SHA512Half hash string for %v", k)
		}
	}
}

func TestHashStrings(t *testing.T) {
	base := "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}
	if hash.String() != base {
		t.Error("Invalid conversion to string")
	}

	text, err := hash.CustomMarshalText()
	if err != nil {
		t.Error(err)
	}

	if string(text) != base {
		t.Errorf("CustomMarshalText failed - %v vs %v", string(text), base)
	}

	text, err = hash.JSONByte()
	if err != nil {
		t.Error(err)
	}

	if string(text) != fmt.Sprintf("\"%v\"", base) {
		t.Errorf("JSONByte failed - %v vs %v", string(text), base)
	}

	str, err := hash.JSONString()
	if err != nil {
		t.Error(err)
	}

	if str != fmt.Sprintf("\"%v\"", base) {
		t.Errorf("JSONString failed - %v vs %v", string(text), base)
	}

	b := new(bytes.Buffer)
	err = hash.JSONBuffer(b)
	if err != nil {
		t.Error(err)
	}

	if string(b.Bytes()) != fmt.Sprintf("\"%v\"", base) {
		t.Errorf("JSONString failed - %v vs %v", string(text), base)
	}
}

func TestIsSameAs(t *testing.T) {
	base := "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}
	hex, err := DecodeBinary(base)
	if err != nil {
		t.Error(err)
	}
	hash2, err := NewShaHash(hex)
	if err != nil {
		t.Error(err)
	}
	if hash.IsSameAs(hash2) == false {
		t.Error("Identical hashes not recognized as such")
	}

	hash3 := hash.Copy()
	if hash.IsSameAs(hash3) == false {
		t.Errorf("Copied hash is not identical")
	}
}

func TestHashMisc(t *testing.T) {
	base := "4040404040404040404040404040404040404040404040404040404040404040"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}
	if hash.String() != base {
		t.Error("Error in String")
	}

	hash2, err := NewShaHashFromStr(base)
	if err != nil {
		t.Error(err)
	}

	if hash2.ByteString() != "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@" {
		t.Errorf("Error in ByteString - received %v", hash2.ByteString())
	}

	h, err := hex.DecodeString(base)
	if err != nil {
		t.Error(err)
	}
	hash = NewHash(h)
	if hash.String() != base {
		t.Error("Error in NewHash")
	}

	//***********************

	if hash.IsSameAs(nil) != false {
		t.Error("Error in IsSameAs")
	}

	//***********************

	minuteHash, err := HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Error(err)
	}
	if minuteHash.IsMinuteMarker() == false {
		t.Error("Error in IsMinuteMarker")
	}

	hash = NewZeroHash()
	if hash.String() != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Error("Error in NewZeroHash")
	}
}

func TestHashIsZero(t *testing.T) {
	strs := []string{
		"0000000000000000000000000000000000000000000000000000000000000001",
		"0000000000000000000000000000000000000000000000000000000000000002",
		"0000000000000000000000000000000000000000000000000000000000000003",
		"0000000000000000000000000000000000000000000000000000000000000004",
		"0000000000000000000000000000000000000000000000000000000000000005",
		"0000000000000000000000000000000000000000000000000000000000000006",
		"0000000000000000000000000000000000000000000000000000000000000007",
		"0000000000000000000000000000000000000000000000000000000000000008",
		"0000000000000000000000000000000000000000000000000000000000000009",
		"000000000000000000000000000000000000000000000000000000000000000a",
		"000000000000000000000000000000000000000000000000000000000000000b",
		"000000000000000000000000000000000000000000000000000000000000000c",
		"000000000000000000000000000000000000000000000000000000000000000d",
		"000000000000000000000000000000000000000000000000000000000000000e",
		"000000000000000000000000000000000000000000000000000000000000000f"}
	for _, str := range strs {
		h, err := NewShaHashFromStr(str)
		if err != nil {
			t.Error(err)
		}
		if h.IsZero() == true {
			t.Errorf("Non-zero hash is zero")
		}
	}

	h, err := NewShaHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Error(err)
	}
	if h.IsZero() == false {
		t.Errorf("Zero hash is non-zero")
	}

}

func TestIsMinuteMarker(t *testing.T) {
	strs := []string{
		"0000000000000000000000000000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000000000000000000000000001",
		"0000000000000000000000000000000000000000000000000000000000000002",
		"0000000000000000000000000000000000000000000000000000000000000003",
		"0000000000000000000000000000000000000000000000000000000000000004",
		"0000000000000000000000000000000000000000000000000000000000000005",
		"0000000000000000000000000000000000000000000000000000000000000006",
		"0000000000000000000000000000000000000000000000000000000000000007",
		"0000000000000000000000000000000000000000000000000000000000000008",
		"0000000000000000000000000000000000000000000000000000000000000009",
		"000000000000000000000000000000000000000000000000000000000000000a",
		"000000000000000000000000000000000000000000000000000000000000000b",
		"000000000000000000000000000000000000000000000000000000000000000c",
		"000000000000000000000000000000000000000000000000000000000000000d",
		"000000000000000000000000000000000000000000000000000000000000000e",
		"000000000000000000000000000000000000000000000000000000000000000f"}
	for _, str := range strs {
		hash, err := HexToHash(str)
		if err != nil {
			t.Errorf("%v", err)
		}
		if hash.IsMinuteMarker() == false {
			t.Errorf("Entry %v is not a minute marker!", str)
		}
	}
	strs = []string{
		"1000000000000000000000000000000000000000000000000000000000000000",
		"0200000000000000000000000000000000000000000000000000000000000000",
		"0030000000000000000000000000000000000000000000000000000000000000",
		"0004000000000000000000000000000000000000000000000000000000000000",
		"0000500000000000000000000000000000000000000000000000000000000000",
		"0000060000000000000000000000000000000000000000000000000000000000",
		"0000007000000000000000000000000000000000000000000000000000000000",
		"0000000800000000000000000000000000000000000000000000000000000000",
		"0000000090000000000000000000000000000000000000000000000000000000",
		"000000000a000000000000000000000000000000000000000000000000000000",
		"0000000000b00000000000000000000000000000000000000000000000000000",
		"00000000000c0000000000000000000000000000000000000000000000000000",
		"000000000000d000000000000000000000000000000000000000000000000000",
		"0000000000000e00000000000000000000000000000000000000000000000000",
		"00000000000000f0000000000000000000000000000000000000000000000000",
		"0000000000000001000000000000000000000000000000000000000000000000",
		"0000000000000000200000000000000000000000000000000000000000000000",
		"0000000000000000030000000000000000000000000000000000000000000000",
		"0000000000000000004000000000000000000000000000000000000000000000",
		"0000000000000000000500000000000000000000000000000000000000000000",
		"0000000000000000000060000000000000000000000000000000000000000000",
		"0000000000000000000007000000000000000000000000000000000000000000",
		"0000000000000000000000800000000000000000000000000000000000000000",
		"0000000000000000000000090000000000000000000000000000000000000000",
		"000000000000000000000000a000000000000000000000000000000000000000",
		"0000000000000000000000000b00000000000000000000000000000000000000",
		"00000000000000000000000000c0000000000000000000000000000000000000",
		"000000000000000000000000000d000000000000000000000000000000000000",
		"0000000000000000000000000000e00000000000000000000000000000000000",
		"00000000000000000000000000000f0000000000000000000000000000000000",
		"0000000000000000000000000000001000000000000000000000000000000000",
		"0000000000000000000000000000000200000000000000000000000000000000",
		"0000000000000000000000000000000030000000000000000000000000000000",
		"0000000000000000000000000000000004000000000000000000000000000000",
		"0000000000000000000000000000000000500000000000000000000000000000",
		"0000000000000000000000000000000000060000000000000000000000000000",
		"0000000000000000000000000000000000007000000000000000000000000000",
		"0000000000000000000000000000000000000800000000000000000000000000",
		"0000000000000000000000000000000000000090000000000000000000000000",
		"000000000000000000000000000000000000000a000000000000000000000000",
		"0000000000000000000000000000000000000000b00000000000000000000000",
		"00000000000000000000000000000000000000000c0000000000000000000000",
		"000000000000000000000000000000000000000000d000000000000000000000",
		"0000000000000000000000000000000000000000000e00000000000000000000",
		"00000000000000000000000000000000000000000000f0000000000000000000",
		"0000000000000000000000000000000000000000000001000000000000000000",
		"0000000000000000000000000000000000000000000000200000000000000000",
		"0000000000000000000000000000000000000000000000030000000000000000",
		"0000000000000000000000000000000000000000000000004000000000000000",
		"0000000000000000000000000000000000000000000000000500000000000000",
		"0000000000000000000000000000000000000000000000000060000000000000",
		"0000000000000000000000000000000000000000000000000007000000000000",
		"0000000000000000000000000000000000000000000000000000800000000000",
		"0000000000000000000000000000000000000000000000000000090000000000",
		"000000000000000000000000000000000000000000000000000000a000000000",
		"0000000000000000000000000000000000000000000000000000000b00000000",
		"00000000000000000000000000000000000000000000000000000000c0000000",
		"000000000000000000000000000000000000000000000000000000000d000000",
		"0000000000000000000000000000000000000000000000000000000000e00000",
		"00000000000000000000000000000000000000000000000000000000000f0000",
		"0000000000000000000000000000000000000000000000000000000000001000",
		"0000000000000000000000000000000000000000000000000000000000000200"}

	for _, str := range strs {
		hash, err := HexToHash(str)
		if err != nil {
			t.Errorf("%v", err)
		}
		if hash.IsMinuteMarker() == true {
			t.Errorf("Entry %v is a minute marker!", str)
		}

		text, err := hash.MarshalText()
		if err != nil {
			t.Errorf("%v", err)
		}
		if string(text) != str {
			t.Errorf("Invalid marshalled text")
		}
	}
}

func TestStringUnmarshaller(t *testing.T) {
	base := "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"
	hash, err := HexToHash(base)
	if err != nil {
		t.Error(err)
	}

	h2 := new(Hash)
	err = h2.UnmarshalText([]byte(base))
	if err != nil {
		t.Error(err)
	}
	if hash.IsSameAs(h2) == false {
		t.Errorf("Hash from UnmarshalText is incorrect - %v vs %v", hash, h2)
	}

	h3 := new(Hash)
	err = json.Unmarshal([]byte("\""+base+"\""), h3)
	if err != nil {
		t.Error(err)
	}
	if hash.IsSameAs(h3) == false {
		t.Errorf("Hash from json.Unmarshal is incorrect - %v vs %v", hash, h3)
	}
}

func TestDoubleSha(t *testing.T) {
	testVector := map[string]string{}
	testVector["abc"] = "4f8b42c22dd3729b519ba6f68d2da7cc5b2d606d05daed5ad5128cc03e6c6358"
	testVector[""] = "5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456"
	testVector["abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"] = "0cffe17f68954dac3a84fb1458bd5ec99209449749b2b308b7cb55812f9563af"
	testVector["abcdefghbcdefghicdefghijdefghijkefghijklfghijklmghijklmnhijklmnoijklmnopjklmnopqklmnopqrlmnopqrsmnopqrstnopqrstu"] = "accd7bd1cb0fcbd85cf0ba5ba96945127776373a7d47891eb43ed6b1e2ee60fe"

	for k, v := range testVector {
		b := DoubleSha([]byte(k))
		h, err := NewShaHash(b)
		if err != nil {
			t.Error(err)
		}
		if h.String() != v {
			t.Errorf("DoubleSha failed %v != %v", h.String(), v)
		}
	}
}

func TestNewShaHashFromStruct(t *testing.T) {
	testVector := map[string]string{}
	testVector["abc"] = "c127d30fe315d2d3f2dfeae6b9d57c6aa6322c73fb3fd868963660d6cdcd471f"
	testVector[""] = "e2854aa639f07056d58cc02ab52d169c48af8b418fcb0df7842f22a1b2ab3ac2"
	testVector["abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"] = "c226baeb2cad51713659f5e111aaaa6a5a4cfffe7d874c3974c212f4c77fe9d7"
	testVector["abcdefghbcdefghicdefghijdefghijkefghijklfghijklmghijklmnhijklmnoijklmnopjklmnopqklmnopqrlmnopqrsmnopqrstnopqrstu"] = "cdc9eb98889856282bf26c78ffde24c46cbeed70442acf25577fd1aef48a5951"

	for k, v := range testVector {
		h, err := NewShaHashFromStruct(k)
		if err != nil {
			t.Error(err)
		}
		if h.String() != v {
			t.Errorf("NewShaHashFromStruct failed %v != %v", h.String(), v)
		}
	}
}
