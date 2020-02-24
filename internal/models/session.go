package models

import "math/rand"

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)
const CookieSessionName = "SCDSESSIONID"

func GenerateSessionCookie() string {
	byteSlice := make([]rune, 64)
	for i := range byteSlice {
		byteSlice[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(byteSlice)
}