# 流水线


```golang
  // 创建流水线
  pipe := core.NewPipeline()
  
  // 创建taskrunner，每个taskrunner中封装用户自定义的task
  // 并将taskrunner 添加到pipeline中
	producer := core.NewTaskRunner(&task.Producer{Count: 10000000})
	add1 := core.NewTaskRunner(&task.AddTask{AddNum: 100})
	add2 := core.NewTaskRunner(&task.AddTask{AddNum: 200})
	outer := core.NewTaskRunner(&task.Outer{})
  pipe.AddNodes(producer, add1, add2, outer)

  // 建立task runner之间的联系
	pipe.Link(producer, add1)
	pipe.Link(producer, add2)
	pipe.Link(add1, outer)
  pipe.Link(add2, outer)

  // 运行 pipeline
  pipe.Execute(context.TODO())
  
  // 如果运行时出现问题，导致整个pipeline结束，
  // 可以通过pipe.Error()获取运行时错误
  err := pipe.Error()
```