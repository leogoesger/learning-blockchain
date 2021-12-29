package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"` // new transactions only (payload)
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

type BlockFS struct {
	Hash  Hash   `json:"hash"`
	Block *Block `json:"block"`
}

func (b *Block) Hash() (Hash, error) {
	bytes, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(bytes), nil
}

func NewBlock(lastestHash Hash, timeNow uint64, txs []Tx) *Block {
	return &Block{
		Header: BlockHeader{
			Parent: lastestHash,
			Time:   timeNow,
		},
		TXs: txs,
	}
}
