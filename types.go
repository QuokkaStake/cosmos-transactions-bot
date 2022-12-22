package main

type TendermintRPCStatus struct {
	Success bool
	Error   error
}

type MessageParser func([]byte, Chain) (Message, error)

type Reporter interface {
	Init()
	Name() string
	Enabled() bool
	Send(Report) error
}

type Link struct {
	Href  string
	Title string
}

type Amount struct {
	Value    float64
	Denom    string
	PriceUSD float64
}
