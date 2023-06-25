package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	NonceLen int = 8
	letters      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GetRandStr(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func HmacSha256(source, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(source))
	signedBytes := h.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func GenAuthToken(id, secret string) string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := GetRandStr(NonceLen)
	source := id + "." + ts + "." + secret + "." + nonce + "." + secret
	sign := strings.ToLower(HmacSha256(source, secret))
	return id + "." + ts + "." + nonce + "." + sign
}

func AuthCheck(authorization, secret string) bool {
	strs := strings.Split(authorization, ".")
	if len(strs) != 4 {
		return false
	}
	id, ts, nonce, sign := strs[0], strs[1], strs[2], strs[3]
	source := id + "." + ts + "." + secret + "." + nonce + "." + secret

	servSign := strings.ToLower(HmacSha256(source, secret))
	return sign == servSign
}
