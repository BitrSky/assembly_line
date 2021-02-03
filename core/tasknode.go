package core

import (
	"context"
	"fmt"
	"reflect"
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
	Run(ctx context.Context) error
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
	parallel  int
	status    int
}

// NewTaskRunner create task runner
func NewTaskRunner(task Task, parallel int) *TaskRunner {
	return &TaskRunner{
		task:     task,
		parallel: parallel,
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
func (runner *TaskRunner) Run(ctx context.Context) error {
	defer runner.end()
	runner.before()

	cctx, cancel := context.WithCancel(ctx)
	em := NewErrMonitor(fmt.Sprintf("taskErrMonotor_%s", reflect.TypeOf(runner.task).String()))
	defer em.Close()
	go em.Monitor(ctx, cancel)

	wg := &sync.WaitGroup{}
	for i := 0; i < runner.parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runner.runTask(cctx); err != nil {
				em.Receive(err)
			}
		}()
	}

	wg.Wait()

	return em.Error()
}

// runTask run task in runner
func (runner *TaskRunner) runTask(ctx context.Context) (err error) {

	errc := make(chan error)
	go func() {
		defer func() {
			if recover := recover(); recover != nil {
				err = fmt.Errorf("panic recover :%v", recover)
			}
		}()
		errc <- runner.task.Run(runner.in, runner.out)
	}()

	select {
	case errv, ok := <-errc:
		if !ok {
			break
		}
		err = errv
	case <-ctx.Done():
	}
	return
}
