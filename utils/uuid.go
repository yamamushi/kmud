package utils

import (
	"encoding/base64"
	"github.com/gofrs/uuid"
	"github.com/yamamushi/kmud-2020/crypt"
)

// GetUUID function
func GetUUID() (id string, err error) {

	formattedid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return formattedid.String(), nil

}

func GetRandomToken() (token string, err error) {
	uuid, err := GetUUID()
	if err != nil {
		return "", err
	}

	sha := crypt.Sha256Sum(uuid)
	baseEncoded := base64.URLEncoding.EncodeToString(sha)
	baseEncoded = RemoveSpecial(baseEncoded)
	return baseEncoded, nil
}
