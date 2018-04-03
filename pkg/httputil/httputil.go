package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewRouter(logger log.Logger) (*mux.Router, []httptransport.ServerOption) {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(ErrorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	return r, options
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(failer); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	errMap := map[string]interface{}{"error": err.Error()}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}

	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)

	enc.Encode(errMap)
}

// failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type failer interface {
	Failed() error
}

type errorWrapper struct {
	Error string `json:"error"`
}

func JSONErrorDecoder(r *http.Response) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("expected JSON formatted error, got Content-Type %s", contentType)
	}
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

func EncodeRequestWithToken(token string, next httptransport.EncodeRequestFunc) httptransport.EncodeRequestFunc {
	return func(ctx context.Context, r *http.Request, request interface{}) error {
		r.SetBasicAuth("micromdm", token)
		return next(ctx, r, request)
	}
}

func CopyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func DecodeJSONRequest(r *http.Request, into interface{}) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(into)
	return err
}

func DecodeJSONResponse(r *http.Response, into interface{}) error {
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return JSONErrorDecoder(r)
	}

	err := json.NewDecoder(r.Body).Decode(into)
	return err
}
