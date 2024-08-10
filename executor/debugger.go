package executor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
	"github.com/peterh/liner"
)

type Debugger struct {
	executor         *Executor
	stackTableWriter *tablewriter.Table
	heapTableWriter  *tablewriter.Table
	stdout           string
}

func NewDebugger(executor *Executor) *Debugger {
	stackTableWriter := tablewriter.NewWriter(os.Stdout)
	stackTableWriter.SetHeader([]string{"Value"})

	heapTableWriter := tablewriter.NewWriter(os.Stdout)
	heapTableWriter.SetHeader([]string{"Address", "Value"})

	debugger := &Debugger{
		executor:         executor,
		stackTableWriter: stackTableWriter,
		heapTableWriter:  heapTableWriter,
	}

	executor.Output = func(value string) {
		debugger.stdout += value
	}

	return debugger
}

func (d *Debugger) Run() error {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	historyPath := filepath.Join(os.TempDir(), ".fflt_debug_history")

	if f, err := os.Open(historyPath); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	d.executor.heap = map[int]int{}
	d.executor.programCounter = 0

	for d.executor.programCounter = 0; d.executor.programCounter < len(d.executor.Instructions); {
		inst := d.executor.Instructions[d.executor.programCounter]
		fmt.Printf("\n")
		fmt.Printf("Current instruction: " + inst.Disassenble() + "\n")
		fmt.Printf("Output: " + d.stdout + "\n")
		d.showStack()

		if err := d.handleCommand(line); err != nil {
			return err
		}
	}

	if f, err := os.Create(historyPath); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	return nil
}

func (d *Debugger) handleCommand(line *liner.State) error {
	for true {
		if command, err := line.Prompt("> "); err == nil {
			switch command {
			case "s", "step":
				err := d.executor.Instructions[d.executor.programCounter].Execute(d.executor)
				if err != nil {
					return err
				}
				d.executor.programCounter++
				return nil
			case "is", "info stack":
				d.showStack()
			case "ih", "info heap":
				d.showHeap()
			}
		} else {
			return err
		}
	}

	return nil
}

func (d *Debugger) showStack() {
	d.stackTableWriter.ClearRows()
	for i := len(d.executor.stack) - 1; i >= 0; i-- {
		d.stackTableWriter.Append([]string{fmt.Sprintf("%d", d.executor.stack[i])})
	}

	fmt.Printf("\n")
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
	d.heapTableWriter.Render()
}
