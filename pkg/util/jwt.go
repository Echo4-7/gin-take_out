package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("shimily")

type Claims struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	//Authority int    `json:"authority"`
	Status string `json:"status"`
	jwt.StandardClaims
}

// GenerateToken 签发token
func GenerateToken(userID uint, status string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		ID: userID,
		//UserName:  userName,
		//Authority: authority,
		Status: status,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Take_Out",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

type EmailClaims struct {
	Nickname      string `json:"nick_name"`
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"` //1 绑定邮箱 2 解绑邮箱 3
	jwt.StandardClaims
}

// GenerateEmailToken 签发email token
func GenerateEmailToken(userId, operationType uint, email, nickname, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(5 * time.Minute)
	claims := EmailClaims{
		Nickname:      nickname,
		UserID:        userId,
		OperationType: operationType,
		Email:         email,
		Password:      password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Take_Out",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseEmailToken 解析email token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
