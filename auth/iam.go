package auth

import (
	"context"
	"io/ioutil"
	"strings"

	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"
)

type iamCredentials struct {
	ycCreds ycsdk.NonExchangeableCredentials
}

func (i iamCredentials) YandexCloudAPICredentials() {
}

func (i iamCredentials) Token(ctx context.Context) (string, error) {
	tokenRes, err := creds.ycCreds.IAMToken(ctx)
	if err != nil {
		return "", err
	}
	return tokenRes.IamToken, nil
}

var creds *iamCredentials

func init() {
	data, err := ioutil.ReadFile("iam.txt")
	if err == nil {
		token := strings.TrimSpace(string(data))
		if token != "" {
			creds = &iamCredentials{
				ycsdk.InstanceServiceAccount(),
			}
			return
		}
	}
	creds = &iamCredentials{
		ycsdk.InstanceServiceAccount(),
	}
}

func GetYdbCredentials() credentials.Credentials {
	return creds
}

func GetIAMToken(ctx context.Context) (string, error) {
	return creds.Token(ctx)
}
