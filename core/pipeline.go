package core

import (
	"context"
	"sync"
)

// IPipeline Pipeline Interface
type IPipeline interface {
	AddNodes(nodes ...ITaskRunner)
	Link(node1 ITaskRunner, node2 ITaskRunner)
	Execute(ctx context.Context) error
}

const (
	// PipelineMaxNodeNum Pipeline max node num
	PipelineMaxNodeNum int = 20
)

// Pipeline pipeline struct
type Pipeline struct {
	nodes      []ITaskRunner
	errMonitor chan error
}

// NewPipeline create pipeline instance
func NewPipeline() *Pipeline {
	return &Pipeline{
		nodes:      make([]ITaskRunner, 0, PipelineMaxNodeNum),
		errMonitor: make(chan error),
	}
}

// Link link two task node
func (pl *Pipeline) Link(node1 ITaskRunner, node2 ITaskRunner) {
	c := make(chan interface{})
	node1.AddOutlet(c)
	node2.AddEntrance(c)
}

// AddNodes add task node to pipeline
func (pl *Pipeline) AddNodes(nodes ...ITaskRunner) {
	pl.nodes = append(pl.nodes, nodes...)
}

// Execute run pipeline
func (pl *Pipeline) Execute(ctx context.Context) error {

	defer close(pl.errMonitor)

	cctx, cancel := context.WithCancel(ctx)

	em := NewErrMonitor("pipelineErrMonitor")
	defer em.Close()
	go em.Monitor(ctx, cancel)

	wg := &sync.WaitGroup{}

	for _, node := range pl.nodes {
		wg.Add(1)
		go func(ctx context.Context, node ITaskRunner) {
			defer wg.Done()
			if err := node.Run(ctx); err != nil {
				em.Receive(err)
			}
		}(cctx, node)
	}

	wg.Wait()
	return em.Error()
}
