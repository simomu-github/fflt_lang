package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newExecutor() *Executor {
	return &Executor{
		heap:           map[int]int{},
		programCounter: 0,
	}
}

func newExecutorWithIOMock(input func() string, output func(string)) *Executor {
	return &Executor{
		Input:          input,
		Output:         output,
		heap:           map[int]int{},
		programCounter: 0,
	}
}

type ioMock struct {
	In  func() string
	Out func(string)
}

func TestPush(t *testing.T) {
	executor := newExecutor()

	push := Push{Value: 1}
	push.Execute(executor)

	assert.Equal(t, []int{1}, executor.stack)
}

func TestDuplicate(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1}

	duplicate := Duplicate{}
	duplicate.Execute(executor)

	assert.Equal(t, []int{1, 1}, executor.stack)
}

func TestSwap(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2}

	swap := Swap{}
	swap.Execute(executor)

	assert.Equal(t, []int{2, 1}, executor.stack)
}

func TestDiscard(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1}

	discard := Discard{}
	discard.Execute(executor)

	assert.Equal(t, []int{}, executor.stack)
}

func TestCopy(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2, 3}

	copy := Copy{Value: 1}
	copy.Execute(executor)

	assert.Equal(t, []int{1, 2, 3, 2}, executor.stack)
}

func TestSlide(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2, 3, 4}

	slide := Slide{Value: 2}
	slide.Execute(executor)

	assert.Equal(t, []int{1, 4}, executor.stack)
}

func TestAddition(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2}

	addition := Addition{}
	addition.Execute(executor)

	assert.Equal(t, []int{3}, executor.stack)
}

func TestSubtraction(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{2, 1}

	subtraction := Subtraction{}
	subtraction.Execute(executor)

	assert.Equal(t, []int{1}, executor.stack)
}

func TestMultiplication(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{2, 2}

	multiplication := Multiplication{}
	multiplication.Execute(executor)

	assert.Equal(t, []int{4}, executor.stack)
}

func TestDivision(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{4, 2}

	divition := Division{}
	divition.Execute(executor)

	assert.Equal(t, []int{2}, executor.stack)

	t.Run("when divide by zero", func(t *testing.T) {
		executor.stack = []int{4, 0}

		divition := Division{}
		err := divition.Execute(executor)
		assert.NotNil(t, err)
	})
}

func TestModulo(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{5, 3}

	modulo := Modulo{}
	modulo.Execute(executor)

	assert.Equal(t, []int{2}, executor.stack)

	t.Run("when modulo by zero", func(t *testing.T) {
		executor.stack = []int{5, 0}

		divition := Modulo{}
		err := divition.Execute(executor)
		assert.NotNil(t, err)
	})
}

func TestStore(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2}

	store := Store{}
	store.Execute(executor)

	assert.Equal(t, map[int]int{1: 2}, executor.heap)
	assert.Equal(t, []int{}, executor.stack)
}

func TestRetrieve(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1}
	executor.heap = map[int]int{1: 2}

	retrieve := Retrieve{}
	retrieve.Execute(executor)

	assert.Equal(t, []int{2}, executor.stack)
}

func TestPutc(t *testing.T) {
	executor := newExecutorWithIOMock(
		func() string { return "" },
		func(str string) {
			assert.Equal(t, "a", str)
		},
	)
	executor.stack = []int{97}

	putc := Putc{}
	putc.Execute(executor)

	assert.Equal(t, []int{}, executor.stack)
}

func TestPutn(t *testing.T) {
	executor := newExecutorWithIOMock(
		func() string { return "" },
		func(str string) {
			assert.Equal(t, "100", str)
		},
	)
	executor.stack = []int{100}

	putn := Putn{}
	putn.Execute(executor)

	assert.Equal(t, []int{}, executor.stack)
}

func TestGetc(t *testing.T) {
	inputChar := "a"
	executor := newExecutorWithIOMock(
		func() string {
			return inputChar
		},
		func(_ string) {},
	)
	executor.stack = []int{1}

	getc := Getc{}
	getc.Execute(executor)

	assert.Equal(t, []int{}, executor.stack)
	assert.Equal(t, map[int]int{1: 97}, executor.heap)
}

func TestGetn(t *testing.T) {
	inputChar := "100"
	executor := newExecutorWithIOMock(
		func() string {
			return inputChar
		},
		func(_ string) {},
	)
	executor.stack = []int{1}

	getc := Getn{}
	getc.Execute(executor)

	assert.Equal(t, []int{}, executor.stack)
	assert.Equal(t, map[int]int{1: 100}, executor.heap)
}

func TestCallSubroutine(t *testing.T) {
	executor := newExecutor()
	executor.Instructions = append(executor.Instructions, MarkLabel{})
	executor.LabelMap = map[string]int{
		"F": 0,
	}
	executor.programCounter = 1

	callSubroutine := CallSubroutine{Label: "F"}
	callSubroutine.Execute(executor)

	assert.Equal(t, 0, executor.programCounter)
}

func TestEndSubroutine(t *testing.T) {
	executor := newExecutor()
	executor.programCounter = 10
	executor.callStack = append(executor.callStack, 0)

	endSubroutine := EndSubroutine{}
	endSubroutine.Execute(executor)

	assert.Equal(t, 0, executor.programCounter)
}

func TestJumpLabel(t *testing.T) {
	executor := newExecutor()
	executor.Instructions = append(executor.Instructions, MarkLabel{})
	executor.LabelMap = map[string]int{
		"F": 0,
	}
	executor.programCounter = 1

	jumpLabel := JumpLabel{Label: "F"}
	jumpLabel.Execute(executor)

	assert.Equal(t, 0, executor.programCounter)
}

func TestJumplabelWhenZero(t *testing.T) {
	executor := newExecutor()
	executor.Instructions = append(executor.Instructions, MarkLabel{})
	executor.LabelMap = map[string]int{
		"F": 0,
	}
	executor.programCounter = 1
	executor.stack = []int{0}

	jumpLabelWhenZero := JumpLabelWhenZero{Label: "F"}
	jumpLabelWhenZero.Execute(executor)

	assert.Equal(t, 0, executor.programCounter)

	executor.programCounter = 1
	executor.stack = []int{1}
	jumpLabelWhenZero.Execute(executor)

	assert.Equal(t, 1, executor.programCounter)
}

func TestJumplabelWhenNegative(t *testing.T) {
	executor := newExecutor()
	executor.Instructions = append(executor.Instructions, MarkLabel{})
	executor.LabelMap = map[string]int{
		"F": 0,
	}
	executor.programCounter = 1
	executor.stack = []int{-1}

	jumpLabelWhenNegative := JumpLabelWhenNegative{Label: "F"}
	jumpLabelWhenNegative.Execute(executor)

	assert.Equal(t, 0, executor.programCounter)

	executor.programCounter = 1
	executor.stack = []int{0}
	jumpLabelWhenNegative.Execute(executor)

	assert.Equal(t, 1, executor.programCounter)

	executor.programCounter = 1
	executor.stack = []int{1}
	jumpLabelWhenNegative.Execute(executor)

	assert.Equal(t, 1, executor.programCounter)
}

func TestEndProgram(t *testing.T) {
	executor := newExecutor()
	executor.Instructions = append(executor.Instructions, MarkLabel{})
	executor.programCounter = 0

	endProgram := EndProgram{}
	endProgram.Execute(executor)

	assert.Equal(t, 1, executor.programCounter)
}

func TestRuntimeError(t *testing.T) {
	executor := newExecutor()

	t.Run("when stack is empty", func(t *testing.T) {
		addition := Addition{}
		err := addition.Execute(executor)
		assert.NotNil(t, err)
	})

	t.Run("when copy to out of range", func(t *testing.T) {
		executor.stack = []int{1}
		copy := Copy{Value: 1}
		err := copy.Execute(executor)
		assert.NotNil(t, err)
	})

	t.Run("when slide to out of range", func(t *testing.T) {
		executor.stack = []int{1, 2, 3, 4}
		slide := Slide{Value: 4}
		err := slide.Execute(executor)
		assert.NotNil(t, err)
	})

	t.Run("when invalid heap access", func(t *testing.T) {
		executor.stack = []int{1}
		executor.heap = map[int]int{}

		retrieve := Retrieve{}
		err := retrieve.Execute(executor)
		assert.NotNil(t, err)
	})

	t.Run("when label is not found", func(t *testing.T) {
		jumpLabel := JumpLabel{Label: "F"}
		err := jumpLabel.Execute(executor)
		assert.NotNil(t, err)
	})
}
