package ports

import (
	"support/internal/core/domain"
	"time"
)

type SecurityRepository interface {
	ObtenerFechaDefuncion(cuil string) (*time.Time, error)
	ObtenerDatosParaComparacion(cuil string) (*domain.DatosComparacion, error)
	RegistrarDatosPersona(cidiData *domain.CidiDatos) (int, string, error)
	ActualizarDatosPersona(cidiData *domain.CidiDatos, perNro int) (string, error)
	ObtenerUsuario(cuil, idSexo, nroDocumento string, idNumero int) (domain.Usuario, error)
	RegistrarRepresentado(userId int) (string, error)
	ActualizarFuncionalidades(cuil *string) (string, error)
	RegistrarUsuarioWeb(perNum int, userName string) (domain.Usuario, error)
	ObtenerTiposUsuarioPorUsuario(idUsuario, idTipoUsuario int) ([]domain.TipoUsuario, error)
	ObtenerPermisos(cuil string, idTipoUsuario int) ([]string, []string, error)
	ObtenerMenu(cuil string, idTipoUsuario int) ([]domain.Menu, error)
}

type RedisRepository interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
}
