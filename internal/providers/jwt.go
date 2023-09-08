package providers

import (
	"fmt"
	"support/internal/core/domain"
	"support/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	config env.EnvApp
}

type CustomClaims struct {
	Usuario      domain.Usuario       `json:"usuario"`
	TiposUsuario []domain.TipoUsuario `json:"tipoUsuario"`
	jwt.RegisteredClaims
}

func NewJwt(config env.EnvApp) *jwtService {
	return &jwtService{config: config}
}

func (j *jwtService) GenerarToken(payload *domain.TokenPayload) (string, error) {

	customClaims := CustomClaims{
		Usuario:      payload.Usuario,
		TiposUsuario: payload.TiposUsuario,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now())},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, customClaims)
	tokenFirmado, err := token.SignedString([]byte(j.config.JWT_SECRET))
	if err != nil {
		fmt.Println("Error al firmar el token:", err)
		return "", err
	}
	return tokenFirmado, nil
}

func (j *jwtService) ValidarToken(tokenString string) (*domain.TokenPayload, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.JWT_SECRET), nil
	})
	if err != nil {
		fmt.Println("Error al validar el token:", err)
		return nil, err
	}
	payload := domain.TokenPayload{}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		payload.TiposUsuario = claims.TiposUsuario
		payload.Usuario = claims.Usuario
	} else {
		fmt.Println("Error al validar el token:", err)
		return nil, err
	}
	return &payload, nil
}
