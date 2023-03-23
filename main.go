package main

import (
	"fmt"
	"github.com/vishal7kumar/eca-sample/cmd"
)

func main() {
	// values to come from caller API
	// write/ read.
	operation := cmd.Operation(cmd.Read)

	dataShards := 5
	parityShards := 3

	inputString := "This is a sample input string"

	filePath := "/home/vishal/code/golang/eca-sample/file.txt"
	cfg := cmd.GetConfig()
	// decide on switch case - decoder/encoder initialization
	switch operation {
	case cmd.Write:
		writer := cmd.NewObjectWriter(dataShards, parityShards, cfg)
		writer.Write(inputString, filePath)
	case cmd.Read:
		reader := cmd.NewObjectReader(dataShards, parityShards, cfg)
		fileContents := reader.Read(filePath)
		fmt.Printf("The contents of the file are: %s \n", string(fileContents))
	default:
		fmt.Printf("Not Implemented")
	}
}
