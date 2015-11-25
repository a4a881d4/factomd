// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid_test

import (
	"fmt"
	"github.com/FactomProject/ed25519"
	. "github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"math/rand"
	"testing"
)

var _ = fmt.Printf
var _ = ed25519.Sign
var _ = rand.New

var sig1 [64]byte
var sig2 [64]byte

var s1, s2 interfaces.ISignature

func Test_Setup_Signature(test *testing.T) {
	sh11 := primitives.Sha([]byte("sig first half  one")).Bytes()
	sh12 := primitives.Sha([]byte("sig second half one")).Bytes()
	sh21 := primitives.Sha([]byte("sig first half  two")).Bytes()
	sh22 := primitives.Sha([]byte("sig second half two")).Bytes()

	copy(sig1[:32], sh11)
	copy(sig1[32:], sh12)
	copy(sig2[:32], sh21)
	copy(sig2[32:], sh22)

	s1 = new(FactoidSignature)
	s1.SetSignature(sig1[:])
	s2 = new(FactoidSignature)
	s2.SetSignature(sig2[:])

	//    txt1,_:=s1.CustomMarshalText()
	//    txt2,_:=s2.CustomMarshalText()
	//    Prtln(string(txt1))
	//    Prtln(string(txt2))
}

/*
func Test_IsEqual_Signature(test *testing.T) {

	if s1.IsEqual(s2) == nil {
		primitives.PrtStk()
		test.Fail()
	}

	s2.SetSignature(sig1[:]) // Set to sig1 for test

	if s1.IsEqual(s2) != nil {
		primitives.PrtStk()
		test.Fail()
	}

	s2.SetSignature(sig2[:]) // Reset it back to Sig2
}

func Test_Marshal_Signature_(test *testing.T) {
	data, err := s1.MarshalBinary()
	s2.UnmarshalBinaryData(data)

	if s1.IsEqual(s2) != nil {
		primitives.PrtStk()
		primitives.Prtln(err)
		test.Fail()
	}

	s2.SetSignature(sig2[:]) // Reset it back to Sig2

}
*/
