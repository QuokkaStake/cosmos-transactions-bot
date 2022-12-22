package main

type TendermintRpcStatus struct {
	Success bool
	Error   error
}

type Message interface {
	Type() string
}

type Report struct {
	Chain      Chain
	Node       string
	Reportable Reportable
}

type Reportable interface {
	Type() string
	GetHash() string
	GetMessages() []Message
}

type Tx struct {
	Hash     string
	Memo     string
	Height   int64
	Messages []Message
}

func (tx Tx) GetMessages() []Message {
	return tx.Messages
}

func (tx Tx) Type() string {
	return "Tx"
}

func (tx Tx) GetHash() string {
	return tx.Hash
}

type TxError struct {
	Error error
}

func (txError TxError) GetMessages() []Message {
	return []Message{}
}

func (txError TxError) Type() string {
	return "TxError"
}

func (txError TxError) GetHash() string {
	return "TxError"
}

type MsgError struct {
	Error error
}

func (m MsgError) Type() string {
	return "MsgError"
}

type MessageParser func([]byte) (Message, error)

type Reporter interface {
	Init()
	Name() string
	Enabled() bool
	Send(Report) error
}
