package main

import (
	"cron/pkg/cron"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("a single argument is required")
	}
	pattern, err := cron.ParseConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Println(pattern)
}
