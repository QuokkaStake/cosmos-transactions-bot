package templates

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/utils"
	"main/templates"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog"
)

type TelegramTemplateManager struct {
	Logger    zerolog.Logger
	Templates map[string]*template.Template
	Timezone  *time.Location
}

func NewTelegramTemplateManager(
	logger *zerolog.Logger,
	timezone *time.Location,
) *TelegramTemplateManager {
	return &TelegramTemplateManager{
		Logger:    logger.With().Str("component", "telegram_template_manager").Logger(),
		Timezone:  timezone,
		Templates: map[string]*template.Template{},
	}
}

func (m *TelegramTemplateManager) GetTemplate(name string) (*template.Template, error) {
	if cachedTemplate, ok := m.Templates[name]; ok {
		m.Logger.Trace().Str("type", name).Msg("Using cached template")
		return cachedTemplate, nil
	}

	m.Logger.Trace().Str("type", name).Msg("Loading template")

	filename := fmt.Sprintf("%s.html", utils.RemoveFirstSlash(name))

	t, err := template.New(filename).Funcs(template.FuncMap{
		"SerializeLink":    m.SerializeLink,
		"SerializeAmount":  m.SerializeAmount,
		"SerializeDate":    m.SerializeDate,
		"SerializeMessage": m.SerializeMessage,
	}).ParseFS(templates.TemplatesFs, "telegram/"+filename)
	if err != nil {
		return nil, err
	}

	m.Templates[name] = t

	return t, nil
}

func (m *TelegramTemplateManager) Render(templateName string, data interface{}) (string, error) {
	reportTemplate, err := m.GetTemplate(templateName)
	if err != nil {
		m.Logger.Error().Err(err).Str("type", templateName).Msg("Error loading template")
		return "", err
	}

	var buffer bytes.Buffer
	err = reportTemplate.Execute(&buffer, data)
	if err != nil {
		m.Logger.Error().Err(err).Str("type", templateName).Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), err
}

func (m *TelegramTemplateManager) SerializeLink(link *configTypes.Link) template.HTML {
	value := link.Title
	if value == "" {
		value = link.Value
	}

	if link.Href != "" {
		return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link.Href, value))
	}

	return template.HTML(value)
}

func (m *TelegramTemplateManager) SerializeAmount(amount amount.Amount) template.HTML {
	if amount.PriceUSD == nil {
		return template.HTML(fmt.Sprintf(
			"%s %s",
			utils.StripTrailingDigits(humanize.BigCommaf(amount.Value), 6),
			amount.Denom,
		))
	}

	return template.HTML(fmt.Sprintf(
		"%s %s ($%s)",
		utils.StripTrailingDigits(humanize.BigCommaf(amount.Value), 6),
		amount.Denom,
		utils.StripTrailingDigits(humanize.BigCommaf(amount.PriceUSD), 3),
	))
}

func (m *TelegramTemplateManager) SerializeDate(date time.Time) template.HTML {
	return template.HTML(date.In(m.Timezone).Format(time.RFC822))
}

func (m *TelegramTemplateManager) SerializeMessage(msg types.Message) template.HTML {
	msgType := msg.Type()

	reporterTemplate, err := m.GetTemplate(msgType)
	if err != nil {
		m.Logger.Error().Err(err).Str("type", msgType).Msg("Error loading template")
		return template.HTML(fmt.Sprintf("Error loading template: <code>%s</code>", html.EscapeString(err.Error())))
	}

	var buffer bytes.Buffer
	err = reporterTemplate.Execute(&buffer, msg)
	if err != nil {
		m.Logger.Error().Err(err).Str("type", msgType).Msg("Error rendering template")
		return template.HTML(fmt.Sprintf("Error rendering template: <code>%s</code>", html.EscapeString(err.Error())))
	}

	return template.HTML(buffer.String())
}
