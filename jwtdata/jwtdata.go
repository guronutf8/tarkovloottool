package jwtdata

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type JWTData struct {
	UserId  string `json:"userId"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func JWTUserFromCtx(ctx context.Context) (JWTData, error) {
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return JWTData{}, errors.New("metadata FromIncomingContext")
	}
	payload := incomingContext.Get("x-jwt-payload")

	return JWTUserFromPayload(payload)

}

func JWTUserFromPayload(payload []string) (JWTData, error) {
	if len(payload) != 1 || payload[0] == "" {
		return JWTData{}, errors.New("meta no have payload")
	}

	sDec, _ := b64.StdEncoding.DecodeString(payload[0])
	data := &JWTData{}
	bytes := []byte("}")
	sDec = append(sDec, bytes[0])
	err := json.Unmarshal(sDec, &data)
	if err != nil {
		return JWTData{}, errors.Join(errors.New("meta DecodeString"), err)
	}

	return *data, nil
}

/*func JWTUserFromCtx(ctx context.Context) (JWTData, error) {
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return JWTData{}, errors.New("metadata FromIncomingContext")
	}
	payload := incomingContext.Get("x-jwtdata-payload")
	if len(payload) != 1 || payload[0] == "" {
		return JWTData{}, errors.New("meta no have payload")
	}

	parser := jwtdata.NewParser(jwtdata.WithoutClaimsValidation())
	token, err := parser.ParseWithClaims(payload[0], &JWTData{}, func(t *jwtdata.Token) (interface{}, error) { return nil, nil })
	if err != nil {
		return JWTData{}, errors.Join(errors.New("ParseWithClaims"), err)
	}

	return token.Claims.(JWTData), nil
}
*/
