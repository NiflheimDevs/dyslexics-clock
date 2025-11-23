package serviceimpl

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/exception"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJWT(Const *bootstrap.Constants) *JWT {
	return &JWT{
		PrivateKey: loadPrivateKey(Const.JWTKeysPath + "/privateKey.pem"),
		PublicKey:  loadPublicKey(Const.JWTKeysPath + "/publicKey.pem"),
	}
}

func loadPrivateKey(keyPath string) *rsa.PrivateKey {
	privateKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func loadPublicKey(keyPath string) *rsa.PublicKey {
	publicKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func (j *JWT) GenerateToken(userID int) (string, string) {
	accessTokenClaims := jwt.MapClaims{
		"iss": "bidlancer",
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(j.PrivateKey)
	if err != nil {
		panic(exception.Exception{
			Tag: exception.INTERNAL_ERROR,
			Errors: []exception.SpecificError{
				exception.AUTH_GENERATE_TOKEN_ERROR,
			},
		})
	}

	refreshTokenClaims := jwt.MapClaims{
		"iss": "bidlancer",
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(j.PrivateKey)
	if err != nil {
		panic(exception.Exception{
			Tag: exception.INTERNAL_ERROR,
			Errors: []exception.SpecificError{
				exception.AUTH_GENERATE_TOKEN_ERROR,
			},
		})
	}

	return accessTokenString, refreshTokenString
}

func (j *JWT) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return j.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("saman says he doesn't like this")
	}
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, nil
		}
	}
	return claims, nil

}

func (j *JWT) RefreshToken(tokenString string) string {
	claims, err := j.VerifyToken(tokenString)
	if claims == nil || err != nil {
		panic(exception.Exception{
			Tag: exception.UNAUTHORIZED,
			Errors: []exception.SpecificError{
				exception.AUTH_TOKEN_EXPIRED,
			},
		})
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		panic(exception.Exception{
			Tag: exception.NOT_FOUND,
			Errors: []exception.SpecificError{
				exception.USER_NOT_FOUND,
			},
		})
	}

	newAccessTokenClaims := jwt.MapClaims{
		"iss": "bidlancer",
		"sub": int(userID),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, newAccessTokenClaims)
	newAccessTokenString, err := newAccessToken.SignedString(j.PrivateKey)
	if err != nil {
		panic(exception.Exception{
			Tag: exception.UNPROCESSABLE,
			Errors: []exception.SpecificError{
				exception.AUTH_GENERATE_TOKEN_ERROR,
			},
		})
	}

	return newAccessTokenString
}
