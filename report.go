package main

type Message interface {
	Type() string
	GetAdditionalData(DataFetcher)
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
	GetAdditionalData(DataFetcher)
}

type Tx struct {
	Hash     Link
	Memo     string
	Height   Link
	Messages []Message
}

func (tx Tx) GetMessages() []Message {
	return tx.Messages
}

func (tx Tx) Type() string {
	return "Tx"
}

func (tx Tx) GetHash() string {
	return tx.Hash.Title
}

func (tx *Tx) GetAdditionalData(fetcher DataFetcher) {
	for _, msg := range tx.Messages {
		msg.GetAdditionalData(fetcher)
	}
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

func (txError *TxError) GetAdditionalData(fetcher DataFetcher) {

}

type MsgError struct {
	Error error
}

func (m MsgError) Type() string {
	return "MsgError"
}

func (m *MsgError) GetAdditionalData(fetcher DataFetcher) {

}
