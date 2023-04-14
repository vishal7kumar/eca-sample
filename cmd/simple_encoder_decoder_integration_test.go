package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

type simpleEncoderDecoderTest struct {
	dataShards      int
	parityShards    int
	inputFileNames  []string
	inputFileSize   []string
	inputFilePath   string
	outputFilePaths []string
}

func NewSimpleEncoderDecoderTest() simpleEncoderDecoderTest {

	initializeConfig()

	dataShards := cfg.GetInt("dataShards")
	parityShards := cfg.GetInt("parityShards")
	inputFileNames := cfg.GetStringSlice("inputFiles.fileNames")
	inputFileSizes := cfg.GetStringSlice("inputFiles.fileSizes")
	inputFilePath := cfg.GetString("inputFiles.filePath")
	outputFilePaths := cfg.GetStringSlice("outputFilePaths")

	return simpleEncoderDecoderTest{dataShards: dataShards, parityShards: parityShards, inputFileNames: inputFileNames, inputFileSize: inputFileSizes,
		inputFilePath: inputFilePath, outputFilePaths: outputFilePaths}
}

func (s simpleEncoderDecoderTest) initialize() {
	// ensure folder structure is present
	createFolderStructure(s.outputFilePaths)

	// TODO: generate a large file
	inputFilePath := filepath.Join("..", s.inputFilePath)

	generateInputFiles(s.inputFileSize, s.inputFileNames, inputFilePath)
}

func (s simpleEncoderDecoderTest) testSimpleEncoderDecoder(t *testing.T) {
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

func (s simpleEncoderDecoderTest) benchmarkSimpleEncoderDecoder(b *testing.B) {
	log.SetOutput(io.Discard)
	for _, fileName := range s.inputFileNames {
		testDataFile := filepath.Join("testdata", fileName)

		// need file size
		file, err := os.Open(testDataFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			panic(err)
		}

		encoder, err := NewSimpleEncoder(s.dataShards, s.parityShards, fileInfo.Size())
		if err != nil {
			panic(err)
		}
		decoder := NewSimpleDecoder(s.dataShards, s.parityShards)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StartTimer()
			encoder.Encode(file, fileName, s.outputFilePaths)
			file.Seek(0, io.SeekStart)
			_, fileSize := decoder.Decode(fileName, s.outputFilePaths)
			b.StopTimer()
			b.SetBytes(int64(fileSize))
		}
	}
}

func BenchmarkSimpleEncoderDecoder(b *testing.B) {
	s := NewSimpleEncoderDecoderTest()
	s.initialize()
	s.benchmarkSimpleEncoderDecoder(b)
}

func TestSimpleEncoderDecoder(t *testing.T) {
	s := NewSimpleEncoderDecoderTest()
	s.initialize()
	s.testSimpleEncoderDecoder(t)
}