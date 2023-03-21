package cmd

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"os"
)

type DecoderService interface {
	Decode(filePath string)
}

type decoderService struct {
	enc          reedsolomon.Encoder
	dataShards   int
	parityShards int
}

func NewDecoder(dataShards int, parityShards int) DecoderService {
	enc, err := reedsolomon.New(dataShards, parityShards)
	checkErr(err)
	return &decoderService{enc: enc, dataShards: dataShards, parityShards: parityShards}
}

func (e decoderService) Decode(filePath string) {
	shards := make([][]byte, e.dataShards+e.parityShards)
	for i := range shards {
		inputFile := fmt.Sprintf("%s.%d", filePath, i)
		fmt.Println("Opening", inputFile)
		var err error
		shards[i], err = os.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
		}
		// TODO: consolidated error handling to be done at the higher level
	}

	// Verify the shards
	ok, err := e.enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		err = e.enc.Reconstruct(shards)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		ok, err = e.enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			os.Exit(1)
		}
		checkErr(err)
	}

	// Join the shards and write them

	fmt.Println("Writing data to", filePath)
	f, err := os.Create(filePath)
	checkErr(err)

	// We don't know the exact filesize.
	err = e.enc.Join(f, shards, len(shards[0])*e.dataShards)
	checkErr(err)
}
