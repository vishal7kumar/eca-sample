package main

import (
	"fmt"
	"github.com/vishal7kumar/eca-sample/cmd"
)

func main() {
	// values to come from caller API
	// write/ read.
	operation := cmd.Operation(cmd.Read)

	dataShards := 5   // always from the parameter of the request for the write operation.
	parityShards := 3 // for the read operation, a module having persistent in-memory DB.

	inputString := "This is a sample input string"

	filePaths := []string{"/home/vishal/code/golang/eca-sample/output/a",
		"/home/vishal/code/golang/eca-sample/output/b"}
	fileName := "file.txt" // assuming to be in the current directory
	cfg := cmd.GetConfig()
	// decide on switch case - decoder/encoder initialization
	switch operation {
	case cmd.Write:
		writer := cmd.NewObjectWriter(dataShards, parityShards, cfg)
		writer.Write(inputString, fileName, filePaths)
	case cmd.Read:
		reader := cmd.NewObjectReader(dataShards, parityShards, cfg)
		fileContents := reader.Read(fileName, filePaths)
		fmt.Printf("The contents of the file are: %s \n", string(fileContents)) // returning bytes
		// read from the streams and write back to it.
		// no need to store in the memory.
	default:
		fmt.Printf("Not Implemented")
	}
}
