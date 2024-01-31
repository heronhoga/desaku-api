package admin

import (
	"net/http"
	"desaku-api/databases"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
)

func GetAllTagihanListrik(c *gin.Context) {

	var TagihanListrik[] struct {
		IdListrik string `json:"id_listrik"`
		TotalTagihanListrik string `json:"total_tagihan_listrik"`
		TanggalTagihan string `json:"tanggal_tagihan"`
		Status string `json:"status"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

currentYear := time.Now().Year()

monthParam := c.Param("month")

formattedDate := fmt.Sprintf("%d-%s", currentYear, monthParam)

tagihanListrikQuery := databases.DB.Raw(`SELECT listrik.id_listrik,
listrik.total_tagihan_listrik, 
listrik.tanggal_tagihan, listrik.status, warga.nama, 
warga.alamat FROM listrik INNER JOIN warga ON listrik.id_warga = 
warga.id_warga WHERE DATE_FORMAT(listrik.tanggal_tagihan, '%Y-%m') = ?`, formattedDate).Scan(&TagihanListrik)

	if tagihanListrikQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanListrikQuery.Error.Error()})
		return
	}

	if tagihanListrikQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": TagihanListrik})
}