package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"olx-clone/conf"
)

var vaultKey string

func init() {
	// readKeys()
}

func readKeys() {
	normal, err := base64.StdEncoding.DecodeString(conf.VaultKey)
	if err != nil {
		panic("error reading vault key")
	}
	vaultKey = string(normal)
	// 32 bit key for AES256
	if len(vaultKey) != 32 {
		panic("invalid vault key.")
	}
}

// AESEncrypt does AES Encryption and return encrypted string
func AESEncrypt(text string) (string, error) {
	textByte := []byte(text)
	key := []byte(vaultKey)
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	enc := gcm.Seal(nonce, nonce, textByte, nil)
	return string(enc), nil

}
func AESEncryptToBase64(text string) (string, error) {
	encrypted, err := AESEncrypt(text)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(encrypted))
	return encoded, nil
}

// AESDecrypt does decryption of AES encrypted string an returns plain string
func AESDecrypt(encrypted string) (string, error) {
	key := []byte(vaultKey)
	ciphertext := []byte(encrypted)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
func AESDecryptBase64(text string) (string, error) {
	normal, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	plaintext, err := AESDecrypt(string(normal))
	if err != nil {
		return "", err
	}
	return plaintext, nil
}
