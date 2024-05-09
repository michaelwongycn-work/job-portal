package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

var secretKey = ""

func SetAuthConfig(key string) {
	secretKey = key
}

func Hash(input string) (string, error) {
	concatenated := input + secretKey
	hash := sha256.New()

	_, err := hash.Write([]byte(concatenated))
	if err != nil {
		return "", err
	}

	hashed := hash.Sum(nil)

	hashedString := hex.EncodeToString(hashed)

	return hashedString, nil
}

func Encrypt(input string) (string, error) {
	key, _ := base64.URLEncoding.DecodeString(secretKey)
	text := []byte(input)

	hashed_key := sha256.Sum256(key)
	block, err := aes.NewCipher(hashed_key[:])
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	ciphertextString := base64.StdEncoding.EncodeToString(ciphertext)
	return base64.URLEncoding.EncodeToString([]byte(ciphertextString)), nil
}

func Decrypt(input string) (string, error) {
	key, _ := base64.URLEncoding.DecodeString(secretKey)
	text, _ := base64.URLEncoding.DecodeString(input)
	text, _ = base64.StdEncoding.DecodeString(string(text))

	hashed_key := sha256.Sum256(key)
	block, err := aes.NewCipher(hashed_key[:])
	if err != nil {
		return "", err
	}
	if len(text) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	return string(text), nil
}
