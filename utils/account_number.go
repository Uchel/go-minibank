package utils

import (
	"log"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

func AccountNumberGenerate() string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Println(err)

	}

	// Set seed for random number generation

	// Generate a new unique ID within a specific range
	min := int64(100000000000)
	max := int64(999999999999)
	id := node.Generate().Int64()
	randomID := min + (id % (max - min + 1))

	AccountNumber := strconv.Itoa(int(randomID))

	return AccountNumber
}
