package admin

import (
	"desaku-api/databases"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func CreateTagihanListrik(c *gin.Context) {
    var WargaList []struct {
        IdWarga string `gorm:"column:id_warga"`
    }

    wargaQuery := databases.DB.Raw(`SELECT id_warga FROM warga`).Scan(&WargaList)
    if wargaQuery.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": wargaQuery.Error.Error()})
        return
    }

    currentYear := time.Now().Year()
    monthParam := c.Param("month")

    for _, warga := range WargaList {
        tagihan := struct {
            IdWarga string `gorm:"column:id_warga"`
            TanggalTagihan string `gorm:"column:tanggal_tagihan"`
        } {
            IdWarga: warga.IdWarga,
            TanggalTagihan: fmt.Sprintf("%d-%s-01", currentYear, monthParam),
        }

        result := databases.DB.Exec(`INSERT INTO listrik (id_warga, tanggal_tagihan) VALUES (?, ?)`, tagihan.IdWarga, tagihan.TanggalTagihan)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Tagihan listrik created for all Warga"})
}
