package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/yamamushi/kmud-2020/types"
	"net/http"
)

func makeAuthEndpoint(svc AccountManagerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(authRequest)
		token, err := svc.Auth(req.Secret, req.Username, req.HashedPass)
		if err != nil {
			return authResponse{token, err.Error()}, nil
		}
		return authResponse{token, ""}, nil
	}
}

func makeAccountInfoEndpoint(svc AccountManagerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(accountInfoRequest)
		field, err := svc.AccountInfo(req.Secret, req.AuthToken, req.Account, req.Field)
		return accountInfoResponse{Account: field, Err: err.Error()}, nil
	}
}

func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request authRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAccountInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request accountInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type authRequest struct {
	Secret     string `json:"secret"`
	Username   string `json:"username"`
	HashedPass string `json:"hashedpass"`
}

type authResponse struct {
	AuthToken string `json:"authtoken"`
	Err       string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type accountInfoRequest struct {
	Secret    string `json:"secret"`
	AuthToken string `json:"authtoken"`
	Account   string `json:"account"`
	Field     string `json:"field"`
}

type accountInfoResponse struct {
	Account types.Account `json:"account"`
	Err     string        `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}
