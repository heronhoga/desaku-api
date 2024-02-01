package admin

import (
	"net/http"
	"desaku-api/databases"

	"github.com/gin-gonic/gin"
)

func GetAllToko (c *gin.Context) {

	var epasar[] struct {
		IdToko string `json:"id_toko"`
		NamaToko string `json:"nama_toko"`
		NamaPedagang string `json:"nama_pedagang"`
		JenisDagangan string `json:"jenis_dagangan"`
		Status string `json:"status"`
	}

	tokoQuery := databases.DB.Raw(`SELECT id_toko, nama_toko, nama_pedagang, jenis_dagangan, status FROM epasar`).Scan(&epasar)

	if tokoQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tokoQuery.Error.Error()})
		return
	}

	if tokoQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": epasar})
}