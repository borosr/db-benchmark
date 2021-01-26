package cmd

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	mysql     cmd = "mysql"
	couchbase cmd = "couchbase"
	cockroach cmd = "cockroach"
)

var (
	all           = []string{string(mysql), string(couchbase), string(cockroach)}
	validCmdRegex = regexp.MustCompile("(" + strings.Join(all, "|") + ")")

	commands = map[cmd]Benchmark{
		mysql:     nil,
		couchbase: nil,
		cockroach: nil,
	}
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

	return commands[cmd(args[1])], nil
}
