package lexer

import (
	"errors"
	"fmt"
)

func lexicalError(lexer *Lexer, message string) error {
	errorMessage := fmt.Sprintf("Syntex error: %s at %s:%d:%d", message, lexer.filename, lexer.currentLine, lexer.currentColumn)
	return errors.New(errorMessage)
}
