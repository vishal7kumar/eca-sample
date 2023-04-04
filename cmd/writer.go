package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Writer interface {
	Write(input io.Reader, fileName string, filePath []string)
}

type objectWriter struct {
	encoder EncoderService
	cfg     *viper.Viper
}

func (o objectWriter) Write(input io.Reader, fileName string, filePath []string) {
	// inputBytes := []byte(input)
	o.encoder.Encode(input, fileName, filePath)
}

func NewObjectWriter(dataShards int, parityShards int, size int64, cfg *viper.Viper) Writer {
	if dataShards == 0 || parityShards == 0 {
		dataShards = cfg.GetInt("dataShards")
		parityShards = cfg.GetInt("parityShards")
	}

	// calculate shard size to determine which Encoder to initialize
	//shards := dataShards + parityShards
	var encoder EncoderService
	var err error
	if size > (10 * 1024 * 1024) {
		encoder, err = NewStreamingEncoder(dataShards, parityShards, size)
		if err != nil {
			log.Fatal("unable to initialize encoder.")
			os.Exit(1)
		}
	} else {
		encoder, err = NewSimpleEncoder(dataShards, parityShards, size)
		if err != nil {
			log.Fatal("unable to initialize encoder.")
			os.Exit(1)
		}
	}

	if err == ErrInvalidDataAndParitySum {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return &objectWriter{encoder: encoder, cfg: cfg}
}
