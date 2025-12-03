package panicwall

import (
	"encoding/json"
	"log"
	"net/http"

	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
)

type PanicWall struct {
}

func NewPanicWall() *PanicWall {
	return &PanicWall{}
}

func (recovery *PanicWall) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println(rec)
				err, ok := rec.(derror.DomainError)
				if ok {
					json, status := recovery.handleRecoveredError(&err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(status)
					w.Write(json)
				} else {
					w.WriteHeader(501)
					errorResponse := map[string]any{
						"error": rec,
					}
					jsonResponse, _ := json.Marshal(errorResponse)
					w.Write(jsonResponse)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (recovery *PanicWall) handleRecoveredError(err *derror.DomainError) ([]byte, int) {
	var code int
	switch {
	case derror.IsType(err, derror.ErrTypeConflict):
		code = http.StatusConflict
	case derror.IsType(err, derror.ErrTypeBadRequest):
		code = http.StatusBadRequest
	case derror.IsType(err, derror.ErrTypeUnauthorized):
		code = http.StatusUnauthorized
	case derror.IsType(err, derror.ErrTypeDB):
		code = http.StatusInternalServerError
	case derror.IsType(err, derror.ErrTypeForeignKey):
		code = http.StatusBadRequest
	case derror.IsType(err, derror.ErrTypeTimeout):
		code = http.StatusGatewayTimeout
	case derror.IsType(err, derror.ErrTypeInternal):
		code = http.StatusInternalServerError
	case derror.IsType(err, derror.ErrTypeNotFound):
		code = http.StatusNotFound
	case derror.IsType(err, derror.ErrTypeForbidden):
		code = http.StatusForbidden
	case derror.IsType(err, derror.ErrTypeUnprocessable):
		code = http.StatusUnprocessableEntity
	default:
		code = 418
	}

	json, _ := json.Marshal(map[string]any{
		"message": err.Msg,		
	})
	return json, code
}
