package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"support/internal/core/domain"
	"support/internal/core/ports"
)

type service struct {
	securityRepository ports.SecurityRepository
	jwtService         ports.JwtService
	redisRepository    ports.RedisRepository
}

func New(securityRepository ports.SecurityRepository, jwtService ports.JwtService, redisRepository ports.RedisRepository) *service {
	return &service{securityRepository: securityRepository, jwtService: jwtService, redisRepository: redisRepository}
}

func (s *service) Menu(idTipoUsuario int, tokenString string) []*domain.Menu {
	token, err := s.jwtService.ValidarToken(tokenString)
	if err != nil {
		fmt.Printf("Error al validar el token: %v\n", err)
		return nil
	}
	cuil := token.Usuario.Cuil
	if token.Usuario.CuilRepresentado != "" {
		cuil = token.Usuario.CuilRepresentado
		idTipoUsuario = 5
	}
	menu, err := s.securityRepository.ObtenerMenu(cuil, idTipoUsuario)

	resMenu := buildMenu(menu)
	if err != nil {
		fmt.Printf("Error al obtener el menu: %v\n", err)
		return nil
	}
	return resMenu
}

func (s *service) Permisos(idTipoUsuario, idUsuario int, tokenString string) []string {
	token, err := s.jwtService.ValidarToken(tokenString)
	if err != nil {
		fmt.Printf("Error al validar el token: %v\n", err)
		return nil
	}
	cuil := token.Usuario.Cuil
	if token.Usuario.CuilRepresentado != "" {
		cuil = token.Usuario.CuilRepresentado
		idTipoUsuario = 5
	}
	urlsFront, urlsRedis, err := s.securityRepository.ObtenerPermisos(cuil, idTipoUsuario)
	if err != nil {
		fmt.Printf("Error al obtener los permisos: %v\n", err)
		return nil
	}
	urlsJson, _ := json.Marshal(urlsRedis)

	err = s.redisRepository.Set(strconv.Itoa(idUsuario), urlsJson)
	if err != nil {
		fmt.Printf("Error al setear en redis: %v\n", err)
		return nil
	}
	return urlsFront
}

func (s *service) UsuarioTemporal() (string, error) {
	payload := domain.TokenPayload{}
	payload.CrearUsuarioTemporal()
	token, err := s.jwtService.GenerarToken(&payload)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) ObtenerRepresentado(cidiData *domain.CidiDatos) (domain.Representado, error) {
	var representado domain.Representado
	representado.HashRepresentado = cidiData.Respuesta.SesionHash
	representado.TieneRepresentado = cidiData.TieneRepresentados
	tieneRepresentado := tieneRepresentado(cidiData)

	if tieneRepresentado {
		representado.RdoCuilCuit = cidiData.Representado.RdoCuilCuit
		representado.RdoNombre = cidiData.Representado.RdoNombre
		representado.RdoDenominacion = cidiData.Representado.RdoDenominacion
		representado.RdoTipo = cidiData.Representado.RdoTipo
		representado.RdoIdEstado = cidiData.Representado.RdoIdEstado
	}

	return representado, nil
}

func (s *service) Login(cidiData *domain.CidiDatos) (string, bool, error) {
	var perNumber int
	var signedToken string
	fechaDefuncion, err := s.securityRepository.ObtenerFechaDefuncion(cidiData.CUIL)
	if err != nil {
		return "", false, err
	}

	if fechaDefuncion.IsZero() {
		// Person is alive
		datosParaComparar, err := s.securityRepository.ObtenerDatosParaComparacion(cidiData.CUIL)
		if err != nil {
			return "", false, err
		}
		if datosParaComparar.PerNroInt != 0 {
			perNumber = datosParaComparar.PerNroInt
			// Person exists, should compare and update if necessary
			cidiDataParaComparar, err := datosParaComparar.New(cidiData, datosParaComparar.PerNroInt)
			if err != nil {
				return "", false, err
			}

			sonIguales := structFieldsEqual(cidiDataParaComparar, datosParaComparar)
			if !sonIguales {
				s.securityRepository.ActualizarDatosPersona(cidiData, datosParaComparar.PerNroInt)
			}
		} else {
			// Person does not exists, should register person data
			id, mensaje, err := s.securityRepository.RegistrarDatosPersona(cidiData)
			if err != nil {
				return "", false, err
			}
			if mensaje == "OK" {
				perNumber = id
			}
		}
		// Should obtain usuario for person
		usuario, err := s.securityRepository.ObtenerUsuario(cidiData.CUIL, cidiData.Id_Sexo, cidiData.NroDocumento, cidiData.Id_Numero)
		if err != nil {
			return "", false, err
		}
		if usuario.IdUsuario != 0 {
			// usuario exists Check if has represented
			tieneRepresentado := tieneRepresentado(cidiData)

			if tieneRepresentado {
				// Has represented should create it
				_, err := s.securityRepository.RegistrarRepresentado(usuario.IdUsuario)
				if err != nil {
					return "", false, err
				}
			}
			// Should update functionalities
			s.securityRepository.ActualizarFuncionalidades(obtenerCuilRepresentadoYFuncionalidad(cidiData, &cidiData.CUIL))

			// Should check if web usuario exists
			if usuario.IdTipoUsuario == 0 {
				usuario, err = s.registrarUsuarioWeb(perNumber, "")
				if err != nil {
					return "", false, err
				}
			}
		} else {
			// usuario does not exists should create
			usuario, err = s.registrarUsuarioWeb(perNumber, "")
			if err != nil {
				return "", false, err

			}
		}

		// Should obtain usuario types
		tiposUsuario, err := s.obtenerTiposUsuario(usuario, cidiData)
		if err != nil {
			return "", false, err
		}
		obtenerCuilRepresentadoYFuncionalidad(cidiData, &usuario.CuilRepresentado)
		usuario.NivelCidi = cidiData.Id_Estado
		signedToken, err = s.generarToken(&domain.TokenPayload{Usuario: usuario, TiposUsuario: tiposUsuario})
		if err != nil {
			return "", false, err
		}
		return signedToken, false, nil
	} else {
		// Person is death
		return "", true, nil
	}
}

func structFieldsEqual(a, b interface{}) bool {
	valA := reflect.ValueOf(a).Elem()
	valB := reflect.ValueOf(b).Elem()

	if valA.Kind() != reflect.Struct || valB.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < valA.NumField(); i++ {
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)
		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			return false
		}
	}

	return true
}

func tieneRepresentado(cidiData *domain.CidiDatos) bool {
	return cidiData.TieneRepresentados == "S" && cidiData.Representado.RdoCuilCuit != "" &&
		cidiData.Representado.RdoCuilCuit != cidiData.CUIL
}

func obtenerCuilRepresentadoYFuncionalidad(cidiData *domain.CidiDatos, cuil *string) *string {
	tieneRepresentado := tieneRepresentado(cidiData)
	if tieneRepresentado {
		cuil = &cidiData.Representado.RdoCuilCuit
	}
	return cuil
}

func (s *service) registrarUsuarioWeb(perNum int, nombreUsuario string) (domain.Usuario, error) {
	usuarioWeb, err := s.securityRepository.RegistrarUsuarioWeb(perNum, nombreUsuario)
	if err != nil {
		return usuarioWeb, err
	}
	return usuarioWeb, nil

}

func (s *service) obtenerTiposUsuario(usuario domain.Usuario, cidiData *domain.CidiDatos) ([]domain.TipoUsuario, error) {
	tieneRepresentado := tieneRepresentado(cidiData)
	idTipoUsuario := 0
	if tieneRepresentado {
		idTipoUsuario = 5
	}
	tiposUsuario, err := s.securityRepository.ObtenerTiposUsuarioPorUsuario(usuario.IdUsuario, idTipoUsuario)
	if err != nil {
		return nil, err
	}
	return tiposUsuario, nil
}

func (s *service) generarToken(payload *domain.TokenPayload) (string, error) {
	tokenFirmado, err := s.jwtService.GenerarToken(payload)
	if err != nil {
		fmt.Println("Error al firmar el token:", err)
		return "", err
	}
	return tokenFirmado, nil
}

func buildMenu(menu []domain.Menu) []*domain.Menu {
	menuItems := make(map[int]*domain.Menu)

	for i, item := range menu {
		menuItems[item.Id] = &menu[i]
	}

	var output []*domain.Menu

	for _, item := range menu {
		if item.IdPadre == 0 {
			output = append(output, menuItems[item.Id])
		} else {
			parent := menuItems[item.IdPadre]
			parent.SubMenu = append(parent.SubMenu, menuItems[item.Id])
		}
	}
	return output
}
