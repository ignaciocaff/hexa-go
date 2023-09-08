package domain

type Menu struct {
	Id      int     `json:"id"`
	Nombre  string  `json:"nombre"`
	Path    string  `json:"path"`
	IdPadre int     `json:"idPadre"`
	Icon    string  `json:"icon"`
	SubMenu []*Menu `json:"subMenu,omitempty"`
}

type Funcionalidad struct {
	Url         string `json:"url"`
	IdMetodoUrl int    `json:"idMetodoUrl"`
}
