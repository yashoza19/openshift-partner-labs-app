package actions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/gobuffalo/envy"
	"io"
	"log"
)

func encryptVal(valToEncrypt string) []byte {
	byteOfVal := []byte(valToEncrypt)
	key := []byte(envy.Get("X_ENCRYPTION_KEY", ""))
	cypher, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
	}
	gcm, err := cipher.NewGCM(cypher)
	if err != nil {
		log.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err)
	}

	return gcm.Seal(nonce, nonce, byteOfVal, nil)
}
