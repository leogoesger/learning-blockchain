package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Account string
type Snapshot [32]byte

type State struct {
	Balances        map[Account]uint
	txMempool       []Tx
	dbFile          *os.File
	latestBlockHash Hash
}

func (s *State) apply(tx *Tx) error {
	if tx.Data == "reward" {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) AddBlock(b *Block) error {
	for i := range b.TXs {
		s.AddTx(&b.TXs[i])
	}

	return nil
}

func (s *State) applyBlock(b *Block) error {
	for i := range b.TXs {
		s.apply(&b.TXs[i])
	}

	return nil
}

func (s *State) AddTx(tx *Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, *tx)
	return nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) Persist() (Hash, error) {
	block := NewBlock(
		s.latestBlockHash,
		uint64(time.Now().Unix()),
		s.txMempool,
	)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, nil
	}

	s.latestBlockHash = blockHash

	blockFs := BlockFS{blockHash, block}
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, nil
	}

	if _, err := s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash
	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

type Genesis struct {
	GenesisTime time.Time        `json:"gensis_time"`
	ChainID     string           `json:"chain_id"`
	Balances    map[Account]uint `json:"balances"`
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	f, err := os.Open(genFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var genesis Genesis
	if err = json.Unmarshal(bytes, &genesis); err != nil {
		log.Fatal(err)
	}

	txFilePath := filepath.Join(cwd, "database", "block.db")
	file, err := os.OpenFile(txFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}

	state := &State{
		Balances:  genesis.Balances,
		txMempool: nil,
		dbFile:    file,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var blockFS BlockFS
		if err = json.Unmarshal(scanner.Bytes(), &blockFS); err != nil {
			log.Fatal(err)
		}
		if err = state.applyBlock(blockFS.Block); err != nil {
			return nil, err
		}
		state.latestBlockHash = blockFS.Hash
	}

	return state, nil
}
