package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
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

var cfg *viper.Viper

func initializeConfig() {
	os.Setenv("STAGE", "test")
	cfg = GetConfig()
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
		defer f.Close()
		reader := NewInfiniteReader(data)
		fileSize, err := getFileSize(fileSizes[i])
		if err != nil {
			return err
		}
		_, err = io.CopyN(f, reader, fileSize)
		if err != nil {
			return err
		}
		err = f.Sync()
		if err != nil {
			panic(err)
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
