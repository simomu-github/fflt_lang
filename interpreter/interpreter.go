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
	dumpOpt    = flag.Bool("dump", false, "disassemble instructions")
	debugOpt   = flag.Bool("debug", false, "debugger")
)

const version = "v0.0.2"

type Interpreter struct {
	stderr io.Writer
}

func New() *Interpreter {
	return &Interpreter{
		stderr: os.Stderr,
	}
}

func (i *Interpreter) Run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  fflt_lang [FILE]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if *versionOpt {
		fmt.Printf("fflt_lang version %s\n", version)
		return 1
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		return 1
	}

	filename := flag.Arg(0)

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

	exe := executor.Executor{
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

	if *dumpOpt {
		exe.Disassenble()
		return 0
	}

	if *debugOpt {
		debugger := executor.NewDebugger(&exe)
		debugger.Run()
		return 0
	}

	errRuntime := exe.Run()
	if errRuntime != nil {
		fmt.Fprintln(i.stderr, errRuntime.Error())
		return 1
	}

	return 0
}
