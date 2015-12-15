package testHelper

import (
	"github.com/FactomProject/factomd/common/directoryBlock/dbInfo"
)

func CreateTestDirBlockInfo(prev *dbInfo.DirBlockInfo) *dbInfo.DirBlockInfo {
	dbi := dbInfo.NewDirBlockInfo()
	if prev == nil {
		dbi.DBHeight = 0
	} else {
		dbi.DBHeight = prev.DBHeight + 1
	}
	height := dbi.DBHeight

	dbi.DBHash.UnmarshalBinary(intToByteSlice(int(height)))
	dbi.Timestamp = int64(height)
	dbi.BTCTxHash.UnmarshalBinary(intToByteSlice(int(height)))
	dbi.BTCTxOffset = int32(int(height))
	dbi.BTCBlockHeight = int32(height)
	dbi.BTCBlockHash.UnmarshalBinary(intToByteSlice(255 - int(height)))
	dbi.DBMerkleRoot.UnmarshalBinary(intToByteSlice(255 - int(height)))
	dbi.BTCConfirmed = height%2 == 0

	return dbi
}

func intToByteSlice(n int) []byte {
	answer := make([]byte, 32)
	for i := range answer {
		answer[i] = byte(n)
	}
	return answer
}
