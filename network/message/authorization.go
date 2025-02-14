package message

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/libp2p/message"
	"github.com/koko1123/flow-go-1/model/messages"
	"github.com/koko1123/flow-go-1/network/channels"
)

type ChannelAuthConfig struct {
	// AuthorizedRoles list of roles authorized to send this message on the channel.
	AuthorizedRoles flow.RoleList

	// AllowedProtocols list of protocols the message is allowed to be sent on. Currently AllowedProtocols is expected to have
	// exactly one element in the list. This is due to the fact that currently there are no messages that are used with both protocols aside from TestMessage.
	AllowedProtocols Protocols
}

var authorizationConfigs map[string]MsgAuthConfig

// MsgAuthConfig contains authorization information for a specific flow message. The authorization
// is represented as a map from network channel -> list of all roles allowed to send the message on
// the channel.
type MsgAuthConfig struct {
	// Name is the string representation of the message type.
	Name string
	// Type is a func that returns a new instance of message type.
	Type func() interface{}
	// Config is the mapping of network channel to list of authorized flow roles.
	Config map[channels.Channel]ChannelAuthConfig
}

// EnsureAuthorized checks if the specified role is authorized to send the message on the provided channel and
// asserts that the message is authorized to be sent on the channel.
// Expected error returns during normal operations:
//   - ErrUnauthorizedMessageOnChannel: the channel is not included in the message's list of authorized channels
//   - ErrUnauthorizedRole: the role is not included in the message's list of authorized roles for the provided channel
//   - ErrUnauthorizedUnicastOnChannel: the message is not authorized to be sent via unicast protocol.
//   - ErrUnauthorizedPublishOnChannel: the message is not authorized to be sent via publish protocol.
func (m MsgAuthConfig) EnsureAuthorized(role flow.Role, channel channels.Channel, protocol Protocol) error {
	authConfig, ok := m.Config[channel]
	if !ok {
		return ErrUnauthorizedMessageOnChannel
	}

	if !authConfig.AuthorizedRoles.Contains(role) {
		return ErrUnauthorizedRole
	}

	if !authConfig.AllowedProtocols.Contains(protocol) {
		return NewUnauthorizedProtocolError(protocol)
	}

	return nil
}

func initializeMessageAuthConfigsMap() {
	authorizationConfigs = make(map[string]MsgAuthConfig)

	// consensus
	authorizationConfigs[BlockProposal] = MsgAuthConfig{
		Name: BlockProposal,
		Type: func() interface{} {
			return new(messages.BlockProposal)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ConsensusCommittee: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolPublish},
			},
			channels.PushBlocks: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias ReceiveBlocks = PushBlocks
		},
	}
	authorizationConfigs[BlockVote] = MsgAuthConfig{
		Name: BlockVote,
		Type: func() interface{} {
			return new(messages.BlockVote)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ConsensusCommittee: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}

	// protocol state sync
	authorizationConfigs[SyncRequest] = MsgAuthConfig{
		Name: SyncRequest,
		Type: func() interface{} {
			return new(messages.SyncRequest)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.SyncCommittee: {
				AuthorizedRoles:  flow.Roles(),
				AllowedProtocols: Protocols{ProtocolPublish},
			},
			channels.SyncClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolPublish},
			},
		},
	}
	authorizationConfigs[SyncResponse] = MsgAuthConfig{
		Name: SyncResponse,
		Type: func() interface{} {
			return new(messages.SyncResponse)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.SyncCommittee: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
			channels.SyncClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}
	authorizationConfigs[RangeRequest] = MsgAuthConfig{
		Name: RangeRequest,
		Type: func() interface{} {
			return new(messages.RangeRequest)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.SyncCommittee: {
				AuthorizedRoles:  flow.Roles(),
				AllowedProtocols: Protocols{ProtocolPublish},
			},
			channels.SyncClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolPublish},
			},
		},
	}
	authorizationConfigs[BatchRequest] = MsgAuthConfig{
		Name: BatchRequest,
		Type: func() interface{} {
			return new(messages.BatchRequest)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.SyncCommittee: {
				AuthorizedRoles:  flow.Roles(),
				AllowedProtocols: Protocols{ProtocolPublish},
			},
			channels.SyncClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolPublish},
			},
		},
	}
	authorizationConfigs[BlockResponse] = MsgAuthConfig{
		Name: BlockResponse,
		Type: func() interface{} {
			return new(messages.BlockResponse)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.SyncCommittee: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}

	// cluster consensus
	authorizationConfigs[ClusterBlockProposal] = MsgAuthConfig{
		Name: ClusterBlockProposal,
		Type: func() interface{} {
			return new(messages.ClusterBlockProposal)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ConsensusClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolPublish},
			},
		},
	}
	authorizationConfigs[ClusterBlockVote] = MsgAuthConfig{
		Name: ClusterBlockVote,
		Type: func() interface{} {
			return new(messages.ClusterBlockVote)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ConsensusClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}
	authorizationConfigs[ClusterBlockResponse] = MsgAuthConfig{
		Name: ClusterBlockResponse,
		Type: func() interface{} {
			return new(messages.ClusterBlockResponse)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.SyncClusterPrefix: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}

	// collections, guarantees & transactions
	authorizationConfigs[CollectionGuarantee] = MsgAuthConfig{
		Name: CollectionGuarantee,
		Type: func() interface{} {
			return new(flow.CollectionGuarantee)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.PushGuarantees: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias ReceiveGuarantees = PushGuarantees
		},
	}
	authorizationConfigs[TransactionBody] = MsgAuthConfig{
		Name: TransactionBody,
		Type: func() interface{} {
			return new(flow.TransactionBody)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.PushTransactions: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias ReceiveTransactions = PushTransactions
		},
	}

	// core messages for execution & verification
	authorizationConfigs[ExecutionReceipt] = MsgAuthConfig{
		Name: ExecutionReceipt,
		Type: func() interface{} {
			return new(flow.ExecutionReceipt)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.PushReceipts: {
				AuthorizedRoles:  flow.RoleList{flow.RoleExecution},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias ReceiveReceipts = PushReceipts
		},
	}
	authorizationConfigs[ResultApproval] = MsgAuthConfig{
		Name: ResultApproval,
		Type: func() interface{} {
			return new(flow.ResultApproval)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.PushApprovals: {
				AuthorizedRoles:  flow.RoleList{flow.RoleVerification},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias ReceiveApprovals = PushApprovals
		},
	}

	// data exchange for execution of blocks
	authorizationConfigs[ChunkDataRequest] = MsgAuthConfig{
		Name: ChunkDataRequest,
		Type: func() interface{} {
			return new(messages.ChunkDataRequest)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.RequestChunks: {
				AuthorizedRoles:  flow.RoleList{flow.RoleVerification},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias RequestChunks = ProvideChunks
		},
	}
	authorizationConfigs[ChunkDataResponse] = MsgAuthConfig{
		Name: ChunkDataResponse,
		Type: func() interface{} {
			return new(messages.ChunkDataResponse)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ProvideChunks: {
				AuthorizedRoles:  flow.RoleList{flow.RoleExecution},
				AllowedProtocols: Protocols{ProtocolUnicast},
			}, // channel alias RequestChunks = ProvideChunks
		},
	}

	// result approvals
	authorizationConfigs[ApprovalRequest] = MsgAuthConfig{
		Name: ApprovalRequest,
		Type: func() interface{} {
			return new(messages.ApprovalRequest)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.RequestApprovalsByChunk: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolPublish},
			}, // channel alias ProvideApprovalsByChunk  = RequestApprovalsByChunk
		},
	}
	authorizationConfigs[ApprovalResponse] = MsgAuthConfig{
		Name: ApprovalResponse,
		Type: func() interface{} {
			return new(messages.ApprovalResponse)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ProvideApprovalsByChunk: {
				AuthorizedRoles:  flow.RoleList{flow.RoleVerification},
				AllowedProtocols: Protocols{ProtocolUnicast},
			}, // channel alias ProvideApprovalsByChunk  = RequestApprovalsByChunk
		},
	}

	// generic entity exchange engines
	authorizationConfigs[EntityRequest] = MsgAuthConfig{
		Name: EntityRequest,
		Type: func() interface{} {
			return new(messages.EntityRequest)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.RequestReceiptsByBlockID: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
			channels.RequestCollections: {
				AuthorizedRoles:  flow.RoleList{flow.RoleAccess, flow.RoleExecution},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}
	authorizationConfigs[EntityResponse] = MsgAuthConfig{
		Name: EntityResponse,
		Type: func() interface{} {
			return new(messages.EntityResponse)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.ProvideReceiptsByBlockID: {
				AuthorizedRoles:  flow.RoleList{flow.RoleExecution},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
			channels.ProvideCollections: {
				AuthorizedRoles:  flow.RoleList{flow.RoleCollection},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}

	// testing
	authorizationConfigs[TestMessage] = MsgAuthConfig{
		Name: TestMessage,
		Type: func() interface{} {
			return new(message.TestMessage)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.TestNetworkChannel: {
				AuthorizedRoles:  flow.Roles(),
				AllowedProtocols: Protocols{ProtocolPublish, ProtocolUnicast},
			},
			channels.TestMetricsChannel: {
				AuthorizedRoles:  flow.Roles(),
				AllowedProtocols: Protocols{ProtocolPublish, ProtocolUnicast},
			},
		},
	}

	// DKG
	authorizationConfigs[DKGMessage] = MsgAuthConfig{
		Name: DKGMessage,
		Type: func() interface{} {
			return new(messages.DKGMessage)
		},
		Config: map[channels.Channel]ChannelAuthConfig{
			channels.DKGCommittee: {
				AuthorizedRoles:  flow.RoleList{flow.RoleConsensus},
				AllowedProtocols: Protocols{ProtocolUnicast},
			},
		},
	}
}

// GetMessageAuthConfig checks the underlying type and returns the correct
// message auth Config.
// Expected error returns during normal operations:
//   - ErrUnknownMsgType : if underlying type of v does  not match any of the known message types
func GetMessageAuthConfig(v interface{}) (MsgAuthConfig, error) {
	switch v.(type) {
	// consensus
	case *messages.BlockProposal:
		return authorizationConfigs[BlockProposal], nil
	case *messages.BlockVote:
		return authorizationConfigs[BlockVote], nil

	// protocol state sync
	case *messages.SyncRequest:
		return authorizationConfigs[SyncRequest], nil
	case *messages.SyncResponse:
		return authorizationConfigs[SyncResponse], nil
	case *messages.RangeRequest:
		return authorizationConfigs[RangeRequest], nil
	case *messages.BatchRequest:
		return authorizationConfigs[BatchRequest], nil
	case *messages.BlockResponse:
		return authorizationConfigs[BlockResponse], nil

	// cluster consensus
	case *messages.ClusterBlockProposal:
		return authorizationConfigs[ClusterBlockProposal], nil
	case *messages.ClusterBlockVote:
		return authorizationConfigs[ClusterBlockVote], nil
	case *messages.ClusterBlockResponse:
		return authorizationConfigs[ClusterBlockResponse], nil

	// collections, guarantees & transactions
	case *flow.CollectionGuarantee:
		return authorizationConfigs[CollectionGuarantee], nil
	case *flow.TransactionBody:
		return authorizationConfigs[TransactionBody], nil

	// core messages for execution & verification
	case *flow.ExecutionReceipt:
		return authorizationConfigs[ExecutionReceipt], nil
	case *flow.ResultApproval:
		return authorizationConfigs[ResultApproval], nil

	// data exchange for execution of blocks
	case *messages.ChunkDataRequest:
		return authorizationConfigs[ChunkDataRequest], nil
	case *messages.ChunkDataResponse:
		return authorizationConfigs[ChunkDataResponse], nil

	// result approvals
	case *messages.ApprovalRequest:
		return authorizationConfigs[ApprovalRequest], nil
	case *messages.ApprovalResponse:
		return authorizationConfigs[ApprovalResponse], nil

	// generic entity exchange engines
	case *messages.EntityRequest:
		return authorizationConfigs[EntityRequest], nil
	case *messages.EntityResponse:
		return authorizationConfigs[EntityResponse], nil

	// testing
	case *message.TestMessage:
		return authorizationConfigs[TestMessage], nil

	// dkg
	case *messages.DKGMessage:
		return authorizationConfigs[DKGMessage], nil

	default:
		return MsgAuthConfig{}, NewUnknownMsgTypeErr(v)
	}
}

// GetAllMessageAuthConfigs returns all the configured message auth configurations.
func GetAllMessageAuthConfigs() []MsgAuthConfig {
	configs := make([]MsgAuthConfig, len(authorizationConfigs))

	i := 0
	for _, config := range authorizationConfigs {
		configs[i] = config
		i++
	}

	return configs
}
