package auth

import (
	"context"
	"sync"

	ycsdk "github.com/yandex-cloud/go-sdk"
)

var (
	sdk     *ycsdk.SDK
	sdkLock sync.Mutex
)

func GetSDK() (*ycsdk.SDK, error) {
	sdkLock.Lock()
	defer sdkLock.Unlock()

	if sdk == nil {
		var err error
		sdk, err = ycsdk.Build(context.Background(), ycsdk.Config{
			Credentials: creds,
		})
		if err != nil {
			return nil, err
		}
	}

	return sdk, nil
}
