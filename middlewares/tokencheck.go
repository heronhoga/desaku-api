package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"desaku-api/databases"
)

func TokenAuthMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
    // Get token from cookies
    token, err := c.Cookie("token")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
        return
    }

    var existingToken string
	err = databases.DB.Raw("SELECT session FROM admin WHERE session = ?", token).Scan(&existingToken).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token found"})
		return
	}



    c.Next()
}
}
