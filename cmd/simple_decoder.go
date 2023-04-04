package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

type DecoderService interface {
	Decode(fileName string, filePaths []string) (io.Reader, int)
}

type simpleDecoder struct {
	enc          reedsolomon.Encoder
	dataShards   int
	parityShards int
}

func NewSimpleDecoder(dataShards int, parityShards int) DecoderService {
	enc, err := reedsolomon.New(dataShards, parityShards)
	checkErr(err)
	return &simpleDecoder{enc: enc, dataShards: dataShards, parityShards: parityShards}
}

func (e simpleDecoder) Decode(fileName string, filePaths []string) (io.Reader, int) {

	shards := make([][]byte, e.dataShards+e.parityShards)
	for i := range shards {
		inputFile := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i], fileName), i)
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
	checkErr(err)
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

	reader, writer := io.Pipe()

	// We don't know the exact filesize.
	fileSize := len(shards[0]) * e.dataShards

	// Join the shards and write them
	go func() {
		defer writer.Close()
		err = e.enc.Join(writer, shards, fileSize)
		checkErr(err)
		// TODO: the err needs to be used in the main program
		// need to use channels
	}()

	return reader, fileSize
}
