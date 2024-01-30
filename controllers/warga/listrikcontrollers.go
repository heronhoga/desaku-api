package warga

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTagihanListrikWarga(c *gin.Context) {
	//ID WARGA
	id := c.Param("id")

	//VARIABLE FOR TAGIHAN
	var tagihanListrik []struct {
		IdListrik string `gorm:"column:id_listrik"`
		IdWarga string `gorm:"column:id_warga"`
		TotalTagihanListrik int64 `gorm:"column:total_tagihan_listrik"`
		TanggalTagihan string `gorm:"column:tanggal_tagihan"`
	}

	//GETTING ALL DATA
	result := databases.DB.Raw("SELECT id_listrik, id_warga, total_tagihan_listrik, tanggal_tagihan FROM listrik WHERE id_warga = ? AND status = 'pending'", id).Scan(&tagihanListrik)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(tagihanListrik) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada tagihan listrik tersedia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": tagihanListrik})
}