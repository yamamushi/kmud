package main

import (
	"encoding/hex"
	"errors"
	"github.com/badoux/checkmail"
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
	Search(string, string, types.Account, *config.Config, *database.DatabaseHandler) ([]types.Account, error)
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
	account = utils.BsonMapToAccount(result)
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

	accountStruct, err := utils.ValidateRequest(secret, token, "", "moderators", conf, DB)
	if err != nil {
		return types.Account{}, err
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

	err = checkmail.ValidateFormat(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	err = checkmail.ValidateHost(email)
	if err != nil {
		return errors.New("invalid email domain")
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
	err = DB.Insert(types.Account{Username: username, Email: email, HashedPass: hexPass, Locked: "false", Permissions: []string{"user"}, Groups: []string{"default"}}, conf.DB.MongoDB, "accounts")
	if err != nil {
		return err
	}

	return utils.EmptyError()
}

func (accountManagerService) Search(secret string, token string, inputAccount types.Account, conf *config.Config, DB *database.DatabaseHandler) ([]types.Account, error) {

	_, err := utils.ValidateRequest(secret, token, "moderators", "", conf, DB)
	if err != nil {
		return []types.Account{}, err
	}

	input := utils.AccountToBson(inputAccount)
	results, err := DB.FindAll(input, conf.DB.MongoDB, "accounts")
	if err != nil {
		return []types.Account{}, errors.New("unauthorized request")
	}

	var output []types.Account
	for _, result := range results {
		converted := utils.BsonMapToAccount(result)
		converted.Token = ""
		converted.HashedPass = ""
		output = append(output, converted)
	}
	return output, utils.EmptyError()
}
