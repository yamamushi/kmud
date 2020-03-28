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
		req := request.(types.AuthRequest)
		token, err := svc.Auth(req.Secret, req.Username, req.HashedPass, conf, db)
		if err != nil {
			return types.AuthResponse{token, err.Error()}, nil
		}
		return types.AuthResponse{token, ""}, nil
	}
}

func makeAccountInfoEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(types.AccountInfoRequest)
		field, err := svc.AccountInfo(req.Secret, req.Token, req.Field, conf, db)
		return types.AccountInfoResponse{Account: field, Err: err.Error()}, nil
	}
}

func makeAccountRegistrationEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(types.AccountRegistrationRequest)
		err := svc.AccountRegistration(req.Secret, req.Username, req.Email, req.HashedPass, conf, db)
		return types.AccountRegistrationResponse{Err: err.Error()}, nil
	}
}

func makeSearchEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(types.SearchRequest)
		accounts, err := svc.Search(req.Secret, req.Token, req.Account, conf, db)
		return types.SearchResponse{Accounts: accounts, Err: err.Error()}, nil
	}
}

func makeModifyEndpoint(svc AccountManagerService, conf *config.Config, db *database.DatabaseHandler) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(types.ModifyRequest)
		account, err := svc.Modify(req.Secret, req.Token, req.Account, conf, db)
		return types.ModifyResponse{Account: account, Err: err.Error()}, nil
	}
}

func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request types.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAccountInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request types.AccountInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAccountRegistrationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request types.AccountRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeModifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request types.ModifyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request types.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
