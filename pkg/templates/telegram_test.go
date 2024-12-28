package templates

import (
	"html/template"
	"main/pkg/config/types"
	loggerPkg "main/pkg/logger"
	"main/pkg/messages"
	amountPkg "main/pkg/types/amount"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelegramTemplateManagerGetTemplateFailedToLoad(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)

	_, err = manager.Render("not-existing", nil)
	require.Error(t, err)
}

func TestTelegramTemplateManagerGetTemplateFailedToRender(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)

	_, err = manager.Render("Tx", nil)
	require.Error(t, err)
}

func TestTelegramTemplateManagerGetTemplateOk(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)

	_, err = manager.Render("Help", "1.2.3")
	require.NoError(t, err)

	_, err = manager.Render("Help", "1.2.4")
	require.NoError(t, err)
}

func TestTelegramTemplateManagerGetTemplateSerializeLink(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)
	require.Equal(t, template.HTML("<a href='https://example.com'>LinkTitle</a>"), manager.SerializeLink(&types.Link{
		Href:  "https://example.com",
		Value: "LinkValue",
		Title: "LinkTitle",
	}))
	require.Equal(t, template.HTML("<a href='https://example.com'>LinkValue</a>"), manager.SerializeLink(&types.Link{
		Href:  "https://example.com",
		Value: "LinkValue",
	}))
	require.Equal(t, template.HTML("LinkValue"), manager.SerializeLink(&types.Link{
		Value: "LinkValue",
	}))
	require.Equal(t, template.HTML("LinkTitle"), manager.SerializeLink(&types.Link{
		Value: "LinkValue",
		Title: "LinkTitle",
	}))
}

func TestTelegramTemplateManagerGetTemplateSerializeAmount(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)
	require.Equal(t, template.HTML("1.234567 DENOM"), manager.SerializeAmount(amountPkg.Amount{
		Value: big.NewFloat(1.23456789),
		Denom: "DENOM",
	}))

	require.Equal(t, template.HTML("1.234567 DENOM ($9.876)"), manager.SerializeAmount(amountPkg.Amount{
		Value:    big.NewFloat(1.23456789),
		Denom:    "DENOM",
		PriceUSD: big.NewFloat(9.876543),
	}))
}

func TestTelegramTemplateManagerGetTemplateSerializeDate(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	date, err := time.Parse(time.RFC3339, "2024-12-27T11:09:00Z")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)
	require.Equal(t, template.HTML("27 Dec 24 14:09 MSK"), manager.SerializeDate(date))
}

func TestTelegramTemplateManagerGetTemplateSerializeMessageFailedToLoad(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)
	require.Equal(
		t,
		template.HTML("Error loading template: <code>template: pattern matches no files: `telegram/MsgNotExistingMessage.html`</code>"),
		manager.SerializeMessage(&messages.MsgNotExistingMessage{}),
	)
}

func TestTelegramTemplateManagerGetTemplateSerializeMessageFailedToRender(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)
	require.Equal(
		t,
		template.HTML("Error rendering template: <code>template: cosmos.bank.v1beta1.MsgSend.html:2:9: executing &#34;cosmos.bank.v1beta1.MsgSend.html&#34; at &lt;SerializeLink .From&gt;: error calling SerializeLink: runtime error: invalid memory address or nil pointer dereference</code>"),
		manager.SerializeMessage(&messages.MsgSend{}),
	)
}

func TestTelegramTemplateManagerGetTemplateSerializeMessageOk(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)

	manager := NewTelegramTemplateManager(loggerPkg.GetNopLogger(), timezone)
	require.Equal(
		t,
		template.HTML("❌ This message type is not supported yet: <code>random</code>\n"),
		manager.SerializeMessage(&messages.MsgUnsupportedMessage{MsgType: "random"}),
	)
}
