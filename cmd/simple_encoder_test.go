package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

const Mb = 1024 * 1024
const Gb = 1024 * Mb

type InfiniteReader struct {
	data []byte
	pos  int
}

func NewInfiniteReader(data []byte) *InfiniteReader {
	return &InfiniteReader{data: data}
}

func (r *InfiniteReader) Read(p []byte) (n int, err error) {
	for len(p) > 0 {
		if r.pos >= len(r.data) {
			r.pos = 0
		}
		chunkSize := copy(p, r.data[r.pos:])
		p = p[chunkSize:]
		r.pos += chunkSize
		n += chunkSize
	}
	return n, nil
}

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

func createFolderStructure(filePaths []string) {
	for _, dir := range filePaths {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Println("directory doesn't exist, creating.")
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			log.Println("error with the path.")
		}
	}
	log.Println("directory structure is intact.")
}

func generateInputFiles(fileSizes []string, fileNames []string, filePath string) error {
	data := []byte("hello world!")
	for i, file := range fileNames {
		f, err := os.Create(filepath.Join(filePath, file))
		if err != nil {
			log.Println("could not create file, error: ", err.Error())
			return err
		}
		reader := NewInfiniteReader(data)
		fileSize, err := getFileSize(fileSizes[i])
		if err != nil {
			return err
		}
		_, err = io.CopyN(f, reader, fileSize)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFileSize(size string) (int64, error) {
    var multiplier int64
    switch size[len(size)-1] {
    case 'G':
        multiplier = Gb
    case 'M':
        multiplier = Mb
    default:
        return 0, fmt.Errorf("unsupported size format: %s", size)
    }
    value, err := strconv.ParseInt(size[:len(size)-1], 10, 64)
    if err != nil {
        return 0, err
    }
    return value * multiplier, nil
}


func TestSimpleEncoder(t *testing.T) {
	s := NewSimpleEncoderTest()
	s.initialize()
	s.testSimpleEncoder(t)
}
