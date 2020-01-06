package ssh

import (
	"encoding/base64"
)

func PublicKeyAsBytesToString(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func PublicKeyToBytes(key PublicKey) []byte {
	return key.Marshal()
}

func PublicKeyToString(key PublicKey) string {
	return PublicKeyAsBytesToString(PublicKeyToBytes(key))
}
