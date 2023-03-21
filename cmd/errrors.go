package cmd

import (
	"errors"
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}

var invalidDataAndParitySumErr = errors.New("error: sum of data and parity shards cannot exceed 256")
