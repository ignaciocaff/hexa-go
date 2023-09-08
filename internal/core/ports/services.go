package ports

import (
	"support/internal/core/domain"
)

type SecurityService interface {
	UsuarioTemporal() (string, error)
	Login(cidiData *domain.CidiDatos) (string, bool, error)
	Permisos(idTipoUsuario, idUsuario int, tokenString string) []string
	Menu(idTipoUsuario int, tokenString string) []*domain.Menu
	ObtenerRepresentado(cidiData *domain.CidiDatos) (domain.Representado, error)
}

type CidiService interface {
	Login(cookie string) (*domain.CidiDatos, error)
}

type JwtService interface {
	GenerarToken(payload *domain.TokenPayload) (string, error)
	ValidarToken(tokenString string) (*domain.TokenPayload, error)
}
