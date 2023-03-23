package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Writer interface {
	Write(input string, filePath string)
}

type objectWriter struct {
	encoder EncoderService
	cfg     *viper.Viper
}

func (o objectWriter) Write(input string, filePath string) {
	inputBytes := []byte(input)
	o.encoder.Encode(inputBytes, filePath)
}

func NewObjectWriter(dataShards int, parityShards int, cfg *viper.Viper) Writer {
	if dataShards == 0 || parityShards == 0 {
		dataShards = cfg.GetInt("dataShards")
		parityShards = cfg.GetInt("parityShards")
	}
	encoder, err := NewEncoder(dataShards, parityShards)
	if err == invalidDataAndParitySumErr {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return &objectWriter{encoder: encoder, cfg: cfg}
}
