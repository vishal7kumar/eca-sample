package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// type Tester interface {
// 	Test(t *testing.T)
// 	Benchmark(b *testing.B)
// }

type simpleEncoderTest struct {
	dataShards      int
	parityShards    int
	inputFileNames  []string
	inputFileSize   []string
	inputFilePath   string
	outputFilePaths []string
}

func NewSimpleEncoderTest() simpleEncoderTest {

	initializeConfig()

	dataShards := cfg.GetInt("dataShards")
	parityShards := cfg.GetInt("parityShards")
	inputFileNames := cfg.GetStringSlice("inputFiles.fileNames")
	inputFileSizes := cfg.GetStringSlice("inputFiles.fileSizes")
	inputFilePath := cfg.GetString("inputFiles.filePath")
	outputFilePaths := cfg.GetStringSlice("outputFilePaths")

	return simpleEncoderTest{dataShards: dataShards, parityShards: parityShards, inputFileNames: inputFileNames, inputFileSize: inputFileSizes,
		inputFilePath: inputFilePath, outputFilePaths: outputFilePaths}
}

func (s simpleEncoderTest) initialize() {
	// ensure folder structure is present
	createFolderStructure(s.outputFilePaths)

	// TODO: generate a large file
	inputFilePath := filepath.Join("..", s.inputFilePath)

	generateInputFiles(s.inputFileSize, s.inputFileNames, inputFilePath)
}

func (s simpleEncoderTest) testSimpleEncoder(t *testing.T) {
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

		encoder, err := NewSimpleEncoder(s.dataShards, s.parityShards, fileInfo.Size())
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

func (s simpleEncoderTest) benchSimpleEncoder(b *testing.B) {
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

		encoder, err := NewSimpleEncoder(s.dataShards, s.parityShards, fileInfo.Size())
		if err != nil {
			panic(err)
		}

		for i := 0; i < b.N; i++ {
			encoder.Encode(file, fileName, s.outputFilePaths)
			file.Seek(0, io.SeekStart)
		}
	}
}

func TestSimpleEncoder(t *testing.T) {
	s := NewSimpleEncoderTest()
	s.initialize()
	s.testSimpleEncoder(t)
}

func BenchmarkSimpleEncoder(b *testing.B) {
	s := NewSimpleEncoderTest()
	s.initialize()
	s.benchSimpleEncoder(b)
}
