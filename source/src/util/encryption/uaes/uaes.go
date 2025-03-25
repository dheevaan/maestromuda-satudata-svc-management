package uaes

import (
	"encoding/json"
	"log"

	"github.com/Luzifer/go-openssl/v4"
)

type AES struct {
	secretKey string
	cg        openssl.CredsGenerator
}

// Constructor for AES
func NewAES(secretKey string) *AES {
	o := new(AES)
	o.secretKey = secretKey

	o.cg = openssl.BytesToKeyMD5 //? Compatible with Crypto JS
	return o
}

func (o *AES) Decrypt(enc string) (res string) {
	dec, err := openssl.New().DecryptBytes(o.secretKey, []byte(enc), o.cg)
	if err != nil {
		log.Println(err, enc)
		return enc //? Return original
	}

	res = string(dec)
	return
}
func (o *AES) Decrypt_ToMap(enc string) (res map[string]any) {
	if err := json.Unmarshal([]byte(o.Decrypt(enc)), &res); err != nil {
		log.Println(err)
		return
	}
	return
}
func (o *AES) Encrypt(target []byte) (res string, err error) {
	resRaw, err := openssl.New().EncryptBytes(o.secretKey, target, o.cg)
	if err != nil {
		log.Println(err)
		return
	}
	res = string(resRaw)
	return
}
func (o *AES) Encrypt_Any(target any) (res string, err error) {
	asJson, err := json.Marshal(target)
	if err != nil {
		log.Println(err)
		return
	}

	resRaw, err := openssl.New().EncryptBytes(o.secretKey, asJson, o.cg)
	if err != nil {
		log.Println(err)
		return
	}
	res = string(resRaw)
	return
}
