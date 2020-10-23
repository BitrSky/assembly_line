package task

import (
	"fmt"
)

type Producer struct {
	Count int
}

func (p *Producer) Run(in chan interface{}, out chan interface{}) error {
	for i := 0; i < p.Count; i++ {
		//	fmt.Printf("in:%d\n", i)
		out <- i
	}
	return nil
}

type AddTask struct {
	AddNum int
}

func (add *AddTask) Run(in chan interface{}, out chan interface{}) error {
	for d := range in {
		//		fmt.Printf("add:%d + %d\n", d, add.AddNum)
		out <- d.(int) + add.AddNum
	}
	return nil
}

type Outer struct {
}

func (o *Outer) Run(in chan interface{}, out chan interface{}) error {
	for d := range in {
		_ = fmt.Sprintf("out: %v\n", d)
	}
	return nil
}
