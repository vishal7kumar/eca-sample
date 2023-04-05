package cmd

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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

	filePathsList := parseParams(filePaths)

	err = validateParams(dataShards, parityShards, filePathsList)
	if err != nil {
		errorResponse := Response{Error: &ErrorResponse{Error: err.Error()}}
		ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	objReader := NewObjectReader(dataShards, parityShards, fileName, filePathsList, GetConfig())
	reader, size := objReader.Read(fileName, filePathsList)
	ctx.Header("Content-Type", "text/plain")

	ctx.DataFromReader(http.StatusOK, size, "text/plain", reader, nil)
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

	filePathsList := parseParams(filePaths)

	err = validateParams(dataShards, parityShards, filePathsList)
	if err != nil {
		errorResponse := Response{Error: &ErrorResponse{Error: err.Error()}}
		ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	requestBody := ctx.Request.Body

	contentLength := ctx.Request.Header.Get("content-length")
	contentLengthInBytes, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Fatal("content length should be a number.")
		os.Exit(1)
	}

	if contentLengthInBytes == 0 {
		log.Fatal("data size is required for encoding operation.")
		os.Exit(1)
	}

	size := int64(contentLengthInBytes)

	writer := NewObjectWriter(dataShards, parityShards, size, GetConfig())
	writer.Write(requestBody, fileName, filePathsList)

	ctx.JSON(http.StatusCreated, nil)
}

func parseParams(params string) []string {
	return strings.Split(params, ",")
}

func validateParams(dataShards int, parityShards int, filePathsList []string) error {
	if len(filePathsList) != dataShards+parityShards {
		return ErrUnequalPathsAndShards
	}
	return nil
}
