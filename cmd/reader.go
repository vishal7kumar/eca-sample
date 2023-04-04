package cmd

import (
	"io"

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

func NewObjectReader(dataShards int, parityShards int, cfg *viper.Viper) Reader {
	if dataShards == 0 && parityShards == 0 {
		dataShards = cfg.GetInt("dataShards")
		parityShards = cfg.GetInt("parityShards")
	}

	// TODO: based on the size of a shard, return suitable decoder
	// return &objectReader{decoder: NewStreamingDecoder(dataShards, parityShards), cfg: cfg}

	return &objectReader{decoder: NewSimpleDecoder(dataShards, parityShards), cfg: cfg}
}
