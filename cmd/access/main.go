package main

import (
	"github.com/koko1123/flow-go-1/cmd"
	nodebuilder "github.com/koko1123/flow-go-1/cmd/access/node_builder"
	"github.com/koko1123/flow-go-1/model/flow"
)

func main() {
	builder := nodebuilder.FlowAccessNode(cmd.FlowNode(flow.RoleAccess.String()))

	builder.PrintBuildVersionDetails()

	// parse all the command line args
	if err := builder.ParseFlags(); err != nil {
		builder.Logger.Fatal().Err(err).Send()
	}

	if err := builder.Initialize(); err != nil {
		builder.Logger.Fatal().Err(err).Send()
	}

	node, err := builder.Build()
	if err != nil {
		builder.Logger.Fatal().Err(err).Send()
	}
	node.Run()
}
