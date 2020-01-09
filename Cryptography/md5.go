package Cryptography

import (
	"crypto/md5"
	"encoding/hex"
)

// StringToMD5 -
func StringToMD5(value string) string {
	hasher := md5.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}
