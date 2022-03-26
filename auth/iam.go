package auth

import (
	"context"
	"io/ioutil"
	"strings"

	ycsdk "github.com/yandex-cloud/go-sdk"
)

var creds ycsdk.NonExchangeableCredentials

func init() {
	data, err := ioutil.ReadFile("iam.txt")
	if err == nil {
		token := strings.TrimSpace(string(data))
		if token != "" {
			creds = ycsdk.NewIAMTokenCredentials(token)
			return
		}
	}
	creds = ycsdk.InstanceServiceAccount()
}

func GetIAMToken(ctx context.Context) (string, error) {
	tokenRes, err := creds.IAMToken(ctx)
	if err != nil {
		return "", err
	}
	return tokenRes.IamToken, nil
}
