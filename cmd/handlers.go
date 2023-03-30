package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// move to responses
type Response struct {
	Data []byte `json:"data"`
}

func ReadHandler(ctx *gin.Context) {
	// fileName := ctx.GetInt()
	fileName := ctx.Param("fileName")

	dataShards, err := strconv.Atoi(ctx.Query("dataShards"))
	if err != nil {
		log.Fatal("input must be a number")
		os.Exit(2)
	}

	parityShards, err := strconv.Atoi(ctx.Query("parityShards"))
	if err != nil {
		log.Fatal("input must be a number")
		os.Exit(2)
	}

	filePaths := ctx.Query("filePaths")

	// input validations reqiured

	filePathsList := parseParams(filePaths)

	// move to validations
	if len(filePathsList) != dataShards+parityShards {
		log.Fatal("number of paths should be equal to the total shards.")
		os.Exit(2)
	}

	reader := NewObjectReader(dataShards, parityShards, GetConfig())
	fileContents := reader.Read(fileName, filePathsList)
	fmt.Printf("The contents of the file are: %s \n", string(fileContents)) // returning bytes

	response := Response{Data: fileContents}

	ctx.JSON(http.StatusOK, response)
}

func WriteHandler(ctx *gin.Context) {
	//writer := NewObjectWriter(dataShards, parityShards, cfg)
	//writer.Write(inputString, fileName, filePaths)
}

// move to utils
func parseParams(params string) []string {
	return strings.Split(params, ",")
}
