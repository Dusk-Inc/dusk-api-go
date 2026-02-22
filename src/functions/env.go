package functions

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EnvModel struct {
	Host string
	Port int
}

func ParseEnv() (EnvModel, error) {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	portRaw := os.Getenv("PORT")
	if portRaw == "" {
		return EnvModel{}, fmt.Errorf("PORT is required")
	}
	port, parseError := strconv.Atoi(portRaw)
	if parseError != nil || port < 1 || port > 65535 {
		return EnvModel{}, fmt.Errorf("PORT must be an integer between 1 and 65535")
	}

	return EnvModel{Host: host, Port: port}, nil
}

func SendNotImplemented(context *gin.Context) {
	context.JSON(501, gin.H{
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Not implemented.",
		},
	})
}
