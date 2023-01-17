package main

import (
	"github.com/koko1123/flow-go-1/cmd"
	"github.com/koko1123/flow-go-1/model/flow"
)

func main() {
	verificationBuilder := cmd.NewVerificationNodeBuilder(
		cmd.FlowNode(flow.RoleVerification.String()))
	verificationBuilder.LoadFlags()

	if err := verificationBuilder.FlowNodeBuilder.Initialize(); err != nil {
		verificationBuilder.FlowNodeBuilder.Logger.Fatal().Err(err).Send()
	}

	verificationBuilder.LoadComponentsAndModules()

	node, err := verificationBuilder.FlowNodeBuilder.Build()
	if err != nil {
		verificationBuilder.FlowNodeBuilder.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
