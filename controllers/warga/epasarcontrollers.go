package warga

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllEpasar(c *gin.Context) {
	var resultStruct []struct {
		IdToko string `gorm:"column:id_toko"`
		NamaToko string `gorm:"column:nama_toko"`
		NamaPedagang string `gorm:"column:nama_pedagang"`
		JenisDagangan string `gorm:"column:jenis_dagangan"`
		Status string `gorm:"column:status"`
	}

	result := databases.DB.Raw("SELECT id_toko, nama_toko, nama_pedagang, jenis_dagangan, status FROM epasar").Scan(&resultStruct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(resultStruct) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada toko tersedia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": resultStruct})
}