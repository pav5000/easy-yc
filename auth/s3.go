package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	region   = "ru-central1"
	endpoint = "https://storage.yandexcloud.net"
)

func GetS3Service() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %s", err.Error())
	}
	return s3.New(sess, aws.NewConfig()), nil
}

// S3CredsFromConfigToEnv помощник для локальной разработки
// Вытаскивает ключи стореджа из ~/.aws/credentials
// И кладет в ENV-переменные, чтобы потом передавать их создаваемым машинам
func S3CredsFromConfigToEnv() error {
	filename := filepath.Join(os.Getenv("HOME"), ".aws/credentials")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	parseParam := func(line string, paramName string) (string, error) {
		components := strings.Split(line, "=")
		if len(components) != 2 {
			return "", errors.New("wrong number of components")
		}
		return strings.TrimSpace(components[1]), nil
	}

	var keyId, secretKey string

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "aws_secret_access_key") {
			secretKey, err = parseParam(line, "aws_secret_access_key")
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "aws_access_key_id") {
			keyId, err = parseParam(line, "aws_access_key_id")
			if err != nil {
				return err
			}
		}
	}

	if keyId == "" {
		return errors.New("empty key id")
	}
	if secretKey == "" {
		return errors.New("empty secret key")
	}

	err = os.Setenv("AWS_ACCESS_KEY_ID", keyId)
	if err != nil {
		return err
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)
	if err != nil {
		return err
	}

	return nil
}
