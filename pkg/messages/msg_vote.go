package messages

import (
	"fmt"
	cosmosGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/responses"
	"strconv"
)

type MsgVote struct {
	Voter      configTypes.Link
	ProposalID configTypes.Link
	Proposal   *responses.Proposal
	Option     string
}

func ParseMsgVote(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosGovTypes.MsgVote
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgVote{
		Voter:      chain.GetWalletLink(parsedMessage.Voter),
		ProposalID: chain.GetProposalLink(strconv.FormatUint(parsedMessage.ProposalId, 10)),
		Option:     parsedMessage.Option.String(),
	}, nil
}

func (m MsgVote) Type() string {
	return "MsgVote"
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

func (m *MsgVote) GetValues() map[string]string {
	return map[string]string{
		"type":        "MsgVote",
		"voter":       m.Voter.Value,
		"proposal_id": m.ProposalID.Value,
	}
}
