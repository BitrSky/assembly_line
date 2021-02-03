package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/BitrSky/assembly_line/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var Err = fmt.Errorf("test error")

type taskWithError struct {
}

func (t *taskWithError) Run(in chan interface{}, out chan interface{}) error {
	i := 0
	for d := range in {
		out <- d.(int) + 1000
		i++
		if i == 10 {
			return Err
		}

	}
	return nil
}

type PipelineSuite struct {
	suite.Suite
}

func (s *PipelineSuite) SetupSuite() {

}

func (s PipelineSuite) TestExecuteWithError() {
	pipe := NewPipeline()
	producer := NewTaskRunner(&task.Producer{Count: 20}, 1)
	add1 := NewTaskRunner(&task.Add{AddNum: 100}, 2)
	add2 := NewTaskRunner(&taskWithError{}, 1)
	outer := NewTaskRunner(&task.Outer{}, 1)
	pipe.AddNodes(producer, add1, add2, outer)

	pipe.Link(producer, add1)
	pipe.Link(producer, add2)
	pipe.Link(add1, outer)
	pipe.Link(add2, outer)
	err := pipe.Execute(context.TODO())

	assert.Equal(s.T(), Err, err)
}

func (s PipelineSuite) TestExecute() {
	pipe := NewPipeline()

	producer := NewTaskRunner(&task.Producer{Count: 20}, 1)
	add1 := NewTaskRunner(&task.Add{AddNum: 100}, 2)
	add2 := NewTaskRunner(&task.Add{AddNum: 200}, 1)
	outer := NewTaskRunner(&task.Outer{}, 1)
	pipe.AddNodes(producer, add1, add2, outer)

	pipe.Link(producer, add1)
	pipe.Link(producer, add2)
	pipe.Link(add1, outer)
	pipe.Link(add2, outer)
	err := pipe.Execute(context.TODO())

	assert.Equal(s.T(), nil, err)
}

func (s PipelineSuite) TestExecute2() {
	pipe := NewPipeline()

	producer := NewTaskRunner(&task.Producer{Count: 10}, 1)
	wait := NewTaskRunner(&task.Wait{}, 5)
	outer := NewTaskRunner(&task.Outer{}, 1)
	pipe.AddNodes(producer, wait, outer)

	pipe.Link(producer, wait)
	pipe.Link(wait, outer)
	err := pipe.Execute(context.TODO())
	//fmt.Println(err)

	assert.Equal(s.T(), nil, err)
}

func TestUnitPipeline(t *testing.T) {
	s := new(PipelineSuite)
	suite.Run(t, s)
}
