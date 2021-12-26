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

type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile *os.File
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

func (s *State) Add(tx *Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, *tx)
	return nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) Persist() error {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for idx := range mempool {
		bytes, err := json.Marshal(mempool[idx])
		if err != nil {
			return err
		}
		if _, err := s.dbFile.Write(append(bytes, '\n')); err != nil {
			return err
		}
		s.txMempool = s.txMempool[1:]
	}

	return nil
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

	txFilePath := filepath.Join(cwd, "database", "tx.db")
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
		var tx Tx
		if err = json.Unmarshal(scanner.Bytes(), &tx); err != nil {
			log.Fatal(err)
		}
		if err = state.apply(&tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}
