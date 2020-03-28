package utils

import (
	"errors"
	"fmt"
	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/database"
	"github.com/yamamushi/kmud-2020/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func BsonMapToAccount(input bson.D) (account types.Account) {
	if input.Map()["username"] != nil {
		account.Username = input.Map()["username"].(string)
	}
	if input.Map()["hashedpass"] != nil {
		account.HashedPass = input.Map()["hashedpass"].(string)
	}
	if input.Map()["token"] != nil {
		account.Token = input.Map()["token"].(string)
	}
	if input.Map()["email"] != nil {
		account.Email = input.Map()["email"].(string)
	}
	if input.Map()["locked"] != nil {
		account.Locked = input.Map()["locked"].(string)
	}
	if input.Map()["permissions"] != nil {
		permissions := input.Map()["permissions"].(primitive.A)
		for _, permission := range permissions {
			account.Permissions = append(account.Permissions, fmt.Sprintf("%v", permission))
		}
	}
	if input.Map()["groups"] != nil {
		groups := input.Map()["groups"].(primitive.A)
		for _, group := range groups {
			account.Groups = append(account.Groups, fmt.Sprintf("%v", group))
		}
	}
	if input.Map()["characters"] != nil {
		characters := input.Map()["characters"].(primitive.A)
		for _, character := range characters {
			account.Characters = append(account.Characters, fmt.Sprintf("%v", character))
		}
	}
	return account
}

func AccountToBson(input types.Account) (output bson.M) {
	//output = bson.M{"username":"", "email":"", "token":"","hashedpass":"","locked":"","characters":primitive.A{""},"permissions":primitive.A{""}}
	output = make(bson.M)
	if input.Username != "" {
		output["username"] = input.Username
	}
	if input.Email != "" {
		output["email"] = input.Email
	}
	if input.Token != "" {
		output["token"] = input.Token
	}
	if input.HashedPass != "" {
		output["hashedpass"] = input.HashedPass
	}
	if input.Locked != "" {
		output["locked"] = input.Locked
	}
	if len(input.Characters) > 0 {
		characters := primitive.A{}
		for _, character := range input.Characters {
			characters = append(characters, character)
		}
		output["characters"] = characters
	}
	if len(input.Permissions) > 0 {
		permissions := primitive.A{}
		for _, permission := range input.Permissions {
			permissions = append(permissions, permission)
		}
		output["permissions"] = permissions
	}
	if len(input.Groups) > 0 {
		groups := primitive.A{}
		for _, group := range input.Groups {
			groups = append(groups, group)
		}
		output["groups"] = groups
	}
	return output
}

func ValidateRequest(secret string, token string, inputgroup string, inputpermission string, conf *config.Config, DB *database.DatabaseHandler) (account types.Account, err error) {

	if secret != conf.Crypt.AccountManagerSecret {
		return types.Account{}, errors.New("unauthorized request")
	}

	tokenFields := strings.Split(token, ":")
	if len(tokenFields) != 2 {
		return types.Account{}, errors.New("invalid token format")
	}

	result, err := DB.FindOne(bson.M{"username": tokenFields[0]}, conf.DB.MongoDB, "accounts")
	if err != nil {
		output := BsonMapToAccount(result)
		return output, errors.New("unauthorized request")
	}

	accountStruct := BsonMapToAccount(result)
	if accountStruct.Token != tokenFields[1] {
		return types.Account{}, errors.New("unauthorized request")
	}

	err = CheckAccountAccess(inputgroup, inputpermission, accountStruct)
	if err != nil {
		output := BsonMapToAccount(result)
		return output, err
	}

	return accountStruct, nil
}

func CheckGroup(inputgroup string, account types.Account) (err error) {
	groupaccess := false
	for _, group := range account.Groups {
		if group == inputgroup {
			groupaccess = true
		}
	}
	if inputgroup == "" {
		groupaccess = true
	}
	if !groupaccess {
		return errors.New("unauthorized request")
	}
	return nil
}

func CheckPermission(inputpermission string, account types.Account) (err error) {
	permissionsaccess := false
	for _, permission := range account.Permissions {
		if permission == inputpermission {
			permissionsaccess = true
		}
	}
	if inputpermission == "" {
		permissionsaccess = true
	}
	if !permissionsaccess {
		return errors.New("unauthorized request")
	}

	return nil
}

func CheckAccountAccess(inputgroup string, inputpermission string, account types.Account) (err error) {
	err = CheckGroup(inputgroup, account)
	if err != nil {
		return err
	}

	err = CheckPermission(inputpermission, account)
	if err != nil {
		return err
	}

	return nil
}
