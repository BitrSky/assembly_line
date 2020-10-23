package main

import (
	"context"

	"github.com/BitrSky/assembly_line/core"
	"github.com/BitrSky/assembly_line/task"
)

func main() {
	pipe := core.NewPipeline()
	producer := core.NewTaskRunner(&task.Producer{Count: 10000000})
	add1 := core.NewTaskRunner(&task.AddTask{AddNum: 100})
	add2 := core.NewTaskRunner(&task.AddTask{AddNum: 200})
	outer := core.NewTaskRunner(&task.Outer{})
	pipe.AddNodes(producer, add1, add2, outer)

	pipe.Link(producer, add1)
	pipe.Link(producer, add2)
	pipe.Link(add1, outer)
	pipe.Link(add2, outer)
	pipe.Execute(context.TODO())
}
