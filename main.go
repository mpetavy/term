package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"strings"
)

var (
	ftext = flag.String("t", "", "term")
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "test", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	text := strings.ReplaceAll(*ftext, " ", "")

	term, err := NewTerm(text)
	if common.Error(err) {
		return err
	}

	result, err := term.Calc()
	if common.Error(err) {
		return err
	}

	fmt.Printf("%v\n", result)

	return nil
}

func main() {
	common.Run([]string{"t"})
}
