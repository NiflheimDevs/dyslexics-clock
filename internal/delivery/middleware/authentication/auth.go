package midauth

import (
	"context"
	"net/http"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	repository "github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository/postgres"
)

type Authentication struct {
	DeviceRepo repository.DeviceRepo
	Constants  *bootstrap.Constants
	JWTService service.JWT
}

func NewAuth(
	deviceRepo repository.DeviceRepo,
	Constants *bootstrap.Constants,
	JWTService service.JWT,
) *Authentication {
	return &Authentication{
		DeviceRepo: deviceRepo,
		Constants:  Constants,
		JWTService: JWTService,
	}
}

func (am *Authentication) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var deviceID uint
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			// deviceID = -2
			panic(derror.New(derror.ErrTypeUnauthorized, "invalid token", nil))
		} else {
			tokenString := authHeader[7:]
			claims, err := am.JWTService.VerifyToken(tokenString)
			if err != nil {
				// deviceID = -2
				panic(derror.New(derror.ErrTypeUnauthorized, "invalid token", nil))
			} else if claims == nil {
				// deviceID = -1
				panic(derror.New(derror.ErrTypeUnauthorized, "invalid token", nil))
			} else {
				deviceID = uint(claims["sub"].(float64))
			}
		}
		_, err := am.DeviceRepo.GetDeviceById(r.Context(), uint(deviceID))
		if err != nil {
			// deviceID = -1
			panic(derror.New(derror.ErrTypeUnauthorized, "invalid token", nil))
		}

		ctx := context.WithValue(r.Context(), am.Constants.Context.DeviceID, deviceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
