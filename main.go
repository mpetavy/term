package main

import (
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"strings"
)

var (
	ftext = flag.String("t", "", "term")
)

func init() {
	common.Init("term", "", "", "", "2018", "test", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
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
