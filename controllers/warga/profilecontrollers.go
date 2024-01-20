package warga

import (
	"github.com/gin-gonic/gin"
)

func LoginWarga(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Login Warga",
	})
}

func RegisterWarga(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Register Warga",
	})
}
