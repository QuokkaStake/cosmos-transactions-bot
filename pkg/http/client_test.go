package http

import (
	"errors"
	"main/assets"
	loggerPkg "main/pkg/logger"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestHttpClientErrorCreating(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	client := NewClient(logger, "", "chain")
	err, queryInfo := client.Get("://test", nil)
	require.Error(t, err)
	require.False(t, queryInfo.Success)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientQueryFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)
	logger := loggerPkg.GetNopLogger()
	client := NewClient(logger, "https://example.com", "chain")

	var response interface{}
	err, queryInfo := client.Get("/", &response)
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
	require.False(t, queryInfo.Success)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientJsonParseFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("invalid-json.json")),
	)
	logger := loggerPkg.GetNopLogger()
	client := NewClient(logger, "https://example.com", "chain")
	var response interface{}

	err, queryInfo := client.Get("/", &response)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid character")
	require.False(t, queryInfo.Success)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientBadHttpCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/",
		httpmock.NewBytesResponder(500, assets.GetBytesOrPanic("error.json")),
	)
	logger := loggerPkg.GetNopLogger()
	client := NewClient(logger, "https://example.com", "chain")

	var response interface{}
	err, queryInfo := client.Get("/", &response)
	require.Error(t, err)
	require.ErrorContains(t, err, "bad HTTP code")
	require.False(t, queryInfo.Success)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("error.json")),
	)
	logger := loggerPkg.GetNopLogger()
	client := NewClient(logger, "https://example.com", "chain")

	var response interface{}
	err, queryInfo := client.GetWithHeaders("/", &response, map[string]string{
		"User-Agent": "custom",
	})
	require.NoError(t, err)
	require.True(t, queryInfo.Success)
}
