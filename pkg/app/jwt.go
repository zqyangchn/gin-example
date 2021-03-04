package app

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gin-example/pkg/setting"
)

type Claims struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(setting.JWTSetting.Secret)
}

func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(setting.JWTSetting.Expire)

	claims := Claims{
		AppKey:    appKey,
		AppSecret: appSecret,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    setting.JWTSetting.Issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(tokenClaims)
	return tokenClaims.SignedString(GetJWTSecret())
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return GetJWTSecret(), nil
		},
	)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func ParseTokenWithoutValid(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return GetJWTSecret(), nil
		},
	)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok {
			return claims, nil
		}
	}
	return nil, err
}
