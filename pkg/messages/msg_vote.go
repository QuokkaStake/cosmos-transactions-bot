package messages

import (
	"fmt"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosGovEvents "github.com/cosmos/cosmos-sdk/x/gov/types"
	cosmosGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"main/pkg/types/responses"
	"strconv"
)

type MsgVote struct {
	Voter      configTypes.Link
	ProposalID configTypes.Link
	Proposal   *responses.Proposal
	Option     cosmosGovTypes.VoteOption
}

func ParseMsgVote(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosGovTypes.MsgVote
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgVote{
		Voter:      chain.GetWalletLink(parsedMessage.Voter),
		ProposalID: chain.GetProposalLink(strconv.FormatUint(parsedMessage.ProposalId, 10)),
		Option:     parsedMessage.Option,
	}, nil
}

func (m MsgVote) Type() string {
	return "/cosmos.gov.v1beta1.MsgVote"
}

func (m *MsgVote) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	proposal, found := fetcher.GetProposal(m.ProposalID.Value)
	if found {
		m.Proposal = proposal
		m.ProposalID.Title = fmt.Sprintf("#%s: %s", m.ProposalID.Value, proposal.Content.Title)
	} else {
		m.ProposalID.Title = fmt.Sprintf("#%s", m.ProposalID.Value)
	}
}

func (m *MsgVote) GetVote() string {
	switch m.Option {
	case cosmosGovTypes.OptionYes:
		return "Yes"
	case cosmosGovTypes.OptionAbstain:
		return "Abstain"
	case cosmosGovTypes.OptionNo:
		return "No"
	case cosmosGovTypes.OptionNoWithVeto:
		return "No with veto"
	default:
		return m.Option.String()
	}
}

func (m *MsgVote) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosGovEvents.EventTypeProposalVote, cosmosGovEvents.AttributeKeyProposalID, m.ProposalID.Value),
		event.From(cosmosGovEvents.EventTypeProposalVote, cosmosGovEvents.AttributeKeyOption, m.Option.String()),
	}
}
