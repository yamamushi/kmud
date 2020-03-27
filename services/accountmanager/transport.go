package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/database"
	"github.com/yamamushi/kmud-2020/types"
	"net/http"
)

func makeAuthEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(authRequest)
		token, err := svc.Auth(req.Secret, req.Username, req.HashedPass, conf, db)
		if err != nil {
			return authResponse{token, err.Error()}, nil
		}
		return authResponse{token, ""}, nil
	}
}

func makeAccountInfoEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(accountInfoRequest)
		field, err := svc.AccountInfo(req.Secret, req.Token, req.Field, conf, db)
		return accountInfoResponse{Account: field, Err: err.Error()}, nil
	}
}

func makeAccountRegistrationEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(accountRegistrationRequest)
		err := svc.AccountRegistration(req.Secret, req.Username, req.Email, req.HashedPass, conf, db)
		return accountRegistrationResponse{Err: err.Error()}, nil
	}
}

func makeSearchEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchRequest)
		account, err := svc.Search(req.Secret, req.Token, req.Account, conf, db)
		return searchResponse{Accounts: account, Err: err.Error()}, nil
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

func decodeAccountRegistrationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request accountRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request searchRequest
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
	Err       string `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}

type accountInfoRequest struct {
	Secret string `json:"secret"`
	Token  string `json:"token"`
	Field  string `json:"field"`
}

type accountInfoResponse struct {
	Account types.Account `json:"account"`
	Err     string        `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}

type accountRegistrationRequest struct {
	Secret     string `json:"secret"`
	Username   string `json:"username"`
	HashedPass string `json:"hashedpass"`
	Email      string `json:"email"`
}

type accountRegistrationResponse struct {
	Err string `json:"error"` // errors don't JSON-marshal, so we use a string
}

type searchRequest struct {
	Secret  string        `json:"secret"`
	Token   string        `json:"token"`
	Account types.Account `json:"account"`
}

type searchResponse struct {
	Accounts []types.Account `json:"accounts"`
	Err      string          `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}
