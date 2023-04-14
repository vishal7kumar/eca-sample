package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type streamingEncoderTest struct {
	dataShards      int
	parityShards    int
	inputFileNames  []string
	inputFileSize   []string
	inputFilePath   string
	outputFilePaths []string
}

func NewStreamingEncoderTest() streamingEncoderTest {

	initializeConfig()

	dataShards := cfg.GetInt("dataShards")
	parityShards := cfg.GetInt("parityShards")
	inputFileNames := cfg.GetStringSlice("inputFiles.fileNames")
	inputFileSizes := cfg.GetStringSlice("inputFiles.fileSizes")
	inputFilePath := cfg.GetString("inputFiles.filePath")
	outputFilePaths := cfg.GetStringSlice("outputFilePaths")

	return streamingEncoderTest{dataShards: dataShards, parityShards: parityShards, inputFileNames: inputFileNames, inputFileSize: inputFileSizes,
		inputFilePath: inputFilePath, outputFilePaths: outputFilePaths}
}

func (s streamingEncoderTest) initialize() {
	// ensure folder structure is present
	createFolderStructure(s.outputFilePaths)

	// TODO: generate a large file
	inputFilePath := filepath.Join("..", s.inputFilePath)

	generateInputFiles(s.inputFileSize, s.inputFileNames, inputFilePath)
}

func (s streamingEncoderTest) testStreamingEncoder(t *testing.T) {
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

		encoder.Encode(file, fileName, s.outputFilePaths)

		// Test: checking for existence of file shards at the file paths specified
		for i, path := range s.outputFilePaths {
			filePath := filepath.Join(path, fileName+"."+strconv.Itoa(i))
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("expected file shard %s but not found", filePath)
			}
		}
	}
}

func (s streamingEncoderTest) benchStreamingEncoder(b *testing.B) {
	for _, fileName := range s.inputFileNames {
		data := []byte("hello world!")
		reader := NewInfiniteReader(data)

		fileSize, err := getFileSize(s.inputFileSize[0])
		if err != nil {
			panic(err)
		}

		encoder, err := NewStreamingEncoder(s.dataShards, s.parityShards, fileSize)
		if err != nil {
			panic(err)
		}

		b.SetBytes(fileSize)
		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			encoder.Encode(reader, fileName, s.outputFilePaths)
		}
		b.StopTimer()
	}
}

func TestStreamingEncoder(t *testing.T) {
	s := NewStreamingEncoderTest()
	s.initialize()
	s.testStreamingEncoder(t)
}

func BenchmarkStreamingEncoder(b *testing.B) {
	s := NewStreamingEncoderTest()
	s.initialize()
	s.benchStreamingEncoder(b)
}
