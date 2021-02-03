package task

import (
	"fmt"
	"time"
)

type Producer struct {
	Count int
}

func (p *Producer) Run(in chan interface{}, out chan interface{}) error {
	for i := 0; i < p.Count; i++ {
		out <- i
	}
	return nil
}

type Square struct {
}

func (square *Square) Run(in chan interface{}, out chan interface{}) error {
	for d := range in {
		out <- d.(int) * d.(int)
	}
	return nil
}

type Add struct {
	AddNum int
}

func (add *Add) Run(in chan interface{}, out chan interface{}) error {
	for d := range in {
		out <- d.(int) + add.AddNum
	}
	return nil
}

type Wait struct {
}

func (wait *Wait) Run(in chan interface{}, out chan interface{}) error {
	for d := range in {
		time.Sleep(1 * time.Second)
		out <- d.(int)
	}
	return nil
}

type Outer struct {
}

func (o *Outer) Run(in chan interface{}, out chan interface{}) error {
	for d := range in {
		fmt.Printf("out: %v\n", d)
	}
	return nil
}
