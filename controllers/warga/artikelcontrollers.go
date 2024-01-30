package warga

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllArtikel(c *gin.Context) {
	var resultStruct []struct {
		IdArtikel string `gorm:"column:id_artikel"`
		Judul string `gorm:"column:judul"`
		IsiArtikel string `gorm:"column:isi"`
	}

	result := databases.DB.Raw("SELECT id_artikel, judul, isi FROM artikel").Scan(&resultStruct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(resultStruct) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada artikel tersedia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": resultStruct})
}