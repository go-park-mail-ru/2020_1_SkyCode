package delivery

import (
	"github.com/2020_1_Skycode/internal/middlewares"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/2020_1_Skycode/internal/tools/CSRFManager"
	"github.com/2020_1_Skycode/internal/tools/requestValidator"
	"github.com/2020_1_Skycode/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

type UserHandler struct {
	userUseCase    users.UseCase
	sessionUseCase sessions.UseCase
	middlewareC    *middlewares.MWController
	v              *requestValidator.RequestValidator
	tM             *CSRFManager.CSRFManager
}

func NewUserHandler(private *gin.RouterGroup, public *gin.RouterGroup, uUC users.UseCase, sUC sessions.UseCase,
	validator *requestValidator.RequestValidator, tM *CSRFManager.CSRFManager,
	middlewareC *middlewares.MWController) *UserHandler {
	uh := &UserHandler{
		userUseCase:    uUC,
		sessionUseCase: sUC,
		middlewareC:    middlewareC,
		v:              validator,
		tM:             tM,
	}

	public.POST("/signup", uh.SignUp())
	public.GET("/profile", uh.GetProfile())

	private.PUT("/profile/bio", uh.EditBio())
	private.PUT("/profile/avatar", uh.EditAvatar())
	private.PUT("/profile/password", uh.ChangePassword())
	private.PUT("/profile/phone", uh.ChangePhoneNumber())

	return uh
}

type signUpRequest struct {
	FirstName string `json:"firstName, omitempty" binding:"required" validate:"min=2"`
	LastName  string `json:"lastName, omitempty" binding:"required" validate:"min=2"`
	Phone     string `json:"phone, omitempty" binding:"required" validate:"min=11,max=15"`
	Password  string `json:"password, omitempty" binding:"required" validate:"passwd"`
}

type editBioRequest struct {
	FirstName string `json:"firstName" binding:"required" validate:"min=2"`
	LastName  string `json:"lastName" binding:"required" validate:"min=2"`
	Email     string `json:"email" validate:"email"`
}

type changePasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required" validate:"passwd"`
}

type changePhoneNumberRequest struct {
	NewPhone string `json:"newPhone" binding:"required" validate:"min=11,max=15"`
}

//@Tags User
//@Summary Sign Up Route
//@Description Signing up
//@Accept json
//@Produce json
//@Param "SignUpReq" body signUpRequest true "New user data"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Router /signup [post]
func (uh *UserHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &signUpRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		errorsList := uh.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.ErrorRequestValidation.Error(),
			})

			return
		}

		_, err := uh.userUseCase.GetUserByPhone(req.Phone)

		if err == nil {
			logrus.Info(tools.UserExists)
			c.JSON(http.StatusOK, tools.Error{
				ErrorMessage: tools.UserExists.Error(),
			})

			return
		}

		u := &models.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
			Password:  req.Password,
		}

		if err := uh.userUseCase.CreateUser(u); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		session, cookie := models.GenerateSession(u.ID)

		if err := uh.sessionUseCase.StoreSession(session); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusNotFound, tools.Error{
				ErrorMessage: tools.SessionStoreError.Error(),
			})

			return
		}

		csrfToken, err := uh.tM.GenerateCSRF(session.UserId, session.Token)

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.Writer.Header().Set("X-Csrf-Token", csrfToken)

		c.SetCookie(cookie.Name,
			cookie.Value,
			cookie.MaxAge,
			cookie.Path,
			cookie.Domain,
			cookie.Secure,
			cookie.HttpOnly)

		c.JSON(http.StatusOK, tools.Message{
			Message: "User has been registered",
		})
	}
}

//@Tags User
//@Summary Edit Bio Route
//@Description Editing bio data of user
//@Accept json
//@Produce json
//@Param "bioReq" body editBioRequest true "Bio data of user"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /profile/bio [put]
func (uh *UserHandler) EditBio() gin.HandlerFunc {
	return func(c *gin.Context) {
		updProfile := &editBioRequest{}

		if err := c.Bind(updProfile); err != nil {
			logrus.Info("Bad params")
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := uh.v.ValidateRequest(updProfile)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		user, err := uh.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		user.FirstName = updProfile.FirstName
		user.LastName = updProfile.LastName
		user.Email = updProfile.Email

		if err := uh.userUseCase.UpdateBio(user); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

//@Tags User
//@Summary Edit Avatar Route
//@Description Changing Avatar Of User
//@Accept mpfd
//@Produce json
//@Param "image" formData file true "New avatar of user"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /profile/avatar [put]
func (uh *UserHandler) EditAvatar() gin.HandlerFunc {
	rootDir, _ := os.Getwd()

	return func(c *gin.Context) {
		user, err := uh.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		file, err := c.FormFile("avatar")

		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		filename := shortuuid.New()

		if err := c.SaveUploadedFile(file, filepath.Join(rootDir, tools.AvatarPath, filename)); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		if user.Avatar != "" {
			if err := os.Remove(filepath.Join(rootDir, tools.AvatarPath, user.Avatar)); err != nil {
				logrus.Info(err)
				c.JSON(http.StatusInternalServerError, tools.Error{
					ErrorMessage: tools.DeleteAvatarError.Error(),
				})

				return
			}
		}

		if err := uh.userUseCase.UpdateAvatar(user.ID, filename); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

//@Tags User
//@Summary Change Phone Number Route
//@Description Changing Phone Number Of User
//@Accept json
//@Produce json
//@Param "phone" body changePhoneNumberRequest true "New phone number of user"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /profile/phone [put]
func (uh *UserHandler) ChangePhoneNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &changePhoneNumberRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := uh.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		user, err := uh.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		validate := validator.New()

		if err := validate.Struct(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if err := uh.userUseCase.UpdatePhoneNumber(user.ID, req.NewPhone); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UpdatePhoneError.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

//@Tags User
//@Summary Change Password Route
//@Description Changing password of user
//@Accept json
//@Produce json
//@Param "password" body changePasswordRequest true "New password"
//@Success 200 object tools.Message
//@Failure 400 object tools.Error
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /profile/password [put]
func (uh *UserHandler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &changePasswordRequest{}

		if err := c.Bind(req); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.BadRequest.Error(),
			})

			return
		}

		errorsList := uh.v.ValidateRequest(req)

		if len(*errorsList) > 0 {
			logrus.Info(tools.NotRequiredFields)
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: tools.NotRequiredFields.Error(),
			})

			return
		}

		user, err := uh.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		if err := uh.userUseCase.UpdatePassword(user.ID, req.NewPassword); err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, tools.Error{
				ErrorMessage: tools.UpdatePhoneError.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.Message{
			Message: "success",
		})
	}
}

//@Tags User
//@Summary Get Profile Route
//@Description Getting Profile Of User
//@Accept json
//@Produce json
//@Success 200 object models.User
//@Failure 401 object tools.Error
//@Failure 500 object tools.Error
//@Security basicAuth
//@Router /profile [get]
func (uh *UserHandler) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := uh.middlewareC.GetUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, tools.Error{
				ErrorMessage: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, tools.UserMessage{
			User: user,
		})
	}
}
