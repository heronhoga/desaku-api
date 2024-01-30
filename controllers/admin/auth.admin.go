package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"desaku-api/databases"
	"desaku-api/middlewares"
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

	c.JSON(http.StatusOK, gin.H{
	    "statusCode": http.StatusOK,
	    "loginStatus": true,
	    "id_admin": adminID,
	    "data": admin,
	})
	}
