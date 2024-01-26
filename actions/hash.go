package actions

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func dohash(valtohash []byte) string {
	hashedval, err := bcrypt.GenerateFromPassword(valtohash, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hashedval)
}

func docompare(hashedval string, plainval []byte) bool {
	byteOfHashedVal := []byte(hashedval)
	err := bcrypt.CompareHashAndPassword(byteOfHashedVal, plainval)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
