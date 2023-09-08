package handlers

import (
	"fmt"
	"support/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	securityService ports.SecurityService
	cidiService     ports.CidiService
}

func New(securityService ports.SecurityService, cidiService ports.CidiService) *HTTPHandler {
	return &HTTPHandler{securityService: securityService, cidiService: cidiService}
}

func (h *HTTPHandler) UsuarioTemporal(c *gin.Context) {
	var json struct {
		Message string `json:"message"`
	}
	signedToken, err := h.securityService.UsuarioTemporal()
	if err != nil {
		fmt.Println("Error al firmar el token:", err)
		return
	}
	if err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}

	c.Header("authorization", signedToken)
	c.JSON(200, gin.H{
		"message": json.Message,
	})

}

func (h *HTTPHandler) Login(c *gin.Context) {

	cookie, err := c.Cookie("CiDi")
	if err != nil {
		// Cookie not found
		c.String(409, "Cookie not found")
		return
	}
	cidiData, err := h.cidiService.Login(cookie)
	if err != nil {
		c.String(403, "Error al obtener los datos del usuario de cidi")
		return
	}
	token, estaFallecido, err := h.securityService.Login(cidiData)

	if err != nil {
		c.String(500, "Error al obtener los datos del usuario")
		return
	}
	if estaFallecido {
		c.JSON(200, gin.H{"deceasedDate": true})
		return
	}
	c.Header("authorization", token)
	c.JSON(200, gin.H{})
}

func (h *HTTPHandler) Permisos(c *gin.Context) {
	tokenString := c.GetHeader("authorization")
	var json PermisosJson

	if err := c.ShouldBindJSON(&json); err != nil {
		c.String(409, "Error al obtener los datos del usuario")
		return
	}
	permisosEncriptados := h.securityService.Permisos(json.IdTipoUsuario, json.IdUsuario, tokenString)
	c.JSON(200, permisosEncriptados)
}

func (h *HTTPHandler) Menu(c *gin.Context) {
	tokenString := c.GetHeader("authorization")
	var json MenuJson

	if err := c.ShouldBindJSON(&json); err != nil {
		c.String(409, "Error al obtener los datos del usuario")
		return
	}
	menu := h.securityService.Menu(json.IdTipoUsuario, tokenString)
	c.JSON(200, menu)
}

func (h *HTTPHandler) ObtenerRepresentado(c *gin.Context) {

	cookie, err := c.Cookie("CiDi")
	if err != nil {
		// Cookie not found
		c.String(409, "Cookie not found")
		return
	}
	cidiData, err := h.cidiService.Login(cookie)
	if err != nil {
		c.String(409, "Error al obtener los datos del usuario de cidi")
		return
	}
	representado, err := h.securityService.ObtenerRepresentado(cidiData)
	if err != nil {
		c.String(409, "Error al obtener los datos del usuario")
		return
	}
	c.JSON(200, representado)
}

type PermisosJson struct {
	IdUsuario     int `json:"idUsuario"`
	IdTipoUsuario int `json:"idTipoUsuario"`
}

type MenuJson struct {
	IdTipoUsuario int `json:"idTipoUsuario"`
}
