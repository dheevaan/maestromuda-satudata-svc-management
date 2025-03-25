package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateMd5(text string) string {
	byteString := []byte(text)
	sum := md5.Sum(byteString)
	hashedString := hex.EncodeToString(sum[:])
	return hashedString
}
