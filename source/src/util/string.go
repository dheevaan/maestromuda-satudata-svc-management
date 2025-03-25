package util

import (
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func GetUpperAlphabetOf(number int32) string {
	upperAlphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	return upperAlphabet[number-1]
}
func GenerateID() (res string) {
	res = strings.ReplaceAll(uuid.New().String(), "-", "")
	return
}
func SetAlias(alias string) (res string) {
	res = slug.Make(alias)
	return
}
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
