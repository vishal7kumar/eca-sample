package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Reader interface {
	Read(fileName string, filePaths []string) (io.Reader, int64)
}

type objectReader struct {
	decoder DecoderService
	cfg     *viper.Viper
}

func (o objectReader) Read(fileName string, filePaths []string) (io.Reader, int64) {
	reader, size := o.decoder.Decode(fileName, filePaths)
	return reader, int64(size)
}

func NewObjectReader(dataShards int, parityShards int, fileName string, filePaths []string, cfg *viper.Viper) Reader {
	if dataShards == 0 && parityShards == 0 {
		dataShards = cfg.GetInt("dataShards")
		parityShards = cfg.GetInt("parityShards")
	}

	var decoder DecoderService
	shardSize, err := getShardSize(fileName, filePaths, dataShards, parityShards)

	if shardSize > 10*1024*1024 || err != nil {
		decoder = NewStreamingDecoder(dataShards, parityShards)
	} else {
		decoder = NewSimpleDecoder(dataShards, parityShards)
	}

	return &objectReader{decoder: decoder, cfg: cfg}
}

func getShardSize(fileName string, filePaths []string, dataShards int, parityShards int) (int64, error) {
	totalShards := dataShards + parityShards

	for i := 0; i < totalShards; i++ {
		inputFileName := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i], fileName), i)
		fileInfo, err := os.Stat(inputFileName)
		if err != nil {
			log.Printf("Error: %v \n", err)
			log.Println("continuing to the next file shard")
			continue
		}

		return fileInfo.Size(), nil

	}
	return 0, ErrNoValidShardFound
}
