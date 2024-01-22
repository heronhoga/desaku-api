package warga

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BayarPajakWarga(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}