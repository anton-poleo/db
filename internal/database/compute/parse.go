package compute

import (
	"errors"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

var EmptyCmdError = errors.New("empty command")
var UnknownCmdError = errors.New("unknown command")
var IncorrectArgsNumberError = errors.New("incorrect number of arguments")
var InvalidArgsError = errors.New("invalid characters for argument")

var validArg = regexp.MustCompile(`^[a-zA-Z0-9/_*]+$`)

type Query struct {
	name string
	args []string
}

func NewQuery(name string, args []string) Query {
	return Query{name, args}
}

func (q *Query) Name() string {
	return q.name
}

func (q *Query) Args() []string {
	return q.args
}

type Parser struct {
	log *zap.Logger
}

func NewParser(log *zap.Logger) *Parser {
	return &Parser{log}
}

func (p *Parser) ParseCommand(input string) (Query, error) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return Query{}, EmptyCmdError
	}
	cmd, args := tokens[0], tokens[1:]
	if cnt := getCommandArgsNum(cmd); cnt == -1 {
		return Query{}, UnknownCmdError
	} else if len(args) != cnt {
		return Query{}, IncorrectArgsNumberError
	} else if !validArg.MatchString(args[0]) {
		return Query{}, InvalidArgsError
	}
	return NewQuery(cmd, args), nil
}
