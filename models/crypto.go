package models

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"encoding/base64"
	"errors"
)
const (
	key string = "239a1c1c-edfc-40"
)

func Encrypt(origData string) (string, error) {
	if origData == "" {
		return "", errors.New("nothing to encrypt")
	}
	crypted, err := aesEncrypt([]byte(origData), []byte(key))
	if err != nil {
		return "", err
	}
	encodeString := base64.StdEncoding.EncodeToString(crypted)
	return encodeString, nil

}

func Decrypt(crypted string) (string, error) {
	if crypted == "" {
		return "", errors.New("nothing to decrypt")
	}
	decodeString, _ := base64.StdEncoding.DecodeString(crypted)
	d, err := aesDecrypt([]byte(decodeString), []byte(key))
	if err != nil {
		return "", err
	}
	return string(d), nil


}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		// crypto/cipher: input not full blocks
		return nil, errors.New("fake encrypted data")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}


func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}