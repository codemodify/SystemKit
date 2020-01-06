package Helpers

import (
	"encoding/base64"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// IsNullOrEmpty -
func IsNullOrEmpty(value string) bool {
	return (len(strings.TrimSpace(value)) <= 0)
}

const alphaAndDigitsBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString2(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = alphaAndDigitsBytes[rand.Intn(len(alphaAndDigitsBytes))]
	}
	return string(b)
}

// RandomString -
func RandomString(length int) string {
	// head -c 1000 /dev/urandom | base64 | tr -cd '[:alnum:]' | cut -c 1-100

	buffer := make([]byte, 1000*length)
	_, err := rand.Read(buffer)
	if err != nil {
		return randomString2(length)
	}

	bufferAsB64String := base64.StdEncoding.EncodeToString(buffer)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return randomString2(length)
	}

	theWholeThing := reg.ReplaceAllString(bufferAsB64String, "")
	if len(theWholeThing) > length {
		return theWholeThing[:length]
	}

	theWholeThing = theWholeThing + randomString2(length)

	return theWholeThing[:length]
}

// Contains -
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}
