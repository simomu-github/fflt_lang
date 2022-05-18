package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanStackManipulationToken(t *testing.T) {
	source := "FFFFFT" + "FTF" + "FTL" + "FTT" + "FLFFLT" + "FLTFLT"

	expectedTokens := []Token{
		Token{Type: Push, Literal: "FF", Line: 1, Column: 2},
		Token{Type: Number, Literal: "FFFT", Line: 1, Column: 6},

		Token{Type: Duplicate, Literal: "FTF", Line: 1, Column: 9},

		Token{Type: Swap, Literal: "FTL", Line: 1, Column: 12},

		Token{Type: Discard, Literal: "FTT", Line: 1, Column: 15},

		Token{Type: Copy, Literal: "FLF", Line: 1, Column: 18},
		Token{Type: Number, Literal: "FLT", Line: 1, Column: 21},

		Token{Type: Slide, Literal: "FLT", Line: 1, Column: 24},
		Token{Type: Number, Literal: "FLT", Line: 1, Column: 27},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan stack manipulate token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanArithmeticToken(t *testing.T) {
	source := "LFFF" + "LFFL" + "LFFT" + "LFLF" + "LFLL"

	expectedTokens := []Token{
		Token{Type: Addition, Literal: "LFFF", Line: 1, Column: 4},
		Token{Type: Subtraction, Literal: "LFFL", Line: 1, Column: 8},
		Token{Type: Multiplication, Literal: "LFFT", Line: 1, Column: 12},
		Token{Type: Division, Literal: "LFLF", Line: 1, Column: 16},
		Token{Type: Modulo, Literal: "LFLL", Line: 1, Column: 20},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan arithmetic token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanHeapAccessToken(t *testing.T) {
	source := "LLF" + "LLL"

	expectedTokens := []Token{
		Token{Type: Store, Literal: "LLF", Line: 1, Column: 3},
		Token{Type: Retrieve, Literal: "LLL", Line: 1, Column: 6},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan heap access token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanIOToken(t *testing.T) {
	source := "LTLF" + "LTLL" + "LTFF" + "LTFL"

	expectedTokens := []Token{
		Token{Type: Getc, Literal: "LTLF", Line: 1, Column: 4},
		Token{Type: Getn, Literal: "LTLL", Line: 1, Column: 8},
		Token{Type: Putc, Literal: "LTFF", Line: 1, Column: 12},
		Token{Type: Putn, Literal: "LTFL", Line: 1, Column: 16},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan IO token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanFlowControllToken(t *testing.T) {
	source := "TFFLT" + "TFLLT" + "TFTLT" + "TLFLT" + "TLLLT" + "TLT" + "TTT"

	expectedTokens := []Token{
		Token{Type: MarkLabel, Literal: "TFF", Line: 1, Column: 3},
		Token{Type: Label, Literal: "LT", Line: 1, Column: 5},

		Token{Type: CallSubroutine, Literal: "TFL", Line: 1, Column: 8},
		Token{Type: Label, Literal: "LT", Line: 1, Column: 10},

		Token{Type: JumpLabel, Literal: "TFT", Line: 1, Column: 13},
		Token{Type: Label, Literal: "LT", Line: 1, Column: 15},

		Token{Type: JumpLabelWhenZero, Literal: "TLF", Line: 1, Column: 18},
		Token{Type: Label, Literal: "LT", Line: 1, Column: 20},

		Token{Type: JumpLabelWhenNegative, Literal: "TLL", Line: 1, Column: 23},
		Token{Type: Label, Literal: "LT", Line: 1, Column: 25},

		Token{Type: EndSubroutine, Literal: "TLT", Line: 1, Column: 28},

		Token{Type: EndProgram, Literal: "TTT", Line: 1, Column: 31},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan flow controll token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanTokenWithOtherCharacter(t *testing.T) {
	source := "hogehogeaaaa\n\n\nFF  F F\nFTaaaaa"

	expectedTokens := []Token{
		Token{Type: Push, Literal: "FF", Line: 4, Column: 2},
		Token{Type: Number, Literal: "FFFT", Line: 5, Column: 2},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanTokenWithDownCaseCharacter(t *testing.T) {
	source := "FfFfft"

	expectedTokens := []Token{
		Token{Type: Push, Literal: "Ff", Line: 1, Column: 2},
		Token{Type: Number, Literal: "Ffft", Line: 1, Column: 6},
	}

	tokens, err := ScanAllTokens(source, "")

	if err != nil {
		t.Errorf("expected scan token, but raise error %s", err.Error())
	}

	assert.Equal(t, expectedTokens, tokens)
}

func TestScanInvalidToken(t *testing.T) {
	source := "Faaa"

	_, err := ScanAllTokens(source, "")

	assert.NotNil(t, err)

	source = "LLhoge"

	_, err = ScanAllTokens(source, "")

	assert.NotNil(t, err)
}
