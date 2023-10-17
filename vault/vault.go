// Package vault contains helper methods to store key vaulue in DB in encrypted format
package vault

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"olx-clone/conf"
	"olx-clone/constants"
	"olx-clone/functions/logger"
	"olx-clone/infra/redis"
	"olx-clone/models/vault"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

var log = logger.Log

func GetFromRedis(key string) (vault.Struct, error) {
	var keyValue vault.Struct
	redisKey := key + constants.VaultKeySuffix
	theValueString, err := redis.Get(redisKey)
	if err != nil {
		return keyValue, err
	}
	decoder := json.NewDecoder(strings.NewReader(theValueString))
	err = decoder.Decode(&keyValue)
	if err != nil {
		return keyValue, err
	}
	return keyValue, nil
}

func SetInRedis(key string, value string) error {
	keyValue := vault.Struct{
		Key:   key,
		Value: value,
	}
	redisKey := key + constants.VaultKeySuffix
	return redis.SetStruct(redisKey, keyValue, time.Minute*10)
}

func Get(context context.Context, key string) (string, error) {
	keyValueData, err := GetFromRedis(key)
	if err != nil {
		keyValueData, err = vault.Get(context, key)
		if err != nil {
			return "", err
		}
		SetInRedis(key, keyValueData.Value)
	}
	encryptedData, err := base64.StdEncoding.DecodeString(keyValueData.Value)
	if err != nil {
		return "", err
	}
	decryptedString, err := AESDecrypt(string(encryptedData))
	if err != nil {
		return "", err
	}
	return decryptedString, nil
}

func GetOrDefaultString(context context.Context, key string, defaultString string) string {
	value, err := Get(context, key)
	if err != nil {
		return defaultString
	}
	return value
}

func Save(context context.Context, key string, value string) (*vault.Struct, error) {
	encrypted, err := AESEncryptToBase64(value)
	if err != nil {
		return nil, err
	}
	keyValueData, err := vault.Create(context, key, encrypted, nil)
	// Update in redis for further getting from redis
	SetInRedis(key, encrypted)
	return keyValueData, err
}

func ValidateKey() {
	// Get the current key passed from Environment variable
	vaultKeyByte, err := base64.StdEncoding.DecodeString(conf.VaultKey)
	if err != nil {
		panic("error reading vault key")
	}
	vaultKey := string(vaultKeyByte)

	key, err := Get(context.TODO(), KeyVaultEncryptionKey)

	// If we don't have key in Vault we'll save the key
	if err == sql.ErrNoRows {
		// Save the private key in vault
		Save(context.TODO(), KeyVaultEncryptionKey, vaultKey)
		log.Println("saved the key in vault")
		return
	}

	// If we face some error we'll panic
	if err != nil {
		sentry.CaptureException(errors.Wrapf(err, "error validating vault key"))
		log.Errorln("error validating vault key")
		panic(err)
	}

	// If the keys found does not matches with the key passed
	if vaultKey != key {
		panic("invalid vault key passed")
	}

	log.Println("vault key validated")
}
