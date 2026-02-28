package main

import (
	"os"
	"strconv"
)

func mustInt(key string) int {
	v, _ := strconv.Atoi(os.Getenv(key))
	return v
}
