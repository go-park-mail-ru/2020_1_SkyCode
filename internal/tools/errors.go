package tools

import "errors"

var (
	BadRequest                 = errors.New("Bad request")
	NoSuchUser                 = errors.New("User with the same name does not exist")
	WrongPassword              = errors.New("Wrong Password")
	NotFound                   = errors.New("Not found")
	SessionStoreError          = errors.New("Store session error")
	NoSuchSession              = errors.New("Session cookie not found")
	DeleteSessionError         = errors.New("Error while deleting session")
	Unauthorized               = errors.New("Unauthorized")
	SessionTypeAssertionErr    = errors.New("Error assert to session type")
	UserTypeAssertionErr       = errors.New("Error assert to user type")
	UpdatePhoneError           = errors.New("Error updating phone number")
	DeleteAvatarError          = errors.New("Error deleting avatar")
	CheckoutOrderError         = errors.New("Error checkout order")
	RestaurantNotFoundError    = errors.New("Restaurant not found")
	RestaurantPermissionsError = errors.New("You doesn't have permissions to manage this restaurant")
	ExpiredCSRFError           = errors.New("Expired CSRF token")
	WrongCSRFtoken             = errors.New("Wrong CSRF token")
	CSRFNotPresented           = errors.New("CSRF not presented")
	HashingError               = errors.New("Hashing error")
	Authorized                 = errors.New("Already authorized")
	NotRequiredFields          = errors.New("No such required fields")
	ErrorRequestValidation     = errors.New("Request validation failed")
	UserExists                 = errors.New("Phone number already exists")
	GetOrdersError             = errors.New("Get orders error")
	ReviewAlreadyExists        = errors.New("Review to this restaurant by this user already exists")
	BadQueryParams             = errors.New("Bad query params")
	ReviewNotFoundError        = errors.New("Review not found")
	PermissionError            = errors.New("You don't have rights for it")
	ApiResponseStatusNotOK     = errors.New("Api answer error")
	ApiAnswerEmptyResult       = errors.New("Api return empty result")
	RestPointNotFound          = errors.New("Restaurant point not found")
	ApiMultiAnswerError        = errors.New("Cant correctly identufy address")
	ApiNotHouseAnswerError     = errors.New("Object not a house")
	GRPCOpertionNotSuccess     = errors.New("GRPC answer not true, operation not success")
	ProductNotFoundError       = errors.New("Product not found")
	BindingError               = errors.New("Error while binding request")
	RestTagNotFound            = errors.New("Restaurant tag not found")
)
