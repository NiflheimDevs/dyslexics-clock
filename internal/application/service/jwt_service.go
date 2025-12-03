package service

type JWT interface {
	GenerateToken(deviceID uint) (string, string)
	VerifyToken(tokenString string) (map[string]any, error)
	RefreshToken(token string) string
}
