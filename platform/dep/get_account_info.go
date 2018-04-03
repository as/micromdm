package dep

import (
	"context"
	"errors"
	"net/http"

	"github.com/as/micromdm/dep"
	"github.com/as/micromdm/pkg/httputil"
	"github.com/go-kit/kit/endpoint"
)

func (svc *DEPService) GetAccountInfo(ctx context.Context) (*dep.Account, error) {
	if svc.client == nil {
		return nil, errors.New("DEP not configured yet. add a DEP token to enable DEP")
	}
	return svc.client.Account()
}

type getAccountInfoRequest struct{}
type getAccountInfoResponse struct {
	*dep.Account
	Err error `json:"err,omitempty"`
}

func (r getAccountInfoResponse) Failed() error { return r.Err }

func decodeGetAccountInfoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetAccountInfoResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getAccountInfoResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeGetAccountInfoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		account, err := svc.GetAccountInfo(ctx)
		return getAccountInfoResponse{Account: account, Err: err}, nil
	}
}

func (e Endpoints) GetAccountInfo(ctx context.Context) (*dep.Account, error) {
	request := getAccountInfoRequest{}
	response, err := e.GetAccountInfoEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(getAccountInfoResponse).Account, response.(getAccountInfoResponse).Err
}
