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

var ErrInvalidDataAndParitySum = errors.New("error: sum of data and parity shards cannot exceed 256")

var ErrUnequalPathsAndShards = errors.New("number of paths should be equal to the total shards")

var ErrNoValidShardFound = errors.New("no readable shard found")