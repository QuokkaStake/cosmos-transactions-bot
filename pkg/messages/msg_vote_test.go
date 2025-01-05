package messages

import (
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/types/event"
	"main/pkg/types/responses"
	"testing"

	cosmosGovEvents "github.com/cosmos/cosmos-sdk/x/gov/types"
	cosmosGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgVoteParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosGovTypes.MsgVote{
		Voter:      "voter",
		ProposalId: 1,
		Option:     cosmosGovTypes.OptionYes,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgVote(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgVote([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgVoteBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosGovTypes.MsgVote{
		Voter:      "voter",
		ProposalId: 1,
		Option:     cosmosGovTypes.OptionYes,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgVote(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.gov.v1beta1.MsgVote", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.gov.v1beta1.MsgVote"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "voter"),
		event.From(cosmosGovEvents.EventTypeProposalVote, cosmosGovEvents.AttributeKeyProposalID, "1"),
		event.From(cosmosGovEvents.EventTypeProposalVote, cosmosGovEvents.AttributeKeyOption, "VOTE_OPTION_YES"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgVoteGetVote(t *testing.T) {
	t.Parallel()

	require.Equal(t, "Yes", (&MsgVote{Option: cosmosGovTypes.OptionYes}).GetVote())
	require.Equal(t, "No", (&MsgVote{Option: cosmosGovTypes.OptionNo}).GetVote())
	require.Equal(t, "No with veto", (&MsgVote{Option: cosmosGovTypes.OptionNoWithVeto}).GetVote())
	require.Equal(t, "Abstain", (&MsgVote{Option: cosmosGovTypes.OptionAbstain}).GetVote())
	require.Equal(t, "Empty", (&MsgVote{Option: cosmosGovTypes.OptionEmpty}).GetVote())
	require.Equal(t, "5", (&MsgVote{Option: 5}).GetVote())
}

func TestMsgVotePopulatePresent(t *testing.T) {
	t.Parallel()

	msg := &cosmosGovTypes.MsgVote{
		Voter:      "voter",
		ProposalId: 1,
		Option:     cosmosGovTypes.OptionYes,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains:      configTypes.Chains{{Name: "chain"}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgVote(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "voter", "voter_alias")
	require.NoError(t, err)

	dataFetcher.Cache.Set("chain_proposal_1", &responses.Proposal{
		ProposalID: "1",
		Content:    responses.ProposalContent{Title: "Title"},
	})

	parsed.GetAdditionalData(dataFetcher, "subscription")

	msgSend, _ := parsed.(*MsgVote)

	require.Equal(t, "voter_alias", msgSend.Voter.Title)
	require.Equal(t, "#1: Title", msgSend.ProposalID.Title)
}

func TestMsgVotePopulateAbsent(t *testing.T) {
	t.Parallel()

	msg := &cosmosGovTypes.MsgVote{
		Voter:      "voter",
		ProposalId: 1,
		Option:     cosmosGovTypes.OptionYes,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains:      configTypes.Chains{{Name: "chain"}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgVote(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "voter", "voter_alias")
	require.NoError(t, err)

	dataFetcher.Cache.Set("chain_proposal_1", nil)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	msgSend, _ := parsed.(*MsgVote)

	require.Equal(t, "voter_alias", msgSend.Voter.Title)
	require.Equal(t, "#1", msgSend.ProposalID.Title)
}
