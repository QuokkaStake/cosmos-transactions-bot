package http

import (
	"encoding/json"
	"fmt"
	"main/pkg/types/query_info"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Client struct {
	logger  zerolog.Logger
	host    string
	timeout time.Duration
}

func NewClient(
	logger *zerolog.Logger,
	host string,
	chainName string,
) *Client {
	return &Client{
		logger: logger.With().
			Str("component", "tendermint_api_client").
			Str("url", host).
			Str("chain", chainName).
			Logger(),
		host:    host,
		timeout: 60 * time.Second,
	}
}

func (c *Client) Get(url string, target interface{}) (error, query_info.QueryInfo) {
	return c.GetWithHeaders(url, target, map[string]string{})
}

func (c *Client) GetWithHeaders(
	relativeURL string,
	target interface{},
	headers map[string]string,
) (error, query_info.QueryInfo) {
	url := fmt.Sprintf("%s%s", c.host, relativeURL)

	client := &http.Client{Timeout: c.timeout}
	start := time.Now()
	queryInfo := query_info.QueryInfo{
		Success: false,
		Node:    c.host,
		Time:    0,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err, queryInfo
	}

	req.Header.Set("User-Agent", "cosmos-transactions-bot")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.logger.Trace().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	queryInfo.Time = time.Since(start)
	if err != nil {
		c.logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err, queryInfo
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		c.logger.Warn().
			Str("url", url).
			Err(err).
			Int("status", res.StatusCode).
			Msg("Query returned bad HTTP code")
		return fmt.Errorf("bad HTTP code: %d", res.StatusCode), queryInfo
	}

	c.logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return err, queryInfo
	}

	queryInfo.Success = true
	return nil, queryInfo
}
