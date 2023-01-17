package main

import (
	"github.com/koko1123/flow-go-1/cmd"
	insecmd "github.com/koko1123/flow-go-1/insecure/cmd"
	"github.com/koko1123/flow-go-1/model/flow"
)

func main() {
	corruptedBuilder := insecmd.NewCorruptedNodeBuilder(flow.RoleVerification.String())
	corruptedVerificationBuilder := cmd.NewVerificationNodeBuilder(corruptedBuilder.FlowNodeBuilder)
	corruptedVerificationBuilder.LoadFlags()

	corruptedBuilder.LoadCorruptFlags()

	if err := corruptedBuilder.Initialize(); err != nil {
		corruptedVerificationBuilder.Logger.Fatal().Err(err).Send()
	}

	corruptedVerificationBuilder.LoadComponentsAndModules()

	node, err := corruptedVerificationBuilder.FlowNodeBuilder.Build()
	if err != nil {
		corruptedVerificationBuilder.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
