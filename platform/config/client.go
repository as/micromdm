package config

import (
	"net/url"

	"github.com/as/micromdm/pkg/httputil"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPClient(instance, token string, logger log.Logger, opts ...httptransport.ClientOption) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var saveEndpoint endpoint.Endpoint
	{
		saveEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/config/certificate"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeSavePushCertificateResponse,
			opts...,
		).Endpoint()
	}

	var applyDEPTokensEndpoint endpoint.Endpoint
	{
		applyDEPTokensEndpoint = httptransport.NewClient(
			"PUT",
			httputil.CopyURL(u, "/v1/dep-tokens"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeApplyDEPTokensResponse,
			opts...,
		).Endpoint()
	}

	var getDEPTokensEndpoint endpoint.Endpoint
	{
		getDEPTokensEndpoint = httptransport.NewClient(
			"GET",
			httputil.CopyURL(u, "/v1/dep-tokens"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeGetDEPTokensResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		SavePushCertificateEndpoint: saveEndpoint,
		ApplyDEPTokensEndpoint:      applyDEPTokensEndpoint,
		GetDEPTokensEndpoint:        getDEPTokensEndpoint,
	}, nil
}
