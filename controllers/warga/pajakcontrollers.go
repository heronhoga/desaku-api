package warga

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTagihanPajakWarga(c *gin.Context) {
    id := c.Param("id")

    var tagihanPajak []struct {
        IdPajak string `gorm:"column:id_pajak"`
		IdWarga string `gorm:"column:id_warga"`
		Tahun string `gorm:"column:tahun"`
		TotalTagihanPajak string `gorm:"column:total_tagihan_pajak"`

    }

    result := databases.DB.Raw("SELECT id_pajak, id_warga, tahun, total_tagihan_pajak FROM pajak WHERE id_warga = ? AND status_bayar = 'pending'", id).Scan(&tagihanPajak)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if len(tagihanPajak) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada tagihan pajak tersedia"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": tagihanPajak})
}

func BayarPajakWarga(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}