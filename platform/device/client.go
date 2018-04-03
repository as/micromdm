package device

import (
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/as/micromdm/pkg/httputil"
)

func NewHTTPClient(instance, token string, logger log.Logger, opts ...httptransport.ClientOption) (Service, error) {
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var listDevicesEndpoint endpoint.Endpoint
	{
		listDevicesEndpoint = httptransport.NewClient(
			"GET",
			httputil.CopyURL(u, "/v1/devices"),
			httputil.EncodeRequestWithToken(token, httptransport.EncodeJSONRequest),
			decodeListDevicesResponse,
			opts...,
		).Endpoint()
	}

	return Endpoints{
		ListDevicesEndpoint: listDevicesEndpoint,
	}, nil

}
