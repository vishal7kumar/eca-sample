package cmd

import (
	"fmt"
	"os"
)

type Writer interface {
	Write(input string, filePath string)
}

type objectWriter struct {
	encoder EncoderService
}

func (o objectWriter) Write(input string, filePath string) {
	inputBytes := []byte(input)
	o.encoder.Encode(inputBytes, filePath)
}

func NewObjectWriter(dataShards int, parityShards int) Writer {
	if dataShards == 0 || parityShards == 0 {
		// read from config
	}
	encoder, err := NewEncoder(dataShards, parityShards)
	if err == invalidDataAndParitySumErr {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}

	return &objectWriter{encoder: encoder}
}
