package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vishal7kumar/eca-sample/cmd"
)

func main() {

	router := gin.Default()

	router.GET("/api/v1/files/:fileName", cmd.ReadHandler)
	router.POST("/api/v1/files/:fileName", cmd.WriteHandler)
	err := router.Run(cmd.RestApiPort)
	if err != nil {
		panic(fmt.Errorf("API server run failed: %v", err))
	}
}
