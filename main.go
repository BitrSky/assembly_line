package main

import (
	"context"
	"fmt"
	"time"

	"github.com/BitrSky/assembly_line/core"
	"github.com/BitrSky/assembly_line/task"
)

func main() {
	pipe := core.NewPipeline()
	producer := core.NewTaskRunner(&task.Producer{Count: 10000000}, 1)
	add := core.NewTaskRunner(&task.Add{AddNum: 100}, 4)
	outer := core.NewTaskRunner(&task.Outer{}, 1)

	pipe.AddNodes(producer, add, outer)
	pipe.Link(producer, add)
	pipe.Link(add, outer)
	start := time.Now().Unix()
	pipe.Execute(context.TODO())
	fmt.Println(time.Now().Unix() - start)
}
