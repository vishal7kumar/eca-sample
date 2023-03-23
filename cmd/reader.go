package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Reader interface {
	Read(filePath string) []byte
}

type objectReader struct {
	decoder DecoderService
	cfg     *viper.Viper
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

func NewObjectReader(dataShards int, parityShards int, cfg *viper.Viper) Reader {
	if dataShards == 0 && parityShards == 0 {
		dataShards = cfg.GetInt("dataShards")
		parityShards = cfg.GetInt("parityShards")
	}
	return &objectReader{decoder: NewDecoder(dataShards, parityShards), cfg: cfg}
}
