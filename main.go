package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/borosr/db-benchmark/cmd"
	. "github.com/borosr/db-benchmark/x"
)

var (
	workers   int
	cycles    int
	waitAfter time.Duration
)

func loadFlags() {
	flag.IntVar(&workers, "workers", 10, "")
	flag.IntVar(&cycles, "cycles", 100, "")
	flag.DurationVar(&waitAfter, "wait", time.Second, "")
	flag.Parse()
}

func main() {
	loadFlags()

	benchmark, err := cmd.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	result := Execute(benchmark, Config{
		WaitAfter: waitAfter,
		Workers:   workers,
		Cycles:    cycles,
	})

	log.Printf("Result is: %s", result.String())
}
