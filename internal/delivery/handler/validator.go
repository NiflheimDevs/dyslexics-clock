package handler

import (
	"encoding/json"
	"log"
	"net/http"

	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
	"github.com/go-playground/validator/v10"
)

func Validated[T any](validate pkg.Validator, r *http.Request) T {
	var params T
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println(err)
		panic(derror.New(derror.ErrTypeBadRequest, "invalid request body", err))
	}

	if err := validate.Struct(params); err != nil {
		log.Println("ValidationError:", err)
		panic(derror.New(derror.ErrTypeBadRequest, "missing parameter", err))

	}

	return params
}

func StructValidator(validate *validator.Validate, s any) error {
	return validate.Struct(s)

}
