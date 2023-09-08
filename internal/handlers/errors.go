package handlers

import (
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
)


func RestartOnErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Verificar si ocurrió un error
		if c.Writer.Status() >= 400 {
			errMsg := c.Errors.String()

			// Verificar si el mensaje de error contiene "ORA-02343"
			if strings.Contains(errMsg, "ORA-02391") {
				fmt.Println("Reiniciando el microservicio en Docker...")
			}
		}
	}
}
