package types

type TelegramConfig struct {
	Chat   int64
	Token  string
	Admins []int64
}

type Reporter struct {
	Name string
	Type string

	TelegramConfig *TelegramConfig
}
