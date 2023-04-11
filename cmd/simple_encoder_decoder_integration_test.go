package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationSimpleEncodeAndDecode(t *testing.T) {
	log.SetOutput(io.Discard)
	// initializeSimpleDecoder()
	// getting the configs ready
	dataShards := cfg.GetInt("dataShards")
	parityShards := cfg.GetInt("parityShards")

	fileName := cfg.GetStringSlice("inputFiles.fileNames")[0]

	testDataFile := filepath.Join("testdata", fileName)

	// opening the file for reading
	file, err := os.Open(testDataFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// need file size
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	filePaths := cfg.GetStringSlice("filePaths")

	encoder, err := NewSimpleEncoder(dataShards, parityShards, fileInfo.Size())
	if err != nil {
		panic(err)
	}

	decoder := NewSimpleDecoder(dataShards, parityShards)

	encoder.Encode(file, fileName, filePaths)
	decoder.Decode(fileName, filePaths)
}

// func BenchmarkSimpleEncodeAndDecode(b *testing.B) {
// 	log.SetOutput(io.Discard)
// 	initializeSimpleDecoder()
// 	// getting the configs ready
// 	dataShards := cfg.GetInt("dataShards")
// 	parityShards := cfg.GetInt("parityShards")

// 	fileName := cfg.GetStringSlice("inputFiles.fileNames")[0]

// 	testDataFile := filepath.Join("testdata", fileName)

// 	// opening the file for reading
// 	file, err := os.Open(testDataFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	// need file size
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		panic(err)
// 	}

// 	filePaths := cfg.GetStringSlice("filePaths")

// 	encoder, err := NewSimpleEncoder(dataShards, parityShards, fileInfo.Size())
// 	if err != nil {
// 		panic(err)
// 	}

// 	decoder := NewSimpleDecoder(dataShards, parityShards)

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		encoder.Encode(file, fileName, filePaths)
// 		_, fileSize := decoder.Decode(fileName, filePaths)
// 		b.SetBytes(int64(fileSize))
// 	}
// }

