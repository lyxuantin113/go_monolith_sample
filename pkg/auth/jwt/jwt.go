package jwt

import (
	"errors"
	"go_monolith_sample/internal/domain/common"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type Claims struct {
	UserID   uint            `json:"user_id"`
	Username string          `json:"username"`
	Role     common.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair tạo cặp Access và Refresh Token
func GenerateTokenPair(userID uint, username string, role common.UserRole) (*TokenPair, error) {
	access, err := generateToken(userID, username, role, AccessTokenDuration, accessSecret)
	if err != nil {
		return nil, err
	}

	refresh, err := generateToken(userID, username, role, RefreshTokenDuration, refreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// generateToken là hàm helper dùng chung để tạo token
func generateToken(userID uint, username string, role common.UserRole, duration time.Duration, secretKey []byte) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateAccessToken dùng cho Middleware để xác thực mỗi request
func ValidateAccessToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, accessSecret)
}

// ValidateRefreshToken dùng cho API /refresh-token để cấp lại access token mới
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, refreshSecret)
}

// parseToken là hàm helper dùng chung để giải mã với secret tương ứng
func parseToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký (Bảo mật)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("thuật toán ký không hợp lệ")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token không hợp lệ")
}
