package core

import (
	"context"
	"fmt"
	"sync"
)

const (
	// PipelineMaxNodeNum Pipeline max node num
	PipelineMaxNodeNum int = 20
)

// Pipeline pipeline struct
type Pipeline struct {
	nodes      []TaskNode
	errMonitor chan error
	err        error
	cancelFun  func()
}

// NewPipeline create pipeline instance
func NewPipeline() *Pipeline {
	return &Pipeline{
		nodes:      make([]TaskNode, 0, PipelineMaxNodeNum),
		errMonitor: make(chan error),
	}
}

// Link link two task node
func (pl *Pipeline) Link(node1 TaskNode, node2 TaskNode) {
	c := make(chan interface{})
	node1.AddOutlet(c)
	node2.AddEntrance(c)
}

// AddNodes add task node to pipeline
func (pl *Pipeline) AddNodes(nodes ...TaskNode) {
	pl.nodes = append(pl.nodes, nodes...)
}

// Execute run pipeline
func (pl *Pipeline) Execute(ctx context.Context) {

	defer close(pl.errMonitor)

	cctx, cancel := context.WithCancel(ctx)
	pl.cancelFun = cancel

	wg := &sync.WaitGroup{}
	go pl.errHandler(ctx)

	for _, node := range pl.nodes {
		wg.Add(1)
		go func(ctx context.Context, node TaskNode) {
			defer wg.Done()
			node.Run(ctx, pl.errMonitor)
		}(cctx, node)
	}

	wg.Wait()
}

// Error return running error in pipeline
func (pl *Pipeline) Error() error {
	return pl.err
}

// errHandler handle running error in pipeline
func (pl *Pipeline) errHandler(ctx context.Context) {
	select {
	case err, isOpen := <-pl.errMonitor:
		if isOpen && err != nil {
			pl.err = err
			pl.cancelFun()
		}
	case <-ctx.Done():
		pl.cancelFun()
		fmt.Println(ctx.Err())
	}
}
