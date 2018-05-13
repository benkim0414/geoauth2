package endpoint

import (
	"context"

	"github.com/benkim0414/geoauth2/client"
	"github.com/benkim0414/geoauth2/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Endpoints collects all of the endpoints that compose a service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Endpoints struct {
	PostClientEndpoint   endpoint.Endpoint
	GetClientEndpoint    endpoint.Endpoint
	DeleteClientEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of
// the expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger) Endpoints {
	postClientEndpoint := MakePostClientEndpoint(svc)
	postClientEndpoint = LoggingMiddleware(log.With(logger, "method", "post_client"))(postClientEndpoint)

	getClientEndpoint := MakeGetClientEndpoint(svc)
	getClientEndpoint = LoggingMiddleware(log.With(logger, "method", "get_client"))(getClientEndpoint)

	deleteClientEndpoint := MakeDeleteClientEndpoint(svc)
	deleteClientEndpoint = LoggingMiddleware(log.With(logger, "method", "delete_client"))(deleteClientEndpoint)

	return Endpoints{
		PostClientEndpoint:   postClientEndpoint,
		GetClientEndpoint:    getClientEndpoint,
		DeleteClientEndpoint: deleteClientEndpoint,
	}
}

// MakePostClientEndpoint constructs a PostClient endpoint wrapping the service.
func MakePostClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostClientRequest)
		client, err := s.PostClient(ctx, req.Client)
		return PostClientResponse{Client: client, Err: err}, nil
	}
}

// MakeGetClientEndpoint constructs a GetClient endpoint wrapping the service.
func MakeGetClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetClientRequest)
		client, err := s.GetClient(ctx, req.ID)
		return GetClientResponse{Client: client, Err: err}, nil
	}
}

// MakeDeleteClientEndpoint constructs a DeleteClient endpoint wrapping the service.
func MakeDeleteClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteClientRequest)
		err = s.DeleteClient(ctx, req.ID)
		return DeleteClientResponse{Err: err}, nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// PostClientRequest collects the request parameters for the PostClient method.
type PostClientRequest struct {
	Client *client.Client `json:"client"`
}

// PostClientResponse collects the response values for the PostClient method.
type PostClientResponse struct {
	Client *client.Client `json:"client"`
	Err    error          `json:"error"`
}

// Failed implements Failer.
func (r PostClientResponse) Failed() error { return r.Err }

// GetClientRequest collects the request parameters for the GetClient method.
type GetClientRequest struct {
	ID string `json:"id"`
}

// GetClientResponse collects the response values for the GetClient method.
type GetClientResponse struct {
	Client *client.Client `json:"client"`
	Err    error          `json:"error"`
}

// Failed implements Failer.
func (r GetClientResponse) Failed() error { return r.Err }

// DeleteClientReqeust collects the request parameters for the DeleteClient method.
type DeleteClientRequest struct {
	ID string `json:"id"`
}

// DeleteClientResponse collects the response values for the DeleteClient method.
type DeleteClientResponse struct {
	Err error `json:"error"`
}

// Failed implements Failer.
func (r DeleteClientResponse) Failed() error { return r.Err }
