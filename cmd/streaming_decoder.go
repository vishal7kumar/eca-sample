package cmd

// import (
// 	"fmt"
// 	"github.com/klauspost/reedsolomon"
// 	"io"
// 	"os"
// 	"path/filepath"
// )

// type streamingDecoder struct {
// 	enc          reedsolomon.StreamEncoder
// 	dataShards   int
// 	parityShards int
// }

// func NewStreamingDecoder(dataShards int, parityShards int) DecoderService {
// 	enc, err := reedsolomon.NewStream(dataShards, parityShards)
// 	checkErr(err)
// 	return &streamingDecoder{enc: enc, dataShards: dataShards, parityShards: parityShards}
// }

// func (e streamingDecoder) Decode(fileName string, filePaths []string) string {
// 	shards, size, err := openInput(e.dataShards, e.parityShards, fileName, filePaths)
// 	checkErr(err)

// 	// Verify the shards
// 	ok, err := e.enc.Verify(shards)
// 	if ok {
// 		fmt.Println("No reconstruction needed")
// 	} else {
// 		fmt.Println("Verification failed. Reconstructing data")
// 		shards, size, err = openInput(e.dataShards, e.parityShards, fileName, filePaths)
// 		checkErr(err)
// 		// Create out destination writers
// 		out := make([]io.Writer, len(shards))
// 		paths := len(filePaths)
// 		for i := range out {
// 			if shards[i] == nil {
// 				outputFile := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i%paths], fileName), i)
// 				fmt.Println("Creating", outputFile)
// 				out[i], err = os.Create(outputFile)
// 				checkErr(err)
// 			}
// 		}
// 		err = e.enc.Reconstruct(shards, out)
// 		if err != nil {
// 			fmt.Println("Reconstruct failed -", err)
// 			os.Exit(1)
// 		}
// 		// Close output.
// 		for i := range out {
// 			if out[i] != nil {
// 				err := out[i].(*os.File).Close()
// 				checkErr(err)
// 			}
// 		}
// 		shards, size, err = openInput(e.dataShards, e.parityShards, fileName, filePaths)
// 		ok, err = e.enc.Verify(shards)
// 		if !ok {
// 			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
// 			os.Exit(1)
// 		}
// 		checkErr(err)
// 	}

// 	// TODO: Write data to the PWD ?
// 	cwd, err := os.Getwd()
// 	checkErr(err)

// 	outFilePath := filepath.Join(cwd, fileName)
// 	fmt.Println("Writing data to", filepath.Join(cwd, fileName))
// 	f, err := os.Create(outFilePath)
// 	checkErr(err)

// 	shards, size, err = openInput(e.dataShards, e.parityShards, fileName, filePaths)
// 	checkErr(err)

// 	// We don't know the exact filesize.
// 	err = e.enc.Join(f, shards, int64(e.dataShards)*size)
// 	checkErr(err)

// 	return outFilePath
// }

// func openInput(dataShards, parShards int, fileName string, filePaths []string) (r []io.Reader, size int64, err error) {
// 	// Create shards and load the data.
// 	paths := len(filePaths)

// 	shards := make([]io.Reader, dataShards+parShards)
// 	for i := range shards {
// 		inputFileName := fmt.Sprintf("%s.%d", filepath.Join(filePaths[i%paths], fileName), i)
// 		fmt.Println("Opening", inputFileName)
// 		f, err := os.Open(inputFileName)
// 		if err != nil {
// 			fmt.Println("Error reading file", err)
// 			shards[i] = nil
// 			continue
// 		} else {
// 			shards[i] = f
// 		}
// 		stat, err := f.Stat()
// 		checkErr(err)
// 		if stat.Size() > 0 {
// 			size = stat.Size()
// 		} else {
// 			shards[i] = nil
// 		}
// 	}
// 	return shards, size, nil
// }
