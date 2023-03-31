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
	fileName := ctx.Param("fileName")

	dataShards, err := strconv.Atoi(ctx.Query("dataShards"))
	if err != nil {
		log.Fatal("input must be a number")
		os.Exit(1)
	}

	parityShards, err := strconv.Atoi(ctx.Query("parityShards"))
	if err != nil {
		log.Fatal("input must be a number")
		os.Exit(1)
	}

	filePaths := ctx.Query("filePaths")

	// input validations reqiured

	filePathsList := parseParams(filePaths)

	// move to validations
	if len(filePathsList) != dataShards+parityShards {
		log.Fatal("number of paths should be equal to the total shards.")
		os.Exit(1)
	}

	reader := NewObjectReader(dataShards, parityShards, GetConfig())
	fileContents := reader.Read(fileName, filePathsList)
	fmt.Printf("The contents of the file are: %s \n", string(fileContents)) // returning bytes

	response := Response{Data: fileContents}

	ctx.JSON(http.StatusOK, response)
}

func WriteHandler(ctx *gin.Context) {
	fileName := ctx.Param("fileName")

	dataShards, err := strconv.Atoi(ctx.Query("dataShards"))
	if err != nil {
		log.Fatal("input must be a number")
		os.Exit(1)
	}

	parityShards, err := strconv.Atoi(ctx.Query("parityShards"))
	if err != nil {
		log.Fatal("input must be a number")
		os.Exit(1)
	}

	filePaths := ctx.Query("filePaths")

	// input validations reqiured

	filePathsList := parseParams(filePaths)

	requestBody := ctx.Request.Body

	contentLength := ctx.Request.Header.Get("content-length")
	contentLengthInBytes, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Fatal("content length should be a number.")
		os.Exit(1)
	}

	// move to validations
	if contentLengthInBytes == 0 {
		log.Fatal("data size is required for encoding operation.")
		os.Exit(1)
	}

	size := int64(contentLengthInBytes)

	if len(filePathsList) != dataShards+parityShards {
		log.Fatal("number of paths should be equal to the total shards.")
		os.Exit(1)
	}

	writer := NewObjectWriter(dataShards, parityShards, size, GetConfig())
	writer.Write(requestBody, fileName, filePathsList)

	ctx.JSON(http.StatusCreated, nil)
}

// move to utils
func parseParams(params string) []string {
	return strings.Split(params, ",")
}
