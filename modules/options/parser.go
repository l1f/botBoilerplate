package options

import (
	"botBoilerplate/messages"
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"strings"
	"time"
)

type CommandOptions struct {
	flagSet *flag.FlagSet
	tokens  []string
	writer  *stringWriter
	Name    string
}

type stringWriter struct {
	io.Writer
	builder strings.Builder
}

func (w *stringWriter) Write(b []byte) (n int, err error) {
	return w.builder.Write(b)
}

func New(commandLine string) (*CommandOptions, error) {
	options := CommandOptions{
		writer: &stringWriter{},
	}

	tokens, err := split(commandLine)
	if err != nil {
		return nil, err
	}

	if len(tokens) < 1 {
		return nil, errors.New(messages.COMMAND_COULDNT_BE_PARSED)
	}

	options.Name = tokens[0]
	options.tokens = tokens[1:]

	options.flagSet = flag.NewFlagSet(options.Name, flag.ContinueOnError)
	options.flagSet.SetOutput(options.writer)

	return &options, nil
}

func split(commandLine string) ([]string, error) {
	commandLine = strings.ReplaceAll(commandLine, "\n", " ")

	r := csv.NewReader(strings.NewReader(commandLine))
	r.Comma = ' '
	rawTokens, err := r.Read()
	if err != nil {
		return nil, err
	}

	var tokens []string
	for _, token := range rawTokens {
		if token != "" {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (c *CommandOptions) Arg(i int) string {
	return c.flagSet.Arg(i)
}

func (c *CommandOptions) Args() []string {
	return c.flagSet.Args()
}

func (c *CommandOptions) Bool(name string, value bool, usage string) *bool {
	return c.flagSet.Bool(name, value, usage)
}

func (c *CommandOptions) BoolVar(p *bool, name string, value bool, usage string) {
	c.flagSet.BoolVar(p, name, value, usage)
}

func (c *CommandOptions) Duration(name string, value time.Duration, usage string) *time.Duration {
	return c.flagSet.Duration(name, value, usage)
}

func (c *CommandOptions) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	c.flagSet.DurationVar(p, name, value, usage)
}

func (c *CommandOptions) Float64(name string, value float64, usage string) *float64 {
	return c.flagSet.Float64(name, value, usage)
}

func (c *CommandOptions) Float64Var(p *float64, name string, value float64, usage string) {
	c.flagSet.Float64Var(p, name, value, usage)
}

func (c *CommandOptions) Int(name string, value int, usage string) *int {
	return c.flagSet.Int(name, value, usage)
}

func (c *CommandOptions) IntVar(p *int, name string, value int, usage string) {
	c.flagSet.IntVar(p, name, value, usage)
}

func (c *CommandOptions) String(name string, value string, usage string) *string {
	return c.flagSet.String(name, value, usage)
}

func (c *CommandOptions) StringVar(p *string, name string, value string, usage string) {
	c.flagSet.StringVar(p, name, value, usage)
}

func (c *CommandOptions) Int64(name string, value int64, usage string) *int64 {
	return c.flagSet.Int64(name, value, usage)
}

func (c *CommandOptions) Int64Var(p *int64, name string, value int64, usage string) {
	c.flagSet.Int64Var(p, name, value, usage)
}

func (c *CommandOptions) UInt(name string, value uint, usage string) *uint {
	return c.flagSet.Uint(name, value, usage)
}

func (c *CommandOptions) UintVar(p *uint, name string, value uint, usage string) {
	c.flagSet.UintVar(p, name, value, usage)
}

func (c *CommandOptions) UInt64(name string, value uint64, usage string) *uint64 {
	return c.flagSet.Uint64(name, value, usage)
}

func (c *CommandOptions) Uint64Var(p *uint64, name string, value uint64, usage string) {
	c.flagSet.Uint64Var(p, name, value, usage)
}

func (c *CommandOptions) Lookup(name string) *flag.Flag {
	return c.flagSet.Lookup(name)
}

func (c *CommandOptions) NArg() int {
	return c.flagSet.NArg()
}

func (c *CommandOptions) NFlag() int {
	return c.flagSet.NFlag()
}

func (c *CommandOptions) Parsed() bool {
	return c.flagSet.Parsed()
}

func (c *CommandOptions) Parse() error {
	return c.flagSet.Parse(c.tokens)
}

func (c *CommandOptions) Help() string {
	c.flagSet.PrintDefaults()

	return c.writer.builder.String()
}
