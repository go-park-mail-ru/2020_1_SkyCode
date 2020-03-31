package tools

import "errors"

var (
	BadRequest = errors.New("Bad request")
	NoSuchUser = errors.New("User with the same name does not exist")
	WrongPassword = errors.New("Wrong Password")
	NotFound = errors.New("Not found")
	SessionStoreError = errors.New("Store session error")
	NoSuchSession = errors.New("Session cookie not found")
	DeleteSessionError = errors.New("Error while deleting session")
	Unauthorized = errors.New("Unauthorized")
	SessionTypeAssertionErr = errors.New("Error assert to session type")
	UserTypeAssertionErr = errors.New("Error assert to user type")
	UpdatePhoneError = errors.New("Error updating phone number")
	DeleteAvatarError = errors.New("Error deleting avatar")
)
