package jwtdata

import "github.com/golang-jwt/jwt/v5"

type JWTData struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Iss    string `json:"iss"`
	Exp    int64  `json:"exp"`
	jwt.RegisteredClaims
}
