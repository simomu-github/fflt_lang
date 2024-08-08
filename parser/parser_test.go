package parser

import (
	"testing"

	"github.com/simomu-github/fflt_lang/executor"
	"github.com/simomu-github/fflt_lang/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseNumber(t *testing.T) {
	state := parseState{
		filename:     "",
		instructions: []executor.Instruction{},
		labelMap:     map[string]int{},
	}

	token := lexer.Token{Type: lexer.Number, Literal: "FLFT", Line: 0, Column: 0}
	value, _ := parseNumber(state, token)
	assert.Equal(t, 2, value)

	token = lexer.Token{Type: lexer.Number, Literal: "FLLLLT", Line: 0, Column: 0}
	value, _ = parseNumber(state, token)
	assert.Equal(t, 15, value)

	token = lexer.Token{Type: lexer.Number, Literal: "LLFT", Line: 0, Column: 0}
	value, _ = parseNumber(state, token)
	assert.Equal(t, -2, value)
}

func TestParseInvalidNumber(t *testing.T) {
	state := parseState{
		filename:     "",
		instructions: []executor.Instruction{},
		labelMap:     map[string]int{},
	}

	token := lexer.Token{Type: lexer.Number, Literal: "T", Line: 0, Column: 0}
	_, err := parseNumber(state, token)
	assert.NotNil(t, err)

	token = lexer.Token{Type: lexer.Number, Literal: "F", Line: 0, Column: 0}
	_, err = parseNumber(state, token)
	assert.NotNil(t, err)

	token = lexer.Token{Type: lexer.Number, Literal: "FT", Line: 0, Column: 0}
	_, err = parseNumber(state, token)
	assert.NotNil(t, err)

	token = lexer.Token{Type: lexer.Number, Literal: "FF", Line: 0, Column: 0}
	_, err = parseNumber(state, token)
	assert.NotNil(t, err)
}

func TestParseLabel(t *testing.T) {
	state := parseState{
		filename:     "",
		instructions: []executor.Instruction{},
		labelMap:     map[string]int{},
	}

	token := lexer.Token{Type: lexer.Label, Literal: "FLFLT", Line: 0, Column: 0}
	label, _ := parseLabel(state, token)
	assert.Equal(t, "FLFL", label)
}

func TestParseInvalidLabel(t *testing.T) {
	state := parseState{
		filename:     "",
		instructions: []executor.Instruction{},
		labelMap:     map[string]int{},
	}

	token := lexer.Token{Type: lexer.Label, Literal: "F", Line: 0, Column: 0}
	_, err := parseLabel(state, token)
	assert.NotNil(t, err)

	token = lexer.Token{Type: lexer.Label, Literal: "FLFL", Line: 0, Column: 0}
	_, err = parseLabel(state, token)
	assert.NotNil(t, err)
}

func TestParseToken(t *testing.T) {
	tokens := []lexer.Token{
		lexer.Token{Type: lexer.Getc, Literal: "LTLF", Line: 1, Column: 4},
		lexer.Token{Type: lexer.Getn, Literal: "LTLL", Line: 1, Column: 8},
		lexer.Token{Type: lexer.Putc, Literal: "LTFF", Line: 1, Column: 12},
		lexer.Token{Type: lexer.Putn, Literal: "LTFL", Line: 1, Column: 16},
	}

	expectedInstructions := []executor.Instruction{
		executor.Getc{Token: tokens[0]},
		executor.Getn{Token: tokens[1]},
		executor.Putc{Token: tokens[2]},
		executor.Putn{Token: tokens[3]},
	}

	instructions, _, err := ParseAll(tokens, "")

	if err != nil {
		t.Errorf("expected parse token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedInstructions, instructions)
}

func TestParseTokenWithParameter(t *testing.T) {
	tokens := []lexer.Token{
		lexer.Token{Type: lexer.Push, Literal: "FF", Line: 1, Column: 4},
		lexer.Token{Type: lexer.Number, Literal: "FLLLLT", Line: 0, Column: 0},

		lexer.Token{Type: lexer.Copy, Literal: "FLF", Line: 1, Column: 4},
		lexer.Token{Type: lexer.Number, Literal: "FLLLT", Line: 0, Column: 0},

		lexer.Token{Type: lexer.Slide, Literal: "FLT", Line: 1, Column: 4},
		lexer.Token{Type: lexer.Number, Literal: "FLLT", Line: 0, Column: 0},

		lexer.Token{Type: lexer.MarkLabel, Literal: "TFF", Line: 1, Column: 3},
		lexer.Token{Type: lexer.Label, Literal: "LT", Line: 1, Column: 5},

		lexer.Token{Type: lexer.JumpLabel, Literal: "TFT", Line: 1, Column: 13},
		lexer.Token{Type: lexer.Label, Literal: "LT", Line: 1, Column: 15},

		lexer.Token{Type: lexer.JumpLabelWhenZero, Literal: "TLF", Line: 1, Column: 18},
		lexer.Token{Type: lexer.Label, Literal: "LT", Line: 1, Column: 20},

		lexer.Token{Type: lexer.JumpLabelWhenNegative, Literal: "TLL", Line: 1, Column: 23},
		lexer.Token{Type: lexer.Label, Literal: "LT", Line: 1, Column: 25},
	}

	expectedInstructions := []executor.Instruction{
		executor.Push{Value: 15},
		executor.Copy{Token: tokens[2], Value: 7},
		executor.Slide{Token: tokens[4], Value: 3},
		executor.MarkLabel{Label: "L"},
		executor.JumpLabel{Token: tokens[8], Label: "L"},
		executor.JumpLabelWhenZero{Token: tokens[10], Label: "L"},
		executor.JumpLabelWhenNegative{Token: tokens[12], Label: "L"},
	}

	expectedLabelMap := map[string]int{
		"L": 3,
	}

	instructions, labelMap, err := ParseAll(tokens, "")

	if err != nil {
		t.Errorf("expected parse token with parameter, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedInstructions, instructions)
	assert.Equal(t, expectedLabelMap, labelMap)
}

func TestParseInvalidParameterToknes(t *testing.T) {
	tokens := []lexer.Token{
		lexer.Token{Type: lexer.Push, Literal: "FF", Line: 1, Column: 4},
	}

	_, _, err := ParseAll(tokens, "")

	assert.NotNil(t, err)

	tokens = []lexer.Token{
		lexer.Token{Type: lexer.Push, Literal: "FF", Line: 1, Column: 4},
		lexer.Token{Type: lexer.Label, Literal: "FT", Line: 1, Column: 4},
	}

	_, _, err = ParseAll(tokens, "")

	assert.NotNil(t, err)

	tokens = []lexer.Token{
		lexer.Token{Type: lexer.MarkLabel, Literal: "TFF", Line: 1, Column: 4},
		lexer.Token{Type: lexer.Number, Literal: "FFT", Line: 1, Column: 4},
	}

	_, _, err = ParseAll(tokens, "")

	assert.NotNil(t, err)
}

func TestParseWhenExistsCallLabelBeforeMarkLabel(t *testing.T) {
	tokens := []lexer.Token{
		lexer.Token{Type: lexer.JumpLabel, Literal: "TFT", Line: 1, Column: 13},
		lexer.Token{Type: lexer.Label, Literal: "FT", Line: 1, Column: 1},

		lexer.Token{Type: lexer.MarkLabel, Literal: "TFF", Line: 1, Column: 1},
		lexer.Token{Type: lexer.Label, Literal: "FT", Line: 1, Column: 1},
	}

	expectedInstructions := []executor.Instruction{
		executor.JumpLabel{Token: tokens[0], Label: "F"},
		executor.MarkLabel{Label: "F"},
	}

	expectedLabelMap := map[string]int{
		"F": 1,
	}

	instructions, labelMap, err := ParseAll(tokens, "")

	if err != nil {
		t.Errorf("expected parse token with label, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedInstructions, instructions)
	assert.Equal(t, expectedLabelMap, labelMap)
}
