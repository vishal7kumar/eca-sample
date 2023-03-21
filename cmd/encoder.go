package cmd

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"os"
	"path/filepath"
)

type EncoderService interface {
	Encode(input []byte, filePath string)
}

type encoderService struct {
	enc reedsolomon.Encoder
}

func NewEncoder(dataShards int, parityShards int) (EncoderService, error) {
	if dataShards+parityShards > 256 {
		return &encoderService{}, invalidDataAndParitySumErr
	}

	enc, err := reedsolomon.New(dataShards, parityShards)
	checkErr(err)

	return &encoderService{enc: enc}, nil
}

func (e encoderService) Encode(data []byte, filePath string) {
	shards, err := e.enc.Split(data)
	checkErr(err)

	// Encode parity
	err = e.enc.Encode(shards)
	checkErr(err)

	dir, file := filepath.Split(filePath)

	for i, shard := range shards {
		outFile := fmt.Sprintf("%s.%d", file, i)

		fmt.Println("Writing to", outFile)
		err = os.WriteFile(filepath.Join(dir, outFile), shard, 0644)
		checkErr(err)
	}
}
