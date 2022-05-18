package executor

import (
	"errors"
	"fmt"

	"github.com/simomu-github/fflt_lang/lexer"
)

func runtimeError(executor *Executor, message string) error {
	errorMessage := fmt.Sprintf("Runtime error: %s", message)
	return errors.New(errorMessage)
}

func runtimeErrorWithToken(executor *Executor, token lexer.Token, message string) error {
	errorMessage := fmt.Sprintf("Runtime error: %s at %s:%d:%d", message, executor.Filename, token.Line, token.Column)
	return errors.New(errorMessage)
}
