package responses

type IbcChannelResponse struct {
	Channel IbcChannel `json:"channel"`
}

type IbcChannel struct {
	ConnectionHops []string `json:"connection_hops"`
}
