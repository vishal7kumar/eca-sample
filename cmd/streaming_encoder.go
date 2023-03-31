package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

type streamingEncoder struct {
	enc reedsolomon.StreamEncoder
	dataShards int
	parityShards int
	size int64
}

func NewStreamingEncoder(dataShards int, parityShards int, size int64) (EncoderService, error) {
	if dataShards+parityShards > 256 {
		return &simpleEncoder{}, invalidDataAndParitySumErr
	}

	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	checkErr(err)

	return &streamingEncoder{enc: enc, dataShards: dataShards, parityShards: parityShards, size: size}, nil
}

func (e streamingEncoder) Encode(data io.Reader, fileName string, filePaths []string) {
	shards := e.dataShards + e.parityShards
	out := make([]*os.File, shards)

	for i := range out {
		outputFile := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i], fileName), i)
		fmt.Println("Creating", outputFile)
		var err error
		out[i], err = os.Create(outputFile)
		checkErr(err)
	}

	// Split into files.
	dataWriter := make([]io.Writer, e.dataShards)
	for i := range dataWriter {
		dataWriter[i] = out[i]
	}

	// Do the split
	err := e.enc.Split(data, dataWriter, e.size)
	checkErr(err)

	// Close and re-open the files.
	input := make([]io.Reader, e.dataShards)

	for i := range dataWriter {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		checkErr(err)
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, e.parityShards)
	for i := range parity {
		parity[i] = out[e.dataShards+i]
		defer out[e.dataShards+i].Close()
	}

	// Encode parity
	err = e.enc.Encode(input, parity)
	checkErr(err)
	fmt.Printf("File split into %d data + %d parity shards.\n", e.dataShards, e.parityShards)
}
