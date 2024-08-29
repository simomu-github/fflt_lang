package executor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/peterh/liner"
)

var (
	commands = []string{
		"s", "step",
		"is", "info stack",
		"ih", "info heap",
		"iv", "info vm",
		"ii", "info instructions",
		"b", "break",
		"d", "delete",
		"c", "continue",
	}
)

const (
	HeaderLength         = 64
	DebuggerInterrupting = "INTERRUPTING"
	DebuggerContinuing   = "CONTINUING"
)

type DebuggerState string

type Debugger struct {
	executor               *Executor
	stackTableWriter       *tablewriter.Table
	heapTableWriter        *tablewriter.Table
	labelMapTableWriter    *tablewriter.Table
	callStackTableWriter   *tablewriter.Table
	breakPoints            mapset.Set[int]
	breakPointsTableWriter *tablewriter.Table
	state                  DebuggerState
	stdin                  *liner.State
	stdout                 string
}

func NewDebugger(executor *Executor) *Debugger {
	stackTableWriter := tablewriter.NewWriter(os.Stdout)
	stackTableWriter.SetHeader([]string{"Value"})

	heapTableWriter := tablewriter.NewWriter(os.Stdout)
	heapTableWriter.SetHeader([]string{"Address", "Value"})

	labelMapTableWriter := tablewriter.NewWriter(os.Stdout)
	labelMapTableWriter.SetHeader([]string{"Label", "Instruction index"})

	callStackTableWriter := tablewriter.NewWriter(os.Stdout)
	callStackTableWriter.SetHeader([]string{"Instruction index"})

	breakPoints := mapset.NewSet[int]()
	breakPoints.Add(0)

	breakPointsTableWriter := tablewriter.NewWriter(os.Stdout)
	breakPointsTableWriter.SetHeader([]string{"Breakpoints"})

	debugger := &Debugger{
		executor:               executor,
		stackTableWriter:       stackTableWriter,
		heapTableWriter:        heapTableWriter,
		labelMapTableWriter:    labelMapTableWriter,
		callStackTableWriter:   callStackTableWriter,
		breakPoints:            breakPoints,
		breakPointsTableWriter: breakPointsTableWriter,
		state:                  DebuggerInterrupting,
		stdin:                  liner.NewLiner(),
	}

	debugger.executor.Input = func() string {
		str, _ := debugger.stdin.Prompt("")
		return str
	}

	debugger.executor.Output = func(value string) {
		if debugger.state == DebuggerInterrupting {
			debugger.stdout += value
		} else {
			fmt.Printf(value)
		}
	}

	return debugger
}

func (d *Debugger) Run() error {
	defer func() {
		d.stdin.Close()
	}()

	d.stdin.SetCtrlCAborts(true)
	d.stdin.SetCompleter(func(line string) (c []string) {
		for _, command := range commands {
			if strings.HasPrefix(command, line) {
				c = append(c, command)
			}
		}
		return
	})

	historyPath := filepath.Join(os.TempDir(), ".fflt_debug_history")

	if f, err := os.Open(historyPath); err == nil {
		d.stdin.ReadHistory(f)
		f.Close()
	}

	d.executor.heap = map[int]int{}
	d.executor.programCounter = 0

	for d.executor.programCounter = 0; d.executor.programCounter < len(d.executor.Instructions); {
		if d.breakPoints.Contains(d.executor.programCounter) {
			d.state = DebuggerInterrupting
		}

		if d.state == DebuggerInterrupting {
			inst := d.executor.Instructions[d.executor.programCounter]
			fmt.Printf("\n")
			fmt.Printf("Program counter: %d\n", d.executor.programCounter)
			fmt.Printf("Current instruction: " + inst.Disassenble() + "\n")
			fmt.Printf("Output: " + d.stdout + "\n")
			d.showStack()

			if err := d.handleCommand(); err != nil {
				return err
			}
		} else {
			err := d.executor.Instructions[d.executor.programCounter].Execute(d.executor)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				d.executor.programCounter++
			}
		}
	}

	if f, err := os.Create(historyPath); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		d.stdin.WriteHistory(f)
		f.Close()
	}

	return nil
}

func (d *Debugger) handleCommand() error {
	for true {
		if command, err := d.stdin.Prompt("> "); err == nil {
			switch command {
			case "s", "step":
				d.stdin.AppendHistory(command)

				err := d.executor.Instructions[d.executor.programCounter].Execute(d.executor)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				} else {
					d.executor.programCounter++
				}
				return nil
			case "is", "info stack":
				d.stdin.AppendHistory(command)
				d.showStack()
			case "ih", "info heap":
				d.stdin.AppendHistory(command)
				d.showHeap()
			case "iv", "info vm":
				d.stdin.AppendHistory(command)
				d.showVM()
			case "ii", "info instructions":
				d.stdin.AppendHistory(command)
				d.showInstructions()
			case "ib", "info breakpoints":
				d.stdin.AppendHistory(command)
				d.showBreakpoints()
			case "c", "continue":
				d.stdin.AppendHistory(command)
				d.state = DebuggerContinuing

				err := d.executor.Instructions[d.executor.programCounter].Execute(d.executor)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				} else {
					d.executor.programCounter++
				}
				return nil
			default:
				split := strings.SplitN(command, " ", 2)
				if len(split) < 2 {
					return nil
				}
				d.handleCommandWithArg(split[0], split[1])
			}
		} else {
			return err
		}
	}

	return nil
}

func (d *Debugger) showVM() {
	fmt.Printf("\n")
	fmt.Printf("Filename: %s\n", d.executor.Filename)
	fmt.Printf("Program counter: %d\n", d.executor.programCounter)
	d.showLabelMap()
	d.showCallStack()
	d.showStack()
	d.showHeap()
}

func (d *Debugger) showStack() {
	d.stackTableWriter.ClearRows()
	for i := len(d.executor.stack) - 1; i >= 0; i-- {
		d.stackTableWriter.Append([]string{fmt.Sprintf("%d", d.executor.stack[i])})
	}

	fmt.Printf("\n")
	outputHeader("Stack")
	d.stackTableWriter.Render()
}

func (d *Debugger) showHeap() {
	d.heapTableWriter.ClearRows()
	for k, v := range d.executor.heap {
		d.heapTableWriter.Append([]string{
			fmt.Sprintf("%d", k),
			fmt.Sprintf("%d", v),
		})
	}

	fmt.Printf("\n")
	outputHeader("Heap")
	d.heapTableWriter.Render()
}

func (d *Debugger) showLabelMap() {
	d.labelMapTableWriter.ClearRows()
	for k, v := range d.executor.LabelMap {
		d.labelMapTableWriter.Append([]string{
			k,
			fmt.Sprintf("%d", v),
		})
	}
	fmt.Printf("\n")
	outputHeader("Label map")
	d.labelMapTableWriter.Render()
}

func (d *Debugger) showCallStack() {
	d.callStackTableWriter.ClearRows()
	for i := len(d.executor.callStack) - 1; i >= 0; i-- {
		d.callStackTableWriter.Append([]string{fmt.Sprintf("%d", d.executor.callStack[i])})
	}
	fmt.Printf("\n")
	outputHeader("Callstack")
	d.callStackTableWriter.Render()
}

func (d *Debugger) showInstructions() {
	fmt.Printf("\n")
	for i, ins := range d.executor.Instructions {
		if d.executor.programCounter == i {
			fmt.Printf("-> %04d %s\n", i, ins.Disassenble())
		} else {
			fmt.Printf("   %04d %s\n", i, ins.Disassenble())
		}
	}
}

func (d *Debugger) showBreakpoints() {
	d.breakPointsTableWriter.ClearRows()
	for b := range d.breakPoints.Iter() {
		d.breakPointsTableWriter.Append([]string{fmt.Sprintf("%d", b)})
	}
	fmt.Printf("\n")
	d.breakPointsTableWriter.Render()
}

func (d *Debugger) handleCommandWithArg(name, arg string) error {
	n, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}

	switch name {
	case "b", "break":
		d.breakPoints.Add(n)
	case "d", "delete":
		d.breakPoints.Remove(n)
	}
	return nil
}

func outputHeader(title string) {
	str := "-- " + title + " "
	remaining := HeaderLength - len(str)
	fmt.Print(str + strings.Repeat("-", remaining) + "\n")
}
