package executor

import (
	"fmt"
	"strconv"

	"github.com/simomu-github/fflt_lang/lexer"
)

type Instruction interface {
	Execute(executor *Executor) error
	Disassenble() string
}

type Push struct {
	Value int
}

func (p Push) Execute(executor *Executor) error {
	executor.stack = append(executor.stack, p.Value)

	return nil
}

func (p Push) Disassenble() string {
	return fmt.Sprintf("PUSH           %d", p.Value)
}

type Swap struct {
	Token lexer.Token
}

func (s Swap) Execute(executor *Executor) error {
	a, errA := executor.Pop()
	if errA != nil {
		return runtimeErrorWithToken(executor, s.Token, "stack is empty")
	}

	b, errB := executor.Pop()
	if errB != nil {
		return runtimeErrorWithToken(executor, s.Token, "stack is empty")
	}

	executor.Push(a)
	executor.Push(b)
	return nil
}

func (s Swap) Disassenble() string {
	return "SWAP"
}

type Duplicate struct {
	Token lexer.Token
}

func (d Duplicate) Execute(executor *Executor) error {
	a, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, d.Token, "stack is empty")
	}

	executor.Push(a)
	executor.Push(a)
	return nil
}

func (d Duplicate) Disassenble() string {
	return "DUP"
}

type Discard struct{}

func (d Discard) Execute(executor *Executor) error {
	executor.Pop()

	return nil
}

func (d Discard) Disassenble() string {
	return "DISCARD"
}

type Copy struct {
	Token lexer.Token
	Value int
}

func (c Copy) Execute(executor *Executor) error {
	if c.Value < 0 {
		return runtimeErrorWithToken(executor, c.Token, "Copy parameter must be a positive number")
	}

	if len(executor.stack) <= c.Value {
		return runtimeErrorWithToken(
			executor,
			c.Token,
			fmt.Sprintf("copy stack[%d] is out of index. stack length: %d", c.Value, len(executor.stack)),
		)
	}

	value := executor.stack[len(executor.stack)-1-c.Value]
	executor.Push(value)
	return nil
}

func (c Copy) Disassenble() string {
	return fmt.Sprintf("COPY           %d", c.Value)
}

type Slide struct {
	Token lexer.Token
	Value int
}

func (s Slide) Execute(executor *Executor) error {
	if s.Value < 0 {
		return runtimeErrorWithToken(executor, s.Token, "Slide parameter must be a positive number")
	}

	if len(executor.stack) <= s.Value {
		return runtimeErrorWithToken(
			executor,
			s.Token,
			fmt.Sprintf("slide length (%d) is out of stack length (%d)", s.Value, len(executor.stack)),
		)
	}

	newStack := executor.stack[:len(executor.stack)-1-s.Value]
	newStack = append(newStack, executor.stack[len(executor.stack)-1])
	executor.stack = newStack

	return nil
}

func (s Slide) Disassenble() string {
	return fmt.Sprintf("SLIDE          %d", s.Value)
}

type Addition struct {
	Token lexer.Token
}

func (a Addition) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return runtimeErrorWithToken(executor, a.Token, "stack is empty")
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return runtimeErrorWithToken(executor, a.Token, "stack is empty")
	}

	executor.Push(lhs + rhs)

	return nil
}

func (a Addition) Disassenble() string {
	return "ADD"
}

type Subtraction struct {
	Token lexer.Token
}

func (s Subtraction) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return runtimeErrorWithToken(executor, s.Token, "stack is empty")
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return runtimeErrorWithToken(executor, s.Token, "stack is empty")
	}

	executor.Push(lhs - rhs)

	return nil
}

func (s Subtraction) Disassenble() string {
	return "SUB"
}

type Multiplication struct {
	Token lexer.Token
}

func (m Multiplication) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return runtimeErrorWithToken(executor, m.Token, "stack is empty")
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return runtimeErrorWithToken(executor, m.Token, "stack is empty")
	}

	executor.Push(lhs * rhs)

	return nil
}

func (m Multiplication) Disassenble() string {
	return "MUL"
}

type Division struct {
	Token lexer.Token
}

func (d Division) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return runtimeErrorWithToken(executor, d.Token, "stack is empty")
	}

	if rhs == 0 {
		return runtimeErrorWithToken(executor, d.Token, "integer divide by zero")
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return runtimeErrorWithToken(executor, d.Token, "stack is empty")
	}

	executor.Push(lhs / rhs)

	return nil
}

func (d Division) Disassenble() string {
	return "DIV"
}

type Modulo struct {
	Token lexer.Token
}

func (m Modulo) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return runtimeErrorWithToken(executor, m.Token, "stack is empty")
	}

	if rhs == 0 {
		return runtimeErrorWithToken(executor, m.Token, "integer divide by zero")
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return runtimeErrorWithToken(executor, m.Token, "stack is empty")
	}

	executor.Push(lhs % rhs)

	return nil
}

func (m Modulo) Disassenble() string {
	return "MOD"
}

type Getc struct {
	Token lexer.Token
}

func (g Getc) Execute(executor *Executor) error {
	text := executor.Input()

	if len(text) == 0 {
		return runtimeError(executor, "input is empty")
	}

	address, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, g.Token, "stack is empty")
	}

	executor.heap[address] = int([]rune(text)[0])
	return nil
}

func (g Getc) Disassenble() string {
	return "GETC"
}

type Getn struct {
	Token lexer.Token
}

func (g Getn) Execute(executor *Executor) error {
	text := executor.Input()
	n, err := strconv.Atoi(text)
	if err != nil {
		return runtimeError(executor, "input character is not numeric")
	}

	address, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, g.Token, "stack is empty")
	}

	executor.heap[address] = n
	return nil
}

func (g Getn) Disassenble() string {
	return "GETN"
}

type Putc struct {
	Token lexer.Token
}

func (p Putc) Execute(executor *Executor) error {
	n, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, p.Token, "stack is empty")
	}

	executor.Output(fmt.Sprintf("%c", n))

	return nil
}

func (p Putc) Disassenble() string {
	return "PUTC"
}

type Putn struct {
	Token lexer.Token
}

func (p Putn) Execute(executor *Executor) error {
	n, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, p.Token, "stack is empty")
	}

	executor.Output(fmt.Sprintf("%d", n))

	return nil
}

func (p Putn) Disassenble() string {
	return "PUTN"
}

type Store struct {
	Token lexer.Token
}

func (s Store) Execute(executor *Executor) error {
	value, errValue := executor.Pop()
	if errValue != nil {
		return runtimeErrorWithToken(executor, s.Token, "stack is empty")
	}

	address, errAddress := executor.Pop()
	if errAddress != nil {
		return runtimeErrorWithToken(executor, s.Token, "stack is empty")
	}

	executor.heap[address] = value
	return nil
}

func (s Store) Disassenble() string {
	return "STORE"
}

type Retrieve struct {
	Token lexer.Token
}

func (r Retrieve) Execute(executor *Executor) error {
	address, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, r.Token, "stack is empty")
	}

	value, ok := executor.heap[address]
	if !ok {
		return runtimeErrorWithToken(executor, r.Token, "invalid heap access")
	}

	executor.Push(value)
	return nil
}

func (r Retrieve) Disassenble() string {
	return "RETRIEVE"
}

type MarkLabel struct {
	Label string
}

func (m MarkLabel) Execute(executor *Executor) error {
	return nil
}

func (m MarkLabel) Disassenble() string {
	return "LABEL          " + m.Label
}

type CallSubroutine struct {
	Token lexer.Token
	Label string
}

func (c CallSubroutine) Execute(executor *Executor) error {
	counter := executor.programCounter
	executor.PushCallStack(counter)

	newCounter, ok := executor.LabelMap[c.Label]
	if !ok {
		return runtimeErrorWithToken(executor, c.Token, fmt.Sprintf("label \"%s\" is not found", c.Label))
	}

	executor.programCounter = newCounter
	return nil
}

func (c CallSubroutine) Disassenble() string {
	return "CALLSUB        " + c.Label
}

type EndSubroutine struct {
	Token lexer.Token
}

func (e EndSubroutine) Execute(executor *Executor) error {
	counter, err := executor.PopCallStack()
	if err != nil {
		return runtimeErrorWithToken(executor, e.Token, "call stack is empty")
	}

	executor.programCounter = counter
	return nil
}

func (e EndSubroutine) Disassenble() string {
	return "ENDSUB"
}

type JumpLabel struct {
	Token lexer.Token
	Label string
}

func (j JumpLabel) Execute(executor *Executor) error {
	newCounter, ok := executor.LabelMap[j.Label]
	if !ok {
		return runtimeErrorWithToken(executor, j.Token, fmt.Sprintf("label \"%s\" is not found", j.Label))
	}

	executor.programCounter = newCounter
	return nil
}

func (j JumpLabel) Disassenble() string {
	return "JUMP           " + j.Label
}

type JumpLabelWhenZero struct {
	Token lexer.Token
	Label string
}

func (j JumpLabelWhenZero) Execute(executor *Executor) error {
	value, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, j.Token, "stack is empty")
	}

	if value != 0 {
		return nil
	}

	newCounter, ok := executor.LabelMap[j.Label]
	if !ok {
		return runtimeErrorWithToken(executor, j.Token, fmt.Sprintf("label \"%s\" is not found", j.Label))
	}

	executor.programCounter = newCounter
	return nil
}

func (j JumpLabelWhenZero) Disassenble() string {
	return "JUMP_WHEN_ZERO " + j.Label
}

type JumpLabelWhenNegative struct {
	Token lexer.Token
	Label string
}

func (j JumpLabelWhenNegative) Execute(executor *Executor) error {
	value, err := executor.Pop()
	if err != nil {
		return runtimeErrorWithToken(executor, j.Token, "stack is empty")
	}

	if value >= 0 {
		return nil
	}

	newCounter, ok := executor.LabelMap[j.Label]
	if !ok {
		return runtimeErrorWithToken(executor, j.Token, fmt.Sprintf("label \"%s\" is not found", j.Label))
	}

	executor.programCounter = newCounter
	return nil
}

func (j JumpLabelWhenNegative) Disassenble() string {
	return "JUMP_WHEN_NEGA " + j.Label
}

type EndProgram struct{}

func (e EndProgram) Execute(executor *Executor) error {
	executor.programCounter = len(executor.Instructions)
	return nil
}

func (e EndProgram) Disassenble() string {
	return "END"
}
