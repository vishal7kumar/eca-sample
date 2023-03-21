package cmd

import (
	"fmt"
	"os"
)

type Reader interface {
	Read(filePath string) []byte
}

type objectReader struct {
	decoder DecoderService
}

func (o objectReader) Read(filePath string) []byte {
	o.decoder.Decode(filePath)
	fileContents, err := readFile(filePath)
	if err != nil {
		fmt.Printf("%v", err)
		// ignoring the error for now
	}
	return fileContents
}

func readFile(filePath string) ([]byte, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("error while reading file")
	}
	return fileContents, nil
}

func NewObjectReader(dataShards int, parityShards int) Reader {
	if dataShards == 0 {
		// read from config
	}
	return &objectReader{decoder: NewDecoder(dataShards, parityShards)}
}
