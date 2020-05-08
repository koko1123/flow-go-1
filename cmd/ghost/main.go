package main

import (
	"github.com/spf13/pflag"

	"github.com/dapperlabs/flow-go/module/metrics"

	"github.com/dapperlabs/flow-go/cmd"
	"github.com/dapperlabs/flow-go/engine/ghost/engine"
	"github.com/dapperlabs/flow-go/module"
	"github.com/dapperlabs/flow-go/network/gossip/libp2p/validators"
)

func main() {
	var (
		rpcConf engine.Config
		err     error
	)

	cmd.FlowNode("ghost").
		ExtraFlags(func(flags *pflag.FlagSet) {
			flags.StringVarP(&rpcConf.ListenAddr, "rpc-addr", "r", "localhost:9000", "the address the GRPC server listens on")
		}).
		Module("message validators", func(node *cmd.FlowNodeBuilder) error {
			node.MsgValidators = []validators.MessageValidator{
				// filter out messages sent by this node itself
				validators.NewSenderValidator(node.Me.NodeID()),
				// but retain all the 1-k messages even if they are not intended for this node
			}
			return nil
		}).
		Module("metrics collector", func(node *cmd.FlowNodeBuilder) error {
			node.Metrics, err = metrics.NewCollector(node.Logger)
			return err
		}).
		Component("RPC engine", func(node *cmd.FlowNodeBuilder) (module.ReadyDoneAware, error) {
			rpcEng, err := engine.New(node.Network, node.Logger, node.Me, rpcConf)
			return rpcEng, err
		}).
		Run("ghost")
}
