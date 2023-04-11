package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

type streamingEncoderDecoderIntegrationTest struct {
	dataShards      int
	parityShards    int
	inputFileNames  []string
	inputFileSize   []string
	inputFilePath   string
	outputFilePaths []string
}

func NewStreamingEncoderDecoderIntegrationTest() streamingEncoderDecoderIntegrationTest {

	initializeConfig()

	dataShards := cfg.GetInt("dataShards")
	parityShards := cfg.GetInt("parityShards")
	inputFileNames := cfg.GetStringSlice("inputFiles.fileNames")
	inputFileSizes := cfg.GetStringSlice("inputFiles.fileSizes")
	inputFilePath := cfg.GetString("inputFiles.filePath")
	outputFilePaths := cfg.GetStringSlice("outputFilePaths")

	return streamingEncoderDecoderIntegrationTest{dataShards: dataShards, parityShards: parityShards, inputFileNames: inputFileNames, inputFileSize: inputFileSizes,
		inputFilePath: inputFilePath, outputFilePaths: outputFilePaths}
}

func (s streamingEncoderDecoderIntegrationTest) initialize() {
	// ensure folder structure is present
	createFolderStructure(s.outputFilePaths)

	// TODO: generate a large file
	inputFilePath := filepath.Join("..", s.inputFilePath)

	generateInputFiles(s.inputFileSize, s.inputFileNames, inputFilePath)
}

func (s streamingEncoderDecoderIntegrationTest) benchmarkStreamingEncoderDecoderIntegrationTest(b *testing.B) {
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

		for i := 0; i < b.N; i++ {
			encoder.Encode(file, fileName, s.outputFilePaths)
			_, fileSize := decoder.Decode(fileName, s.outputFilePaths)
			b.SetBytes(int64(fileSize))
		}
	}
}

func BenchmarkStreamingEncoderDecoderIntegrationTest(b *testing.B) {
	s := NewStreamingEncoderDecoderIntegrationTest()
	s.initialize()
	s.benchmarkStreamingEncoderDecoderIntegrationTest(b)
}
