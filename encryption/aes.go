package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptAES(key, text []byte) ([]byte, error) {
	// add padding if len() % 16 != 0
	if mod := len(text) % aes.BlockSize; mod != 0 {
		padding := make([]byte, aes.BlockSize-mod)
		//add zero padding
		text = append(text, padding...)
	}

	// required key length is 16 or 24 or 32
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cbEncrypt := cipher.NewCBCEncrypter(block, iv)
	cbEncrypt.CryptBlocks(ciphertext[aes.BlockSize:], text)
	return ciphertext, nil
}

func DecryptAES(key, text []byte) ([]byte, error) {
	if len(text)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("[err] invalid ciphertext")
	}

	// required key length is 16 or 24 or 32
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := text[:aes.BlockSize]
	ciphertext := text[aes.BlockSize:]
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	// delete zero padding
	return bytes.Trim(plaintext, "\x00"), nil
}
