package common

import (
	"bytes"
	"encoding/json"
	gconf "github.com/DSiSc/craft/config"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto/sha3"
	"math/big"
)

type NodeType int

const (
	UnknownNode   NodeType = iota // UnknownNode --> 0, Unknown node type
	ConsensusNode                 // ConsensusNode --> 1, Consensus function node
	FullNode                       // FullNode --> 2, Full node
	LightNode                      // LightNode --> 3, Light node with simply function
	MaxNodeType                    // MaxNodeType is the boundary of node type
)

const (
	InvalidInt        int    = 255
	BlankString       string = ""
	FBFTRoundInterval        = 500
)

// Sum returns the first 32 bytes of hash of the bz.
func Sum(bz []byte) []byte {
	var alg string
	if value, ok := gconf.GlobalConfig.Load(gconf.HashAlgName); ok {
		alg = value.(string)
	} else {
		alg = "SHA256"
	}
	hasher := sha3.NewHashByAlgName(alg)
	hasher.Write(bz)
	hash := hasher.Sum(nil)
	return hash[:types.HashLength]
}

func TxHash(tx *types.Transaction) (hash types.Hash) {
	jsonByte, _ := json.Marshal(tx)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

func HeaderHash(block *types.Block) (hash types.Hash) {
	//var defaultHash types.Hash
	if !(block.HeaderHash == types.Hash{}) {
		copy(hash[:], block.HeaderHash[:])
		return
	}
	jsonByte, _ := json.Marshal(block.Header)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

func CopyBytes(b []byte) (copiedBytes []byte) {
	var temp []byte
	if bytes.Equal(b, temp) {
		log.Error("src byte is nil, please confirm.")
		return
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)
	return
}

// New a transaction
func newTransaction(nonce uint64, to *types.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from *types.Address) *types.Transaction {
	if len(data) > 0 {
		data = CopyBytes(data)
	}
	d := types.TxData{
		AccountNonce: nonce,
		Recipient:    to,
		From:         from,
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gasLimit,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &types.Transaction{Data: d}
}

func NewTransaction(nonce uint64, to types.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from types.Address) *types.Transaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data, &from)
}

type MsgType uint8

const (
	MsgNull MsgType = iota
	MsgBlockCommitSuccess
	MsgBlockCommitFailed
	MsgBlockVerifyFailed
	MsgNodeServiceStopped
	MsgRoundRunFailed
	MsgToConsensusFailed
	MsgChangeMaster
)

type SysConfig struct {
	LogLevel log.Level
	LogPath  string
	LogStyle string
}
