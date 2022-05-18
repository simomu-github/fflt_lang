package parser

import (
	"errors"
	"fmt"

	"github.com/simomu-github/fflt_lang/lexer"
)

func parseError(state parseState, token lexer.Token, message string) error {
	errorMessage := fmt.Sprintf("Syntex error: %s at %s:%d:%d", message, state.filename, token.Line, token.Column)
	return errors.New(errorMessage)
}
