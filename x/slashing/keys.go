package slashing

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stake "github.com/cosmos/cosmos-sdk/x/stake/types"
)

// key prefix bytes
var (
	ValidatorSigningInfoKey         = []byte{0x01} // Prefix for signing info
	ValidatorMissedBlockBitArrayKey = []byte{0x02} // Prefix for missed block bit array
	ValidatorSlashingPeriodKey      = []byte{0x03} // Prefix for slashing period
	AddrPubkeyRelationKey           = []byte{0x04} // Prefix for address-pubkey relation
	SlashRecordKey                  = []byte{0x05} // Prefix for slash record
)

// stored by *Tendermint* address (not operator address)
func GetValidatorSigningInfoKey(consAddr []byte) []byte {
	return append(ValidatorSigningInfoKey, consAddr...)
}

// stored by *Tendermint* address (not operator address)
func GetValidatorMissedBlockBitArrayPrefixKey(v sdk.ConsAddress) []byte {
	return append(ValidatorMissedBlockBitArrayKey, v.Bytes()...)
}

// stored by *Tendermint* address (not operator address)
func GetValidatorMissedBlockBitArrayKey(v sdk.ConsAddress, i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return append(GetValidatorMissedBlockBitArrayPrefixKey(v), b...)
}

// stored by *Tendermint* address (not operator address)
func GetValidatorSlashingPeriodPrefix(v sdk.ConsAddress) []byte {
	return append(ValidatorSlashingPeriodKey, v.Bytes()...)
}

// stored by *Tendermint* address (not operator address) followed by start height
func GetValidatorSlashingPeriodKey(v sdk.ConsAddress, startHeight int64) []byte {
	b := make([]byte, 8)
	// this needs to be height + ValidatorUpdateDelay because the slashing period for genesis validators starts at height -ValidatorUpdateDelay
	binary.BigEndian.PutUint64(b, uint64(startHeight+stake.ValidatorUpdateDelay))
	return append(GetValidatorSlashingPeriodPrefix(v), b...)
}

func getAddrPubkeyRelationKey(address []byte) []byte {
	return append(AddrPubkeyRelationKey, address...)
}

func GetSlashRecordKey(consAddr []byte, infractionType byte, infractionHeight uint64) []byte {
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, infractionHeight)
	return append(GetSlashRecordsByAddrAndTypeIndexKey(consAddr, infractionType), heightBz...)
}

func GetSlashRecordsByAddrAndTypeIndexKey(sideConsAddr []byte, infractionType byte) []byte {
	return append(GetSlashRecordsByAddrIndexKey(sideConsAddr), []byte{infractionType}...)
}

func GetSlashRecordsByAddrIndexKey(sideConsAddr []byte) []byte {
	return append(SlashRecordKey, sideConsAddr...)
}
