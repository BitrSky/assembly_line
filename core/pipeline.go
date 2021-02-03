package core

import (
	"context"
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
func (pl *Pipeline) Execute(ctx context.Context) error {

	defer close(pl.errMonitor)

	cctx, cancel := context.WithCancel(ctx)

	em := NewErrMonitor("pipelineErrMonitor")
	defer em.Close()
	go em.Monitor(ctx, cancel)

	wg := &sync.WaitGroup{}

	for _, node := range pl.nodes {
		wg.Add(1)
		go func(ctx context.Context, node TaskNode) {
			defer wg.Done()
			if err := node.Run(ctx); err != nil {
				em.Receive(err)
			}
		}(cctx, node)
	}

	wg.Wait()
	return em.Error()
}
