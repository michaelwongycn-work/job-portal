package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("")
var accessTokenDuration = time.Duration(0)
var refreshTokenDuration = time.Duration(0)

func SetAuthConfig(key string, newAccessTokenDuration, newRefreshTokenDuration time.Duration) {
	secretKey = []byte(key)
	accessTokenDuration = newAccessTokenDuration
	refreshTokenDuration = newRefreshTokenDuration
}

func CreateToken(currTime time.Time, ID, role int) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "job-portal"
	claims["aud"] = "job-portal:crypto"
	claims["sub"] = ID
	claims["rle"] = role
	claims["exp"] = currTime.Add(time.Minute * accessTokenDuration).Unix()
	claims["iat"] = currTime.Unix()

	accessTokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = ID
	rtClaims["exp"] = currTime.Add(time.Minute * refreshTokenDuration).Unix()

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
