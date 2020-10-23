package core

import (
	"context"
	"fmt"
	"sync"
)

const (
	TaskRunning int = 1
	TaskEnd     int = 2
)

// TaskNode task node interface
type TaskNode interface {
	AddEntrance(c chan interface{})
	AddOutlet(c chan interface{})
	Run(ctx context.Context, errMonitor chan error)
}

// Task task interface
type Task interface {
	Run(in chan interface{}, out chan interface{}) error
}

// TaskRunner task runner
type TaskRunner struct {
	entrances []chan interface{}
	outlets   []chan interface{}
	in        chan interface{}
	out       chan interface{}
	task      Task
	status    int
}

// NewTaskRunner create task runner
func NewTaskRunner(task Task) *TaskRunner {
	return &TaskRunner{
		task: task,
	}
}

// AddEntrance add entrance
func (runner *TaskRunner) AddEntrance(entrance chan interface{}) {
	if runner.entrances == nil {
		runner.entrances = make([]chan interface{}, 0, 10)
	}
	runner.entrances = append(runner.entrances, entrance)
}

// AddOutlet add outlet
func (runner *TaskRunner) AddOutlet(outlet chan interface{}) {
	if runner.outlets == nil {
		runner.outlets = make([]chan interface{}, 0, 10)
	}
	runner.outlets = append(runner.outlets, outlet)
}

// before preparatory work
func (runner *TaskRunner) before() {
	runner.status = TaskRunning
	if len(runner.entrances) == 1 {
		runner.in = runner.entrances[0]
	} else if len(runner.entrances) > 1 {
		runner.in = make(chan interface{})
		go runner.collect()
	}

	if len(runner.outlets) == 1 {
		runner.out = runner.outlets[0]
	} else if len(runner.outlets) > 1 {
		runner.out = make(chan interface{})
		go runner.distribution()
	}

}

// end tail-in work
func (runner *TaskRunner) end() {
	runner.status = TaskEnd
	if runner.out != nil {
		close(runner.out)
	}
}

// collect collect multi entrance data
func (runner *TaskRunner) collect() {
	defer close(runner.in)
	wg := &sync.WaitGroup{}
	for _, c := range runner.entrances {
		wg.Add(1)
		go func(c chan interface{}) {
			defer wg.Done()
			for d := range c {
				runner.in <- d
			}
		}(c)
	}
	wg.Wait()
}

// distribution distribution data to multi outlet
func (runner *TaskRunner) distribution() {
	for d := range runner.out {
		for _, c := range runner.outlets {
			c <- d
		}
	}
	for _, c := range runner.outlets {
		close(c)
	}
}

// Run tasn runner run
func (runner *TaskRunner) Run(ctx context.Context, errMonitor chan error) {
	defer runner.end()
	runner.before()

	errCh := make(chan error)
	go runner.runTask(errCh)

	select {
	case err, isOpen := <-errCh:
		if isOpen && err != nil {
			errMonitor <- err
		}
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

// runTask run task in runner
func (runner *TaskRunner) runTask(errMonitor chan error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover; error:%v\n", err)
			errMonitor <- fmt.Errorf("%v", err)
		}
	}()

	errMonitor <- runner.task.Run(runner.in, runner.out)
}
