package main

import (
	"github.com/koko1123/flow-go-1/cmd"
	"github.com/koko1123/flow-go-1/model/flow"
)

func main() {
	exeBuilder := cmd.NewExecutionNodeBuilder(cmd.FlowNode(flow.RoleExecution.String()))
	exeBuilder.LoadFlags()

	if err := exeBuilder.FlowNodeBuilder.Initialize(); err != nil {
		exeBuilder.FlowNodeBuilder.Logger.Fatal().Err(err).Send()
	}

	exeBuilder.LoadComponentsAndModules()

	node, err := exeBuilder.FlowNodeBuilder.Build()
	if err != nil {
		exeBuilder.FlowNodeBuilder.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
