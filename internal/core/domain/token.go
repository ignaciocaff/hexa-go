package domain

type Persona struct {
	PerNroInt    int    `json:"perNroInt"`
	Celular      string `json:"celular"`
	Telefono     string `json:"telefono"`
	Correo       string `json:"correo"`
	IdLocalidad  int    `json:"idLocalidad"`
	Barrio       string `json:"barrio"`
	Calle        string `json:"calle"`
	Altura       string `json:"altura"`
	Depto        string `json:"depto"`
	Piso         string `json:"piso"`
	CodigoPostal int    `json:"codigoPostal"`
}

type Usuario struct {
	IdUsuario        int    `json:"id"`
	IdTipoUsuario    int    `json:"idTipoUsuario"`
	Apellido         string `json:"apellido"`
	Nombre           string `json:"nombre"`
	NroDocumento     string `json:"nroDocumento"`
	Cuil             string `json:"cuil"`
	CuilRepresentado string `json:"cuilRepresentado"`
	NivelCidi        int    `json:"nivelCidi"`
}

type TipoUsuario struct {
	Id           int    `json:"id"`
	Nombre       string `json:"nombre"`
	NombreImagen string `json:"nombreImagen"`
}

type TokenPayload struct {
	TiposUsuario []TipoUsuario `json:"tipoUsuario"`
	Usuario      Usuario       `json:"usuario"`
}

func (t *TokenPayload) CrearUsuarioTemporal() {
	t.Usuario.Nombre = "temporal"
	t.Usuario.Apellido = "temporal"
	t.Usuario.NroDocumento = "temporal"
	t.Usuario.Cuil = "temporal"
	t.Usuario.CuilRepresentado = "temporal"
	t.Usuario.NivelCidi = 2
	t.TiposUsuario = []TipoUsuario{}
}
