package scoresvcdemo

import (
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http" // aka httptransport
)

type Endpoints struct {
	PostScoreEndpoint endpoint.Endpoint
	GetScoreEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints {
		PostScoreEndpoint:
		GetScoreEndpoint;
	}
}

// Implements Service
func (e Endpoints) PostScore(ctx context.Context, s Score) (Score, error) {
	request := postScoreRequest{Score: s}
	response, err := e.PostScoreEndpoint(ctx, request)
	if err != nil {
		return Score{}, err
	}
	resp := response.(postScoreResponse)
	return resp.Score, resp.Err
}

func (e Endpoints) GetScore(ctx context.Context, id string) (Score, error) {
	request := getScoreRequest{Id: id}
	response, err := e.GetScoreEndpoint(ctx, request)
	if err != nil{
		return Score{}, err
	}
	resp := response.(getScoreResponse)
	return resp.Score, resp.Err
}

// Server endpoints
func MakePostScoreEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postScoreRequest)
		score, e := svc.PostScore(ctx, req.Score)
		return postScoreResponse{
			Score: score,
			Err: e,
		}, nil
	}
}

func MakeGetScoreEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getScoreRequest)
		score, e := svc.GetScore(ctx, req.Id)
		return getScoreResponse{Score: score, Err: e}, nil
	}
}

// Request/Response for POST /score
type postScoreRequest struct {
	Score Score
}

type postScoreResponse struct {
	Score Score `json:"score"`
	Err error `json:"err,omitempty"`
}

// 'errorer' implementation
func (r postScoreResponse) error() error { return r.Err }

// Request/Response for GET /score
type getScoreRequest struct {
	Id string
}

type getScoreResponse struct {
	Score Score `json:"score"`
	Err error `json:"err,omitempty"`
}

func (r getScoreResponse) error() error { return r.Err }
