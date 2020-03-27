package main

import (
	"encoding/hex"
	"errors"
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/database"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type AccountManagerService interface {
	Auth(string, string, string, *config.Config, *database.DatabaseHandler) (string, error)
	AccountInfo(string, string, string, *config.Config, *database.DatabaseHandler) (types.Account, error)
	AccountRegistration(string, string, string, string, *config.Config, *database.DatabaseHandler) error
}

type accountManagerService struct {
}

func (accountManagerService) Auth(secret string, username string, hashedpass string, conf *config.Config, DB *database.DatabaseHandler) (string, error) {

	if secret != conf.Crypt.AccountManagerSecret {
		return "", errors.New("unauthorized request")
	}

	account := types.Account{}
	result, err := DB.FindOne(bson.M{"username": username}, conf.DB.MongoDB, "accounts")
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return "", errors.New("account not found")
		} else {
			return "", err
		}
	}
	account = BsonMapToAccount(result)
	inputPass := hex.EncodeToString([]byte(hashedpass))

	if inputPass != account.HashedPass {
		return "", errors.New("invalid password")
	}

	if account.Token == "" {
		token, err := utils.GetRandomToken()
		if err != nil {
			return "", errors.New("error creating user token: " + err.Error())
		}
		account.Token = token
	}

	err = DB.UpdateOne(bson.M{"username": username}, account, conf.DB.MongoDB, "accounts")
	if err != nil {
		return "", err
	}

	auth := account.Username + ":" + account.Token
	return auth, utils.EmptyError()
}

func (accountManagerService) AccountInfo(secret string, token string, field string, conf *config.Config, DB *database.DatabaseHandler) (types.Account, error) {

	if secret != conf.Crypt.AccountManagerSecret {
		return types.Account{}, errors.New("unauthorized request")
	}

	tokenfields := strings.Split(token, ":")
	if len(tokenfields) != 2 {
		return types.Account{}, errors.New("invalid token format")
	}

	if field == "" {
		field = "all"
	}

	result, err := DB.FindOne(bson.M{"username": tokenfields[0]}, conf.DB.MongoDB, "accounts")
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return types.Account{}, errors.New("unauthorized request")
		} else {
			return types.Account{}, err
		}
	}

	accountStruct := BsonMapToAccount(result)
	if accountStruct.Token != tokenfields[1] {
		return types.Account{}, errors.New("unauthorized request")
	}

	output := types.Account{}
	fields := strings.Split(field, "|")
	var unrecognized string
	output.Username = accountStruct.Username
	for _, filter := range fields {
		filter = strings.ToLower(filter)
		var found bool
		if filter == "email" || filter == "all" {
			output.Email = accountStruct.Email
			found = true
		}
		if filter == "permissions" || filter == "all" {
			output.Permissions = accountStruct.Permissions
			found = true
		}
		if filter == "characters" || filter == "all" {
			output.Characters = accountStruct.Characters
			found = true
		}
		if filter == "locked" || filter == "all" {
			output.Locked = accountStruct.Locked
			found = true
		}
		if !found {
			unrecognized = unrecognized + filter + ","
		}
	}
	if unrecognized != "" {
		unrecognized = utils.RemoveLastChar(unrecognized)
		return types.Account{}, errors.New("unrecognized fields: " + unrecognized)
	}

	return output, utils.EmptyError()
}

func (accountManagerService) AccountRegistration(secret string, username string, email string, hashedpass string, conf *config.Config, DB *database.DatabaseHandler) (err error) {

	if secret != conf.Crypt.AccountManagerSecret {
		return errors.New("unauthorized request")
	}

	if secret == "" || username == "" || email == "" || hashedpass == "" {
		return errors.New("invalid request")
	}

	_, err = DB.FindOne(bson.M{"username": username}, conf.DB.MongoDB, "accounts")
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	} else {
		return errors.New("account with username " + username + " already exists")
	}

	_, err = DB.FindOne(bson.M{"email": email}, conf.DB.MongoDB, "accounts")
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	} else {
		return errors.New("account with email " + email + " already exists")
	}

	hexPass := hex.EncodeToString([]byte(hashedpass))
	err = DB.Insert(types.Account{Username: username, Email: email, HashedPass: hexPass, Locked: "false", Permissions: []string{"user"}}, conf.DB.MongoDB, "accounts")
	if err != nil {
		return err
	}

	return utils.EmptyError()
}
