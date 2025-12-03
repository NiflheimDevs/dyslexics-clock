package serviceimpl

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Constants  *bootstrap.Constants
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJWT(Const *bootstrap.Constants) *JWT {
	return &JWT{
		Constants:  Const,
		PrivateKey: loadPrivateKey(Const.JWT.JWTKeysPath + "/privateKey.pem"),
		PublicKey:  loadPublicKey(Const.JWT.JWTKeysPath + "/publicKey.pem"),
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

func (j *JWT) GenerateToken(deviceID uint) (string, string) {
	accessTokenClaims := jwt.MapClaims{
		"iss": j.Constants.JWT.Issuer,
		"sub": deviceID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(j.PrivateKey)
	if err != nil {
		panic(derror.New(derror.ErrTypeInternal, "error generating access token", err))
	}

	refreshTokenClaims := jwt.MapClaims{
		"iss": j.Constants.JWT.Issuer,
		"sub": deviceID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(j.PrivateKey)
	if err != nil {
		panic(derror.New(derror.ErrTypeInternal, "error generating refresh token", err))
	}

	return accessTokenString, refreshTokenString
}

func (j *JWT) VerifyToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			panic(derror.New(derror.ErrTypeUnauthorized, "unexpected signing method:"+token.Header["alg"].(string), nil))
		}
		return j.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, derror.New(derror.ErrTypeUnauthorized,"failed to parse claims", nil)
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
		panic (derror.New(derror.ErrTypeUnauthorized, "error verifying token", err))
	}

	deviceID, ok := claims["sub"].(uint)
	if !ok {
		panic (derror.New(derror.ErrTypeUnauthorized, "error verifying token", err))
	}

	newAccessTokenClaims := jwt.MapClaims{
		"iss": j.Constants.JWT.Issuer,
		"sub": uint(deviceID),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, newAccessTokenClaims)
	newAccessTokenString, err := newAccessToken.SignedString(j.PrivateKey)
	if err != nil {
		panic(derror.New(derror.ErrTypeInternal, "error generating access token", err))
	}
	return newAccessTokenString
}
