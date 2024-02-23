package interpreter

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/simomu-github/fflt_lang/executor"
	"github.com/simomu-github/fflt_lang/lexer"
	"github.com/simomu-github/fflt_lang/parser"
)

var (
	versionOpt = flag.Bool("v", false, "display version information")
)

const version = "v0.0.2"

type Interpreter struct {
	args   []string
	stderr io.Writer
}

func New() *Interpreter {
	return &Interpreter{
		args:   os.Args,
		stderr: os.Stderr,
	}
}

func (i *Interpreter) Run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  fflt_lang [FILE]\n", os.Args[0])
		flag.PrintDefaults()
	}

	if len(i.args) < 2 {
		flag.Usage()
		return 1
	}

	flag.Parse()
	if *versionOpt {
		fmt.Printf("fflt_lang version %s\n", version)
		return 1
	}

	filename := i.args[1]

	bytes, errReadFile := ioutil.ReadFile(filename)
	if errReadFile != nil {
		fmt.Fprintf(i.stderr, "%s can not read\n", filename)
		return 1
	}

	tokens, lexErr := lexer.ScanAllTokens(string(bytes), filename)
	if lexErr != nil {
		fmt.Fprintln(i.stderr, lexErr.Error())
		return 1
	}

	instructions, labelMap, parseErr := parser.ParseAll(tokens, filename)
	if parseErr != nil {
		fmt.Fprintln(i.stderr, parseErr.Error())
		return 1
	}

	executor := executor.Executor{
		Filename:     filename,
		Instructions: instructions,
		LabelMap:     labelMap,
		Input: func() string {
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			return stdin.Text()
		},
		Output: func(str string) {
			fmt.Printf(str)
		},
	}
	errRuntime := executor.Run()
	if errRuntime != nil {
		fmt.Fprintln(i.stderr, errRuntime.Error())
		return 1
	}

	return 0
}
