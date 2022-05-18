package main

import (
	"os"

	"github.com/simomu-github/fflt_lang/interpreter"
)

func main() {
	interpreter := interpreter.New()
	os.Exit(interpreter.Run())
}
