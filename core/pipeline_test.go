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
	producer := NewTaskRunner(&task.Producer{Count: 20})
	add1 := NewTaskRunner(&task.AddTask{AddNum: 100})
	add2 := NewTaskRunner(&taskWithError{})
	outer := NewTaskRunner(&task.Outer{})
	pipe.AddNodes(producer, add1, add2, outer)

	pipe.Link(producer, add1)
	pipe.Link(producer, add2)
	pipe.Link(add1, outer)
	pipe.Link(add2, outer)
	pipe.Execute(context.TODO())

	assert.Equal(s.T(), Err, pipe.Error())
}

func (s PipelineSuite) TestExecute() {
	pipe := NewPipeline()

	producer := NewTaskRunner(&task.Producer{Count: 20})
	add1 := NewTaskRunner(&task.AddTask{AddNum: 100})
	add2 := NewTaskRunner(&task.AddTask{AddNum: 200})
	outer := NewTaskRunner(&task.Outer{})
	pipe.AddNodes(producer, add1, add2, outer)

	pipe.Link(producer, add1)
	pipe.Link(producer, add2)
	pipe.Link(add1, outer)
	pipe.Link(add2, outer)
	pipe.Execute(context.TODO())

	assert.Equal(s.T(), nil, pipe.Error())
}

func TestUnitPipeline(t *testing.T) {
	s := new(PipelineSuite)
	suite.Run(t, s)

}
