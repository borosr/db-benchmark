package cmd

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/borosr/db-benchmark/drivers/cockroach"
	"github.com/borosr/db-benchmark/drivers/mysql"
)

const (
	mysqlCmd     cmd = "mysql"
	couchbaseCmd cmd = "couchbase"
	cockroachCmd cmd = "cockroach"
)

var (
	all           = []string{string(mysqlCmd), string(couchbaseCmd), string(cockroachCmd)}
	validCmdRegex = regexp.MustCompile("(" + strings.Join(all, "|") + ")")

	//commands = map[cmd]Benchmark{
	//	mysqlCmd:     mysql.New(),
	//	couchbaseCmd: nil,
	//	cockroachCmd: cockroach.New(),
	//}
)

type cmd string

type Benchmark interface {
	Write(ctx context.Context) (string, error)
	Read(ctx context.Context, ID string) error
}

func Parse(args []string) (Benchmark, error) {
	if l := len(args); l < 2 {
		return nil, errors.New("missing command line argument")
	} else if l > 2 {
		return nil, errors.New("too many arguments")
	}

	if !validCmdRegex.MatchString(args[1]) {
		return nil, fmt.Errorf("invalid cmd [%s], should be one of [%s]", args[1], strings.Join(all, ","))
	}

	//return commands[cmd(args[1])], nil
	switch cmd(args[1]) {
	case mysqlCmd:
		return mysql.New(), nil
	case cockroachCmd:
		return cockroach.New(), nil
	}

	return nil, errors.New("invalid command type")
}
