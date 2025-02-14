package main

import (
	"github.com/koko1123/flow-go-1/cmd"
	insecmd "github.com/onflow/flow-go/insecure/cmd"
	"github.com/koko1123/flow-go-1/model/flow"
)

func main() {
	corruptedBuilder := insecmd.NewCorruptedNodeBuilder(flow.RoleExecution.String())
	corruptedExecutionBuilder := cmd.NewExecutionNodeBuilder(corruptedBuilder.FlowNodeBuilder)
	corruptedExecutionBuilder.LoadFlags()

	corruptedBuilder.LoadCorruptFlags()

	if err := corruptedBuilder.Initialize(); err != nil {
		corruptedExecutionBuilder.Logger.Fatal().Err(err).Send()
	}

	corruptedExecutionBuilder.LoadComponentsAndModules()

	node, err := corruptedExecutionBuilder.FlowNodeBuilder.Build()
	if err != nil {
		corruptedExecutionBuilder.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
