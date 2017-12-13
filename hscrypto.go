package hld

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log"
)

//PKCS5Padding pack PKCS5
func PKCS5Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := blockSize - (srcLen % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

//PKCS5Unpadding unpack PKCS5
func PKCS5Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, errors.New("")
	}
	return src[:srcLen-paddingLen], nil
}

//GenerateRandomBytes generate random bytes for crypto usage
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

//GenerateAESKeys generate random key/iv for AES
func GenerateAESKeys(strength int) ([]byte, []byte) {
	if (strength != 32) && (strength != 16) {
		log.Fatal("AES Strength must be either 16 or 32")
	}
	bytes, _ := GenerateRandomBytes(64)
	sha := sha256.Sum256(bytes)

	key := make([]byte, strength)
	copy(key, sha[0:strength])
	bytes, _ = GenerateRandomBytes(64)
	sha = sha256.Sum256(bytes)
	iv := make([]byte, 16)
	copy(iv, sha[0:16])
	return key, iv
}

//EncryptAES one time encryption
func EncryptAES(key []byte, iv []byte, data []byte) ([]byte, error) {
	var err error
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	data = PKCS5Padding(data, aes.BlockSize)
	ciphertext := make([]byte, len(data))
	mode.CryptBlocks(ciphertext, data)
	return ciphertext, nil
}

//DecryptAES one time decryption
func DecryptAES(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {
	var err error

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	data := make([]byte, len(ciphertext))
	mode.CryptBlocks(data, ciphertext)

	padded, err := PKCS5Unpadding(data, aes.BlockSize)

	if err == nil {
		return padded, nil
	}
	return nil, err
}
