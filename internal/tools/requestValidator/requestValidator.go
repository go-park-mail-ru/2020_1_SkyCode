package requestValidator

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"strconv"
)

type RequestValidator struct {
	validate   *validator.Validate
	translator *ut.Translator
}

func NewRequestValidator() *RequestValidator {
	v := validator.New()

	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	_ = v.RegisterTranslation("passwd", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} is not strong enough", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})

	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	return &RequestValidator{
		validate:   v,
		translator: &trans,
	}
}

func (v *RequestValidator) ValidateRequest(data interface{}) *map[string]string {
	errorsMap := map[string]string{}
	err := v.validate.Struct(data)

	if err != nil {
		fmt.Println(err)
		for _, err := range err.(validator.ValidationErrors) {
			if err.Type() == reflect.TypeOf(0.1) {
				if err.Param() != strconv.FormatFloat(err.Value().(float64), 'f', -1, 64) {
					errorsMap[err.Field()] = err.Translate(*v.translator)
				}

				continue
			}
			if err.Param() != err.Value().(string) {
				errorsMap[err.Field()] = err.Translate(*v.translator)
			}
		}
	}

	return &errorsMap
}
