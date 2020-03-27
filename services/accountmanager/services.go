package main

import (
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type AccountManagerService interface {
	Auth(string, string, string) (string, error)
	AccountInfo(string, string, string, string) (types.Account, error)
}

type accountManagerService struct{}

func (accountManagerService) Auth(secret string, username string, hashedpass string) (string, error) {

	// Check Hashed Pass
	// Write token to table entry
	// Notify User Manager that user is online

	return "username:sha256token", utils.EmptyError()
}

func (accountManagerService) AccountInfo(secret string, token string, account string, field string) (types.Account, error) {

	return types.Account{Username: "Test", Email: "test@test.com"}, utils.EmptyError()
}
