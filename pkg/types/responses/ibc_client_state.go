package responses

type IbcClientStateResponse struct {
	IdentifiedClientState IbcIdentifiedClientState `json:"identified_client_state"`
}

type IbcIdentifiedClientState struct {
	ClientState IbcClientState `json:"client_state"`
}

type IbcClientState struct {
	ChainId string `json:"chain_id"`
}
