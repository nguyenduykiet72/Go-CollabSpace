package server

import "github.com/gin-gonic/gin"

func resolveGinMode(mode string) string {
	switch mode {
	case "development", "dev":
		return gin.DebugMode
	case "production", "prod":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
