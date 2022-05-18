package parser

import (
	"fmt"

	"github.com/simomu-github/fflt_lang/executor"
	"github.com/simomu-github/fflt_lang/lexer"
)

type parseState struct {
	filename     string
	instructions []executor.Instruction
	labelMap     map[string]int
}

func ParseAll(tokens []lexer.Token, filename string) ([]executor.Instruction, map[string]int, error) {
	state := parseState{
		filename:     filename,
		instructions: []executor.Instruction{},
		labelMap:     map[string]int{},
	}

	for i := 0; i < len(tokens); i++ {
		var err error

		state, i, err = parseToken(state, tokens, i, filename)

		if err != nil {
			return nil, nil, err
		}
	}
	return state.instructions, state.labelMap, nil
}

func parseToken(state parseState, tokens []lexer.Token, index int, filename string) (parseState, int, error) {
	var err error

	if requireArgumentTokens(tokens[index]) {
		if index+1 >= len(tokens) {
			return state, index, parseError(state, tokens[index], "expected parameter token")
		}
		state, err = parseTokenWithParameter(state, tokens[index], tokens[index+1])
		index++
	} else {
		state, err = parseSingleToken(state, tokens[index])
	}

	return state, index, err
}

func parseSingleToken(state parseState, token lexer.Token) (parseState, error) {
	switch token.Type {
	case lexer.Swap:
		state.instructions = append(state.instructions, executor.Swap{Token: token})
	case lexer.Duplicate:
		state.instructions = append(state.instructions, executor.Duplicate{Token: token})
	case lexer.Discard:
		state.instructions = append(state.instructions, executor.Discard{})

	case lexer.Addition:
		state.instructions = append(state.instructions, executor.Addition{Token: token})
	case lexer.Subtraction:
		state.instructions = append(state.instructions, executor.Subtraction{Token: token})
	case lexer.Multiplication:
		state.instructions = append(state.instructions, executor.Multiplication{Token: token})
	case lexer.Division:
		state.instructions = append(state.instructions, executor.Division{Token: token})
	case lexer.Modulo:
		state.instructions = append(state.instructions, executor.Modulo{Token: token})

	case lexer.Store:
		state.instructions = append(state.instructions, executor.Store{Token: token})
	case lexer.Retrieve:
		state.instructions = append(state.instructions, executor.Retrieve{Token: token})

	case lexer.Putc:
		state.instructions = append(state.instructions, executor.Putc{Token: token})
	case lexer.Putn:
		state.instructions = append(state.instructions, executor.Putn{Token: token})
	case lexer.Getc:
		state.instructions = append(state.instructions, executor.Getc{Token: token})
	case lexer.Getn:
		state.instructions = append(state.instructions, executor.Getn{Token: token})

	case lexer.EndSubroutine:
		state.instructions = append(state.instructions, executor.EndSubroutine{Token: token})
	case lexer.EndProgram:
		state.instructions = append(state.instructions, executor.EndProgram{})
	}

	return state, nil
}

func parseTokenWithParameter(state parseState, token lexer.Token, nextToken lexer.Token) (parseState, error) {
	switch token.Type {
	case lexer.Push:
		value, err := parseNumber(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.Push{Value: value})
	case lexer.Copy:
		value, err := parseNumber(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.Copy{Token: token, Value: value})
	case lexer.Slide:
		value, err := parseNumber(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.Slide{Token: token, Value: value})
	case lexer.MarkLabel:
		label, err := parseLabel(state, nextToken)
		if err != nil {
			return state, err
		}

		state.labelMap[label] = len(state.instructions)
		state.instructions = append(state.instructions, executor.MarkLabel{})
	case lexer.CallSubroutine:
		label, err := parseLabel(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.CallSubroutine{Token: token, Label: label})
	case lexer.JumpLabel:
		label, err := parseLabel(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.JumpLabel{Token: token, Label: label})
	case lexer.JumpLabelWhenZero:
		label, err := parseLabel(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.JumpLabelWhenZero{Token: token, Label: label})
	case lexer.JumpLabelWhenNegative:
		label, err := parseLabel(state, nextToken)
		if err != nil {
			return state, err
		}

		state.instructions = append(state.instructions, executor.JumpLabelWhenNegative{Token: token, Label: label})
	}

	return state, nil
}

func parseNumber(state parseState, token lexer.Token) (int, error) {
	if token.Type != lexer.Number {
		return 0, parseError(state, token, fmt.Sprintf("expected number parameter, but actual %s token", token.Type))
	}

	if len(token.Literal) == 0 {
		return 0, parseError(state, token, "expected number parameter")
	}

	char := token.Literal[0]
	var sign int

	switch string(char) {
	case lexer.UpperF, lexer.LowerF:
		sign = 1
	case lexer.UpperL, lexer.LowerL:
		sign = -1
	default:
		return 0, parseError(state, token, "expected sign")
	}

	n, err := parseBinaryNumber(state, token, 0, 0)

	return sign * n, err
}

func parseBinaryNumber(state parseState, token lexer.Token, n int, digits int) (int, error) {
	if digits+1 >= len(token.Literal) {
		return 0, parseError(state, token, fmt.Sprintf("expected numeric parameters end with a \"T\" or \"t\""))
	}

	char := token.Literal[digits+1]
	switch string(char) {
	case lexer.UpperF, lexer.LowerF:
		return parseBinaryNumber(state, token, 2*n, digits+1)
	case lexer.UpperL, lexer.LowerL:
		return parseBinaryNumber(state, token, 2*n+1, digits+1)
	case lexer.UpperT, lexer.LowerT:
		if digits > 0 {
			return n, nil
		} else {
			return 0, parseError(state, token, fmt.Sprintf("expected number parameter"))
		}
	default:
		return 0, parseError(state, token, fmt.Sprintf("expected numeric parameters end with a \"T\" or \"t\""))
	}
}

func parseLabel(state parseState, token lexer.Token) (string, error) {
	if token.Type != lexer.Label {
		return "", parseError(state, token, fmt.Sprintf("expected label parameter, but actual %s token", token.Type))
	}

	if len(token.Literal) <= 1 {
		return "", parseError(state, token, fmt.Sprintf("expected label parameters"))
	}

	lastChar := string(token.Literal[len(token.Literal)-1])
	if lastChar != lexer.UpperT && lastChar != lexer.LowerT {
		return "", parseError(state, token, fmt.Sprintf("expected label parameters end with a \"T\" or \"t\""))
	}

	label := token.Literal[:len(token.Literal)-1]
	return label, nil
}

func requireArgumentTokens(token lexer.Token) bool {
	return token.Type == lexer.Push ||
		token.Type == lexer.Copy ||
		token.Type == lexer.Slide ||
		token.Type == lexer.MarkLabel ||
		token.Type == lexer.CallSubroutine ||
		token.Type == lexer.JumpLabel ||
		token.Type == lexer.JumpLabelWhenZero ||
		token.Type == lexer.JumpLabelWhenNegative
}
