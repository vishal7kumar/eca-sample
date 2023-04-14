package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

type streamingEncoderDecoderIntegration struct {
	dataShards      int
	parityShards    int
	inputFileNames  []string
	inputFileSize   []string
	inputFilePath   string
	outputFilePaths []string
}

func NewStreamingEncoderDecoderIntegration() streamingEncoderDecoderIntegration {

	initializeConfig()

	dataShards := cfg.GetInt("dataShards")
	parityShards := cfg.GetInt("parityShards")
	inputFileNames := cfg.GetStringSlice("inputFiles.fileNames")
	inputFileSizes := cfg.GetStringSlice("inputFiles.fileSizes")
	inputFilePath := cfg.GetString("inputFiles.filePath")
	outputFilePaths := cfg.GetStringSlice("outputFilePaths")

	return streamingEncoderDecoderIntegration{dataShards: dataShards, parityShards: parityShards, inputFileNames: inputFileNames, inputFileSize: inputFileSizes,
		inputFilePath: inputFilePath, outputFilePaths: outputFilePaths}
}

func (s streamingEncoderDecoderIntegration) initialize() {
	// ensure folder structure is present
	createFolderStructure(s.outputFilePaths)

	inputFilePath := filepath.Join("..", s.inputFilePath)

	generateInputFiles(s.inputFileSize, s.inputFileNames, inputFilePath)
}

func (s streamingEncoderDecoderIntegration) testStreamingEncoderDecoderIntegration(t *testing.T) {
	for _, fileName := range s.inputFileNames {
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

		encoder, err := NewStreamingEncoder(s.dataShards, s.parityShards, fileInfo.Size())
		if err != nil {
			panic(err)
		}

		decoder := NewStreamingDecoder(s.dataShards, s.parityShards)
		if err != nil {
			panic(err)
		}

		encoder.Encode(file, fileName, s.outputFilePaths)
		decoder.Decode(fileName, s.outputFilePaths)
	}
}

func (s streamingEncoderDecoderIntegration) benchmarkStreamingEncoderDecoderIntegrationTest(b *testing.B) {
	for _, fileName := range s.inputFileNames {
		data := []byte("hello")
		reader := NewInfiniteReader(data)

		fileSize, err := getFileSize(s.inputFileSize[0])
		if err != nil {
			panic(err)
		}

		encoder, err := NewStreamingEncoder(s.dataShards, s.parityShards, fileSize)
		if err != nil {
			panic(err)
		}

		decoder := NewStreamingDecoder(s.dataShards, s.parityShards)
		if err != nil {
			panic(err)
		}

		b.SetBytes(int64(fileSize))
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			encoder.Encode(reader, fileName, s.outputFilePaths)
			decoder.Decode(fileName, s.outputFilePaths)
		}
		b.StopTimer()
	}
}

func BenchmarkStreamingEncoderDecoderIntegrationTest(b *testing.B) {
	s := NewStreamingEncoderDecoderIntegration()
	s.initialize()
	s.benchmarkStreamingEncoderDecoderIntegrationTest(b)
}

func TestStreamingEncoderDecoderIntegration(t *testing.T) {
	s := NewStreamingEncoderDecoderIntegration()
	s.initialize()
	s.testStreamingEncoderDecoderIntegration(t)
}
