package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

type EncoderService interface {
	Encode(input io.Reader, fileName string, filePath []string)
}

type simpleEncoder struct {
	enc  reedsolomon.Encoder
	size int64
}

func NewSimpleEncoder(dataShards int, parityShards int, size int64) (EncoderService, error) {
	if dataShards+parityShards > 256 {
		return &simpleEncoder{}, ErrInvalidDataAndParitySum
	}

	enc, err := reedsolomon.New(dataShards, parityShards)
	checkErr(err)

	return &simpleEncoder{enc: enc, size: size}, nil
}

func (e simpleEncoder) Encode(dataReader io.Reader, fileName string, filePaths []string) {
	data, err := io.ReadAll(dataReader)
	if err != nil {
		log.Fatal("Read error.")
		os.Exit(1)
	}

	shards, err := e.enc.Split(data)
	checkErr(err)

	// Encode parity
	err = e.enc.Encode(shards)
	checkErr(err)

	for i, shard := range shards {
		outFile := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i], fileName), i)

		fmt.Println("Writing to", outFile)
		err = os.WriteFile(outFile, shard, 0644)
		checkErr(err)
	}
}
