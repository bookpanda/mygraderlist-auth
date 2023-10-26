package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
)

func Encrypt(secret []byte, content string) (encoded string, err error) {
	plainText := []byte(content)

	block, err := aes.NewCipher(secret)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(secret []byte, secure string) (decoded string, err error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	if err != nil {
		return
	}

	block, err := aes.NewCipher(secret)

	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("invalid ciphertext")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

func Hash(bv []byte) string {
	h := sha256.New()
	h.Write(bv)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
