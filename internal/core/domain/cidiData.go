package domain

import (
	"strconv"
)

type CidiDatos struct {
	CUIL               string       `json:"CUIL"`
	CuilFormateado     string       `json:"CuilFormateado"`
	NroDocumento       string       `json:"NroDocumento"`
	Apellido           string       `json:"Apellido"`
	Nombre             string       `json:"Nombre"`
	NombreFormateado   string       `json:"NombreFormateado"`
	FechaNacimiento    string       `json:"FechaNacimiento"`
	Id_Sexo            string       `json:"Id_Sexo"`
	PaiCodPais         string       `json:"PaiCodPais"`
	Id_Numero          int          `json:"Id_Numero"`
	Id_Estado          int          `json:"Id_Estado"`
	Nivel1Seguro       string       `json:"Nivel1Seguro"`
	Estado             string       `json:"Estado"`
	Email              string       `json:"Email"`
	TelArea            string       `json:"TelArea"`
	TelNro             string       `json:"TelNro"`
	TelFormateado      string       `json:"TelFormateado"`
	CelArea            string       `json:"CelArea"`
	CelNro             string       `json:"CelNro"`
	CelFormateado      string       `json:"CelFormateado"`
	Empleado           string       `json:"Empleado"`
	Id_Empleado        string       `json:"Id_Empleado"`
	FechaRegistro      string       `json:"FechaRegistro"`
	FechaBloqueo       string       `json:"FechaBloqueo"`
	IdAplicacionOrigen string       `json:"IdAplicacionOrigen"`
	Domicilio          Domicilio    `json:"Domicilio"`
	CodigoIngresado    string       `json:"CodigoIngresado"`
	Constatado         string       `json:"Constatado"`
	Representado       Representado `json:"Representado"`
	Respuesta          Respuesta    `json:"Respuesta"`
	TieneRepresentados string       `json:"TieneRepresentados,omitempty"`
}

type Domicilio struct {
	IdVin          int    `json:"IdVin"`
	IdPais         int    `json:"IdPais,omitempty"`
	Pais           string `json:"Pais,omitempty"`
	IdProvincia    string `json:"IdProvincia"`
	Provincia      string `json:"Provincia"`
	IdDepartamento int    `json:"IdDepartamento"`
	Departamento   string `json:"Departamento"`
	IdLocalidad    int    `json:"IdLocalidad"`
	Localidad      string `json:"Localidad"`
	IdBarrio       int    `json:"IdBarrio"`
	Barrio         string `json:"Barrio"`
	IdCalle        int    `json:"IdCalle"`
	Calle          string `json:"Calle"`
	Altura         string `json:"Altura"`
	CodigoPostal   string `json:"CodigoPostal"`
	Piso           string `json:"Piso"`
	Depto          string `json:"Depto"`
	Torre          string `json:"Torre"`
	Manzana        string `json:"Manzana,omitempty"`
	Lote           string `json:"Lote"`
}

type Representado struct {
	RdoCuilCuit       string `json:"rdoCuilCuit,omitempty"`
	RdoNombre         string `json:"rdoNombre,omitempty"`
	RdoDenominacion   string `json:"rdoDenominacion,omitempty"`
	RdoTipo           string `json:"rdoTipo,omitempty"`
	RdoIdEstado       int    `json:"rdoIdEstado,omitempty"`
	HashRepresentado  string `json:"hashRepresentado,omitempty"`
	TieneRepresentado string `json:"tieneRepresentado,omitempty"`
}

type Respuesta struct {
	Resultado     string  `json:"Resultado,omitempty"`
	CodigoError   string  `json:"CodigoError,omitempty"`
	SesionHash    string  `json:"SesionHash,omitempty"`
	ExisteUsuario string  `json:"ExisteUsuario,omitempty"`
	Servidor      string  `json:"Servidor,omitempty"`
	TimeStamp     *string `json:"TimeStamp,omitempty" time:"20060102150405.000"`
}

type DatosComparacion struct {
	PerNroInt    int    `json:"PerNroInt,omitempty"`
	Celular      string `json:"Celular,omitempty"`
	Telefono     string `json:"Telefono,omitempty"`
	Correo       string `json:"Correo,omitempty"`
	IdLocalidad  int    `json:"IdLocalidad,omitempty"`
	Barrio       string `json:"Barrio,omitempty"`
	Calle        string `json:"Calle,omitempty"`
	Altura       string `json:"Altura,omitempty"`
	Depto        string `json:"Depto,omitempty"`
	Piso         string `json:"Piso,omitempty"`
	CodigoPostal int    `json:"CodigoPostal,omitempty"`
}

type Comparer interface {
	New(CidiDatos, int) (*DatosComparacion, error)
}

func (d DatosComparacion) New(cidiData *CidiDatos, perNro int) (*DatosComparacion, error) {
	num := 5000
	if cidiData.Domicilio.CodigoPostal != "" {
		num, _ = strconv.Atoi(cidiData.Domicilio.CodigoPostal)

	}
	return &DatosComparacion{PerNroInt: perNro,
		Celular: cidiData.CelFormateado, Telefono: cidiData.TelFormateado, Correo: cidiData.Email,
		IdLocalidad: cidiData.Domicilio.IdLocalidad, Barrio: cidiData.Domicilio.Barrio,
		Calle: cidiData.Domicilio.Calle, Altura: cidiData.Domicilio.Altura, Depto: cidiData.Domicilio.Depto,
		Piso: cidiData.Domicilio.Piso, CodigoPostal: num}, nil
}
