package main

import (
	"fmt"
	"github.com/yamamushi/kmud-2020/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if input.Map()["characters"] != nil {
		account.Characters = input.Map()["characters"].([]string)
	}
	return account
}
