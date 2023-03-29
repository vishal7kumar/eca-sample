package cmd

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"os"
	"path/filepath"
)

type DecoderService interface {
	Decode(fileName string, filePaths []string) string
}

type simpleDecoder struct {
	enc          reedsolomon.Encoder
	dataShards   int
	parityShards int
}

func NewDecoder(dataShards int, parityShards int) DecoderService {
	enc, err := reedsolomon.New(dataShards, parityShards)
	checkErr(err)
	return &simpleDecoder{enc: enc, dataShards: dataShards, parityShards: parityShards}
}

func (e simpleDecoder) Decode(fileName string, filePaths []string) string {
	totalPaths := len(filePaths)

	shards := make([][]byte, e.dataShards+e.parityShards)
	for i := range shards {
		// Round Robin allocation of shards
		inputFile := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i%totalPaths], fileName), i)
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

	// TODO: Write data to the PWD ?
	cwd, err := os.Getwd()
	checkErr(err)

	outFilePath := filepath.Join(cwd, fileName)
	fmt.Println("Writing data to", filepath.Join(cwd, fileName))
	f, err := os.Create(outFilePath)
	checkErr(err)

	// We don't know the exact filesize.
	err = e.enc.Join(f, shards, len(shards[0])*e.dataShards)
	checkErr(err)

	return outFilePath
}
