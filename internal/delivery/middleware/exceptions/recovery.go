package panicwall

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/exception"
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
				if r.Header.Get("Upgrade") == "websocket" {
					log.Fatal(rec)
					next.ServeHTTP(w, r)
					return
				}
				log.Println(rec)
				err, ok := rec.(exception.Exception)
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

func (recovery *PanicWall) handleRecoveredError(err *exception.Exception) ([]byte, int) {
	var code int
	switch err.Tag {
case exception.CONFLICT_ERROR:
		code = 409
	case exception.INTERNAL_ERROR:
		code = 500
	case exception.NOT_FOUND:
		code = 404
	case exception.BAD_REQUEST:
		code = 400
	case exception.LIMIT_EXCEED:
		code = 429
	case exception.UNAUTHORIZED:
		code = 401
	case exception.FORBIDDEN:
		code = 403
	case exception.UNPROCESSABLE:
		code = 422
	default:
		code = 418
	}

	json, _ := json.Marshal(err)
	return json, code
}
