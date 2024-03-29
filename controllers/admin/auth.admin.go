package admin

import (
	"desaku-api/databases"
	"desaku-api/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginAdmin(c *gin.Context) {
	var loginData struct {
	    Username string `json:"username"`
	    Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	    return
	}

	if loginData.Username == "" || loginData.Password == "" {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Both 'username' and 'password' fields must be filled"})
	    return
	}

	var admin struct {
	    Username string `json:"username"`
	    Password string `json:"password"`
	}

	result := databases.DB.Raw("SELECT * FROM admin WHERE username = ? AND password = ?", loginData.Username, loginData.Password).Scan(&admin)

	if result.Error != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	    return
	}

	if admin.Username == "" || admin.Password == "" {
	    c.JSON(http.StatusBadRequest, gin.H{
	        "error": "Invalid username or password",
	        "loginStatus": false,
	    })
	    return
	}

	var adminID int64
	row := databases.DB.Raw("SELECT id FROM admin WHERE username = ? AND password = ?", loginData.Username, loginData.Password).Row()
	err := row.Scan(&adminID)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	    return
	}

	token, err := middlewares.GenerateSessionToken()
	
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	    return
	}

	tokenSet := databases.DB.Exec("UPDATE admin SET session = ? WHERE id = ?", token, adminID)
	if tokenSet.Error != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": tokenSet.Error.Error()})
	    return
	}

	expire := time.Now().Add(24 * time.Hour)
	c.SetCookie("token", token, int(time.Until(expire).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
	    "statusCode": http.StatusOK,
	    "loginStatus": true,
	    "id_admin": adminID,
	    "data": admin,
	})
	}

func CheckAdmin(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"loginStatus": true,
		})
	}

func LogoutAdmin(c *gin.Context) {
	token, err := c.Cookie("token")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
        return
    }

	err = databases.DB.Exec("UPDATE admin SET session = NULL WHERE session = ?", token).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"logoutStatus": true,
	})
}
