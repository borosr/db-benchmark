package x

import (
	"context"
	"log"
	"time"

	"github.com/borosr/db-benchmark/cmd"
)

type Config struct {
	WaitAfter time.Duration
	Workers   int
	Cycles    int
}

type Result struct {
	All      time.Duration
	AvgCycle time.Duration
	AvgWrite time.Duration
	AvgRead  time.Duration
}

type subResult struct {
	all   time.Duration
	write time.Duration
	read  time.Duration
	err   error
}

func Execute(b cmd.Benchmark, c Config) Result {
	start := time.Now()
	ctx := context.Background()

	jobs := make(chan struct{}, c.Cycles)
	results := make(chan subResult, c.Cycles)

	for i := 0; i < c.Workers; i++ {
		go execute(ctx, b, c, jobs, results)
	}

	for i := 0; i < c.Cycles; i++ {
		jobs <- struct{}{}
	}

	close(jobs)

	var sumCycles int64
	var sumWrite int64
	var sumRead int64
	var steps int
	for i := 0; i < c.Cycles; i++ {
		r := <- results
		if r.err != nil {
			log.Printf("Error in subresult: %v", r.err)

			continue
		}
		sumCycles += int64(r.all)
		sumWrite += int64(r.write)
		sumRead += int64(r.read)

		steps++
	}

	return Result{
		All:      time.Since(start),
		AvgCycle: time.Duration(sumCycles / int64(steps)),
		AvgWrite: time.Duration(sumWrite / int64(steps)),
		AvgRead:  time.Duration(sumRead / int64(steps)),
	}
}

func execute(ctx context.Context, b cmd.Benchmark, c Config, jobs chan struct{}, results chan subResult) {
	for range jobs {
		start := time.Now()
		id, err := b.Write(ctx)
		subRes := subResult{write: time.Since(start)}
		if err != nil {
			subRes.err = err
			subRes.all = time.Since(start)
			results <- subRes
			continue
		}
		rStart := time.Now()
		if err := b.Read(ctx, id); err != nil {
			subRes.err = err
			subRes.all = time.Since(start)
			subRes.read = time.Since(rStart)
			results <- subRes
			continue
		}
		subRes.all = time.Since(start)
		subRes.read = time.Since(rStart)
		results <- subRes
		time.Sleep(c.WaitAfter)
	}
}
