package lexer

const (
	UpperF = "F"
	LowerF = "f"
	UpperL = "L"
	LowerL = "l"
	UpperT = "T"
	LowerT = "t"
	LF     = "\n"
)

type TokenType string

const (
	Push      = TokenType("Push")
	Copy      = TokenType("Copy")
	Swap      = TokenType("Swap")
	Duplicate = TokenType("Duplicate")
	Discard   = TokenType("Discard")
	Slide     = TokenType("Slide")

	Addition       = TokenType("Addition")
	Subtraction    = TokenType("Subtraction")
	Multiplication = TokenType("Multiplication")
	Division       = TokenType("Division")
	Modulo         = TokenType("DivisModulo")

	Store    = TokenType("Store")
	Retrieve = TokenType("Retrieve")

	Putc = TokenType("Putc")
	Putn = TokenType("Putn")
	Getc = TokenType("Getc")
	Getn = TokenType("Getn")

	MarkLabel             = TokenType("MarkLabel")
	CallSubroutine        = TokenType("CallSubroutine")
	JumpLabel             = TokenType("JumpLabel")
	JumpLabelWhenZero     = TokenType("JumpLabelWhenZero")
	JumpLabelWhenNegative = TokenType("JumpLabelWhenNegative")
	EndSubroutine         = TokenType("EndSubroutine")
	EndProgram            = TokenType("EndProgram")

	Number = TokenType("Number")
	Label  = TokenType("Label")
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func isAcceptableCharacter(char string) bool {
	return char == UpperF ||
		char == LowerF ||
		char == UpperL ||
		char == LowerL ||
		char == UpperT ||
		char == LowerT
}

func nullString() string {
	return string([]byte{0})
}
