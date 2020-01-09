package Cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

type AES256Cryptor struct {
	Salt       []byte
	Iv         []byte
	Passphrase string
}

func NewAES256Cryptor(text string) (*AES256Cryptor, error) {

	aes256 := &AES256Cryptor{}
	aes256.Iv = []byte{0xD8, 0xAA, 0x29,
		0x82, 0x81, 0x92, 0x68, 0xD4,
		0x17, 0x0E, 0x5A, 0x40, 0x48,
		0x2B, 0x9C, 0x4E}
	aes256.Passphrase = text
	aes256.Salt = make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, aes256.Salt)
	if err != nil {
		// fallback
		aes256.Salt = []byte{0x3B, 0xC1, 0xFD, 0x23, 0x1D, 0xBF, 0xF8, 0x70}
		e1 := errors.New("Randomize error")
		return aes256, e1
	}

	return aes256, nil
}

func Decrypt(cipherbytes []byte, aes256Cryptor *AES256Cryptor) []byte {

	key := deriveKey(aes256Cryptor)

	b, _ := aes.NewCipher(key)

	aescfb := cipher.NewCFBDecrypter(b, aes256Cryptor.Iv)
	decrypted := make([]byte, len(cipherbytes))

	aescfb.XORKeyStream(decrypted, cipherbytes)
	return pkcs5Trimming(decrypted)
}

func deriveKey(aes256Cryptor *AES256Cryptor) []byte {
	return pbkdf2.Key([]byte(aes256Cryptor.Passphrase), aes256Cryptor.Salt, 19, 32, sha1.New)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
