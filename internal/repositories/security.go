package repositories

import (
	"fmt"
	"regexp"
	"support/internal/core/domain"
	"time"

	"github.com/ignaciocaff/oraclesp"
	"github.com/mitchellh/mapstructure"
)

type repository struct {
}

func New() *repository {
	return &repository{}
}

func (s *repository) ObtenerFechaDefuncion(cuil string) (*time.Time, error) {
	var res Persona
	err := oraclesp.Execute("PKG_TRAMITES_CONSULTAS.PR_OBT_DATOS_FALLECIMIENTO", &res, cuil)
	if err != nil {
		return nil, err
	}
	return &res.FechaDefuncion, nil
}

func (s *repository) ObtenerDatosParaComparacion(cuil string) (*domain.DatosComparacion, error) {
	var entity DatosComparacion
	var res *domain.DatosComparacion
	err := oraclesp.Execute("PKG_SEGURIDAD_PERSONAS.PR_OBTENER_DATOS_COMPARAR", &entity, cuil)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(entity, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *repository) RegistrarDatosPersona(cidiData *domain.CidiDatos) (int, string, error) {
	var res StoreProcedureResultWithId
	idLocalidad := cidiData.Domicilio.IdLocalidad
	if idLocalidad == 0 {
		idLocalidad = -1
	}
	err := oraclesp.Execute("PKG_SEGURIDAD_PERSONAS.PR_REGISTRAR_DATOS_PERSONA", &res, cidiData.CUIL, cidiData.NroDocumento, cidiData.Id_Sexo,
		cidiData.Apellido, cidiData.Nombre, cidiData.FechaNacimiento, cidiData.TelFormateado, cidiData.CelFormateado, cidiData.Email, idLocalidad,
		cidiData.Domicilio.Barrio, cidiData.Domicilio.Depto, cidiData.Domicilio.Piso, cidiData.Domicilio.CodigoPostal, cidiData.Domicilio.Calle,
		cidiData.Domicilio.Altura, cidiData.Id_Estado)
	if err != nil {
		return 0, "ERROR", err
	}
	return res.Id, res.Mensaje, nil
}

func (s *repository) ActualizarDatosPersona(cidiData *domain.CidiDatos, perNro int) (string, error) {
	var res StoreProcedureResult
	idLocalidad := cidiData.Domicilio.IdLocalidad
	if idLocalidad == 0 {
		idLocalidad = -1
	}
	err := oraclesp.Execute("PKG_SEGURIDAD_PERSONAS.PR_ACTUALIZAR_DATOS_PERSONA", &res, perNro, cidiData.CelFormateado, cidiData.TelFormateado, cidiData.Email, cidiData.Domicilio.IdLocalidad,
		cidiData.Domicilio.Barrio, cidiData.Domicilio.Calle, cidiData.Domicilio.Altura, cidiData.Domicilio.Depto, cidiData.Domicilio.Piso,
		cidiData.Domicilio.CodigoPostal, cidiData.Id_Estado)
	if err != nil {
		return "ERROR", err
	}
	return res.Mensaje, nil
}

func (s *repository) ObtenerUsuario(cuil, idSexo, nroDocumento string, idNumero int) (domain.Usuario, error) {
	var entity Usuario
	var res domain.Usuario
	err := oraclesp.Execute("PKG_SEGURIDAD_CONSULTAS.PR_OBTENER_USUARIO_X_PERS", &entity, cuil, idSexo, nroDocumento, idNumero)
	if err != nil {
		return res, err
	}
	err = mapstructure.Decode(entity, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *repository) RegistrarRepresentado(userId int) (string, error) {
	var res StoreProcedureResult
	err := oraclesp.Execute("PKG_SEGURIDAD_ABM.PR_REGISTRAR_REPRESENTANTE_CIDI", &res, userId)
	if err != nil {
		return "ERROR", err
	}
	return res.Mensaje, nil
}

func (s *repository) ActualizarFuncionalidades(cuil *string) (string, error) {
	var res StoreProcedureResult
	err := oraclesp.Execute("PKG_SEGURIDAD_PERSONAS.PR_ACTUALIZAR_FUNC_X_CUIL", &res, cuil)
	if err != nil {
		return "ERROR", err
	}
	return res.Mensaje, nil
}

func (s *repository) RegistrarUsuarioWeb(perNum int, userName string) (domain.Usuario, error) {
	var entity Usuario
	var res domain.Usuario
	err := oraclesp.Execute("PKG_SEGURIDAD_PERSONAS.PR_REGISTRAR_USR_WEB", &entity, perNum, userName)
	if err != nil {
		return res, err
	}
	err = mapstructure.Decode(entity, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *repository) ObtenerTiposUsuarioPorUsuario(idUsuario, idTipoUsuario int) ([]domain.TipoUsuario, error) {
	var entity []TipoUsuario
	var res []domain.TipoUsuario

	err := oraclesp.Execute("PKG_SEGURIDAD_CONSULTAS.PR_OBT_TIPOS_USU_X_USR", &entity, idUsuario, idTipoUsuario)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(entity, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *repository) ObtenerPermisos(cuil string, idTipoUsuario int) ([]string, []string, error) {
	var entity []Funcionalidad
	err := oraclesp.Execute("PKG_SEGURIDAD_CONSULTAS.PR_OBT_FUNCIONALIDAD_X_CUIL", &entity, cuil, idTipoUsuario)
	if err != nil {
		return nil, nil, err
	}
	var urlsFront, urlRedis []string
	for _, elem := range entity {
		elem.Url = setUrl(elem.Url)
		if elem.Url != "" {
			if elem.IdMetodoUrl == 1 {
				urlsFront = append(urlsFront, elem.Url)
			} else {
				urlRedis = append(urlRedis, elem.Url)
			}
		}
	}
	return urlsFront, urlRedis, nil
}

func (s *repository) ObtenerMenu(cuil string, idTipoUsuario int) ([]domain.Menu, error) {
	var entity []Menu
	var res []domain.Menu
	err := oraclesp.Execute("PKG_SEGURIDAD_CONSULTAS.PR_OBT_MENU_X_CUIL", &entity, cuil, idTipoUsuario)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(entity, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func setUrl(url string) string {
	re := regexp.MustCompile("^/[^/]+/?([^/]+)/?")
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		url = match[0]
	}
	return url
}

type Persona struct {
	FechaDefuncion time.Time `oracle:"FechaFallecimiento"`
}

type DatosComparacion struct {
	PerNroInt    int    `oracle:"PerNroInt"`
	Celular      string `oracle:"Celular"`
	Telefono     string `oracle:"Telefono"`
	Correo       string `oracle:"Correo"`
	IdLocalidad  int    `oracle:"IdLocalidad"`
	Barrio       string `oracle:"Barrio"`
	Calle        string `oracle:"Calle"`
	Altura       string `oracle:"Altura"`
	Depto        string `oracle:"Depto"`
	Piso         string `oracle:"Piso"`
	CodigoPostal int    `oracle:"CodigoPostal"`
}

type TipoUsuario struct {
	Id           int    `oracle:"Id"`
	Nombre       string `oracle:"Nombre"`
	NombreImagen string `oracle:"NombreImagen"`
}
type StoreProcedureResultWithId struct {
	Id      int    `oracle:"Id"`
	Mensaje string `oracle:"Mensaje"`
}

type StoreProcedureResult struct {
	Mensaje string `oracle:"Mensaje"`
}

type Res5 struct {
	Id        int    `oracle:"Id"`
	Nombre    string `oracle:"Nombre"`
	ImageName string `oracle:"NombreImagen"`
}

type Usuario struct {
	IdUsuario     int    `oracle:"IdUsuario"`
	Apellido      string `oracle:"Apellido"`
	Nombre        string `oracle:"Nombre"`
	NroDocumento  string `oracle:"NroDocumento"`
	Cuil          string `oracle:"Cuil"`
	IdTipoUsuario int    `oracle:"IdTipoUsuario"`
}

type PersonaCidi struct {
	PerNroInt    int    `oracle:"PerNroInt"`
	Celular      string `oracle:"Celular"`
	Telefono     string `oracle:"Telefono"`
	Correo       string `oracle:"Correo"`
	IdLocalidad  int    `oracle:"IdLocalidad"`
	Barrio       string `oracle:"Barrio"`
	Calle        string `oracle:"Calle"`
	Altura       string `oracle:"Altura"`
	Depto        string `oracle:"Depto"`
	Piso         string `oracle:"Piso"`
	CodigoPostal int    `oracle:"CodigoPostal"`
}

type Funcionalidad struct {
	Url         string `oracle:"Url"`
	IdMetodoUrl int    `oracle:"IdMetodoUrl"`
}

type Menu struct {
	Id      int    `oracle:"Id"`
	Nombre  string `oracle:"Nombre"`
	Path    string `oracle:"Path"`
	IdPadre int    `oracle:"IdPadre"`
	Icon    string `oracle:"Icono"`
	SubMenu []Menu `oracle:"SubMenu,omitempty"`
}
