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
	CheckoutOrderError = errors.New("Error checkout order")
	RestaurantNotFoundError = errors.New("Restaurant not found")
	RestaurantPermissionsError = errors.New("You doesn't have permissions to manage this restaurant")
	ExpiredCSRFError = errors.New("Expired CSRF token")
	WrongCSRFtoken = errors.New("Wrong CSRF token")
	CSRFNotPresented = errors.New("CSRF not presented")
	HashingError = errors.New("Hashing error")
	Authorized = errors.New("Already authorized")
)
