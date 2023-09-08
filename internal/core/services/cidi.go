package services

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"support/internal/core/domain"
	"support/internal/env"
	"time"
)

type cidiService struct {
	ctx    context.Context
	config env.EnvApp
}

func NewCidi(ctx context.Context, config env.EnvApp) *cidiService {
	return &cidiService{ctx, config}
}

type EntradaCidi struct {
	IdAplicacion int    `json:"IdAplicacion"`
	Contrasenia  string `json:"Contrasenia"`
	HashCookie   string `json:"HashCookie"`
	TimeStamp    string `json:"TimeStamp"`
	TokenValue   string `json:"TokenValue"`
}

const (
	UserGetApplicationUser = "/Usuario/Obtener_Usuario_Aplicacion"
)

func (s *cidiService) Login(cookie string) (*domain.CidiDatos, error) {
	cidiApp, err := strconv.Atoi(s.config.ID_APP)
	if err != nil {
		fmt.Println("Error al convertir el ID_APP:", err)
	}
	token := genericRequest(cookie, s.config.CIDI_KEY, s.config.CIDI_PASS, cidiApp)
	body, err := json.Marshal(token)
	if err != nil {
		fmt.Printf("Error al armar el body: %v", err)
		return nil, err
	}

	client := http.Client{Timeout: 40 * time.Second}
	req, err := http.NewRequest("POST", s.config.BASE_CIDI_URI+UserGetApplicationUser, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error en generacion de peticion CIDI: %v \n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error en ejecucion de peticion CIDI: %v \n", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Decodificar la respuesta JSON y devolver un objeto CidiDatos
	var cidiData domain.CidiDatos
	err = json.NewDecoder(resp.Body).Decode(&cidiData)
	if err != nil {
		fmt.Printf("Error al decodificar la respuesta JSON: %v \n", err)
		return nil, err
	}

	return &cidiData, nil
}

func genericRequest(cookie, cidiKey, cidiPass string, idApp int) *EntradaCidi {
	timeStamp := getTimeStamp()
	token := generateTokenSHA(timeStamp, cidiKey)
	entrada := EntradaCidi{
		IdAplicacion: idApp,
		Contrasenia:  cidiPass,
		HashCookie:   cookie,
		TimeStamp:    timeStamp,
		TokenValue:   token,
	}
	return &entrada
}

func getTimeStamp() string {
	return time.Now().Format("20060102150405.000")
}

func generateTokenSHA(timeStamp, cidiKey string) string {
	strToHash := timeStamp + cidiKey
	hash := sha1.Sum([]byte(strToHash))
	return fmt.Sprintf("%X", hash)
}
