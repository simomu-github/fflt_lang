package lexer

type Lexer struct {
	filename      string
	source        string
	currentIndex  int
	newLine       bool
	currentColumn int
	currentLine   int
	currentToken  string
	tokens        []Token
}

func ScanAllTokens(source string, filename string) ([]Token, error) {
	lexer := &Lexer{
		filename:      filename,
		source:        source,
		currentIndex:  -1,
		currentLine:   1,
		currentColumn: 0,
		newLine:       false,
	}

	var allTokens []Token

	for lexer.currentIndex < len(lexer.source) {
		tokens, err := lexer.scanToken()
		if err != nil {
			return []Token{}, err
		}

		allTokens = append(allTokens, tokens...)
	}

	return allTokens, nil
}

func (lexer *Lexer) scanToken() ([]Token, error) {
	lexer.currentToken = ""
	char := lexer.readNextChar()
	switch char {
	case UpperF, LowerF:
		return lexer.scanStackManipulation()
	case UpperL, LowerL:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			return lexer.scanArithmetic()
		case UpperL, LowerL:
			return lexer.scanHeapAccess()
		case UpperT, LowerT:
			return lexer.scanIO()
		}
	case UpperT, LowerT:
		return lexer.scanFlowControll()
	}

	return lexer.tokens, nil
}

func (lexer *Lexer) scanStackManipulation() ([]Token, error) {
	char := lexer.readNextChar()
	switch char {
	case UpperF, LowerF:
		pushToken := Token{Type: Push, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
		numberToken := lexer.scanValue(Number)
		return []Token{pushToken, numberToken}, nil
	case UpperL, LowerL:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			copyToken := Token{Type: Copy, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			numberToken := lexer.scanValue(Number)
			return []Token{copyToken, numberToken}, nil
		case UpperT, LowerT:
			slideToken := Token{Type: Slide, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			numberToken := lexer.scanValue(Number)
			return []Token{slideToken, numberToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected stack manipulation command")
		}
	case UpperT, LowerT:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			duplicateToken := Token{Type: Duplicate, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{duplicateToken}, nil
		case UpperL, LowerL:
			swapToken := Token{Type: Swap, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{swapToken}, nil
		case UpperT, LowerT:
			discardToken := Token{Type: Discard, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{discardToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected stack manipulation command")
		}
	default:
		return []Token{}, lexicalError(lexer, "expected stack manipulation command")
	}
}

func (lexer *Lexer) scanArithmetic() ([]Token, error) {
	char := lexer.readNextChar()
	switch char {
	case UpperF, LowerF:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			additionToken := Token{Type: Addition, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{additionToken}, nil
		case UpperL, LowerL:
			subtractionToken := Token{Type: Subtraction, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{subtractionToken}, nil
		case UpperT, LowerT:
			multiplicationToken := Token{Type: Multiplication, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{multiplicationToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected artithemetic command")
		}
	case UpperL, LowerL:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			divisionToken := Token{Type: Division, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{divisionToken}, nil
		case UpperL, LowerL:
			moduloToken := Token{Type: Modulo, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{moduloToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected artithemetic command")
		}
	default:
		return []Token{}, lexicalError(lexer, "expected artithemetic command")
	}
}

func (lexer *Lexer) scanHeapAccess() ([]Token, error) {
	char := lexer.readNextChar()
	switch char {
	case UpperF, LowerF:
		storeToken := Token{Type: Store, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
		return []Token{storeToken}, nil
	case UpperL, LowerL:
		retrieveToken := Token{Type: Retrieve, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
		return []Token{retrieveToken}, nil
	default:
		return []Token{}, lexicalError(lexer, "expected heap access command")
	}
}

func (lexer *Lexer) scanIO() ([]Token, error) {
	char := lexer.readNextChar()
	switch char {
	case UpperF, LowerF:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			putcToken := Token{Type: Putc, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{putcToken}, nil
		case UpperL, LowerL:
			putnToken := Token{Type: Putn, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{putnToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected IO command")
		}
	case UpperL, LowerL:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			getcToken := Token{Type: Getc, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{getcToken}, nil
		case UpperL, LowerL:
			getnToken := Token{Type: Getn, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{getnToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected IO command")
		}
	default:
		return []Token{}, lexicalError(lexer, "expected IO command")
	}
}

func (lexer *Lexer) scanFlowControll() ([]Token, error) {
	char := lexer.readNextChar()
	switch char {
	case UpperF, LowerF:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			markLabelToken := Token{Type: MarkLabel, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			labelToken := lexer.scanValue(Label)
			return []Token{markLabelToken, labelToken}, nil
		case UpperL, LowerL:
			callSubroutineToken := Token{Type: CallSubroutine, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			labelToken := lexer.scanValue(Label)
			return []Token{callSubroutineToken, labelToken}, nil
		case UpperT, LowerT:
			jumpLabelToken := Token{Type: JumpLabel, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			labelToken := lexer.scanValue(Label)
			return []Token{jumpLabelToken, labelToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected flow controll command")
		}
	case UpperL, LowerL:
		char = lexer.readNextChar()
		switch char {
		case UpperF, LowerF:
			jumpLabelWhenZeroToken := Token{Type: JumpLabelWhenZero, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			labelToken := lexer.scanValue(Label)
			return []Token{jumpLabelWhenZeroToken, labelToken}, nil
		case UpperL, LowerL:
			jumpLabelWhenNegativeToken := Token{Type: JumpLabelWhenNegative, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			labelToken := lexer.scanValue(Label)
			return []Token{jumpLabelWhenNegativeToken, labelToken}, nil
		case UpperT, LowerT:
			endSubroutineToken := Token{Type: EndSubroutine, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{endSubroutineToken}, nil
		default:
			return []Token{}, lexicalError(lexer, "expected flow controll command")
		}
	case UpperT, LowerT:
		char = lexer.readNextChar()
		if char == UpperT || char == LowerT {
			endProgramToken := Token{Type: EndProgram, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
			return []Token{endProgramToken}, nil
		} else {
			return []Token{}, lexicalError(lexer, "expected flow controll command")
		}
	default:
		return []Token{}, lexicalError(lexer, "expected flow controll command")
	}
}

func (lexer *Lexer) scanValue(tokenType TokenType) Token {
	lexer.currentToken = ""
	var result string
	char := lexer.readNextChar()
	result = char
	for char != UpperT && char != LowerT && char != nullString() {
		char = lexer.readNextChar()
		result += char
	}

	return Token{Type: tokenType, Literal: lexer.currentToken, Line: lexer.currentLine, Column: lexer.currentColumn}
}

func (lexer *Lexer) readNextChar() string {
	lexer.currentIndex++
	lexer.currentColumn++

	if lexer.currentIndex >= len(lexer.source) {
		return nullString()
	}

	if lexer.currentChar() == LF {
		lexer.currentLine++
		lexer.currentColumn = 0
	}

	for !isAcceptableCharacter(lexer.currentChar()) {
		lexer.currentIndex++
		lexer.currentColumn++

		if lexer.currentIndex >= len(lexer.source) {
			return nullString()
		}

		if lexer.currentChar() == LF {
			lexer.currentLine++
			lexer.currentColumn = 0
		}
	}

	lexer.currentToken += lexer.currentChar()
	return lexer.currentChar()
}

func (lexer Lexer) currentChar() string {
	return string(lexer.source[lexer.currentIndex])
}

func (lexer *Lexer) addToken(token Token) {
	lexer.tokens = append(lexer.tokens, token)
}
