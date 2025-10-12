package conf

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

const ApiKeyFile = "./data/apikey.txt"

var ApiKeys = make(map[string]bool)

func AddAPIKey(key string) {
	ApiKeys[key] = true
}

func IsValidAPIKey(key string) bool {
	return ApiKeys[key]
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 4) // 4字节 = 8位hex字符串
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func InitAPIKey() string {
	data, err := os.ReadFile(ApiKeyFile)
	if err == nil {
		key := string(data)
		AddAPIKey(key)
		return key
	}

	key, err := GenerateAPIKey()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(ApiKeyFile, []byte(key), 0644)
	if err != nil {
		panic(err)
	}

	AddAPIKey(key)
	return key
}
