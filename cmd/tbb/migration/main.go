package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leogoesger/learning-blockchain/database"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		log.Fatal(err)
	}

	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			*database.NewTx("andrej", "andrej", 3, ""),
			*database.NewTx("andrej", "andrej", 700, "reward"),
		},
	)

	if err = state.AddBlock(block0); err != nil {
		log.Fatal(err)
	}

	block0Hash, err := state.Persist()
	if err != nil {
		log.Fatal(err)
	}

	block1 := database.NewBlock(
		block0Hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			*database.NewTx("andrej", "babayaga", 2000, ""),
			*database.NewTx("andrej", "andrej", 100, "reward"),
			*database.NewTx("babayaga", "andrej", 1, ""),
			*database.NewTx("babayaga", "caesar", 1000, ""),
			*database.NewTx("babayaga", "andrej", 50, ""),
			*database.NewTx("andrej", "andrej", 600, "reward"),
		})
	state.AddBlock(block1)
	state.Persist()

	fmt.Println(block0Hash)
}
