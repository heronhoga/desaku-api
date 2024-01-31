package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"desaku-api/databases"
	"time"
)

func TokenAuthMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
	// Get token from cookies
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

	var existingToken string
	err = databases.DB.Raw("SELECT session FROM admin WHERE session = ?", token).Scan(&existingToken).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if existingToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token found"})
		c.Abort()
		return
	}

	expire := time.Now().Add(24 * time.Hour)
	c.SetCookie("token", token, int(time.Until(expire).Seconds()), "/", "", false, true)

	c.Next()
}
}
