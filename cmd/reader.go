package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Reader interface {
	Read(fileName string, filePaths []string) []byte
}

type objectReader struct {
	decoder DecoderService
	cfg     *viper.Viper
}

func (o objectReader) Read(fileName string, filePaths []string) []byte {
	outFilePath := o.decoder.Decode(fileName, filePaths)
	fileContents, err := readFile(outFilePath)
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

func NewObjectReader(dataShards int, parityShards int, cfg *viper.Viper) Reader {
	if dataShards == 0 && parityShards == 0 {
		dataShards = cfg.GetInt("dataShards")
		parityShards = cfg.GetInt("parityShards")
	}

	// TODO: based on the size of a shard, return suitable decoder
	// return &objectReader{decoder: NewStreamingDecoder(dataShards, parityShards), cfg: cfg}

	return &objectReader{decoder: NewSimpleDecoder(dataShards, parityShards), cfg: cfg}
}
