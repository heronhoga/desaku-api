package admin

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTagihanWifi (c *gin.Context) {

	var WargaWifi []struct {
		IdTagihan string `json:"id_tagihan"`
		TotalTagihanWifi string `json:"total_tagihan_wifi"`
		TanggalTagihan string `json:"tanggal_tagihan"`
		Status string `json:"status"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

	tagihanWifiQuery := databases.DB.Raw(`SELECT tagihan_wifi.id_tagihan, 
	tagihan_wifi.total_tagihan_wifi, tagihan_wifi.tanggal_tagihan, tagihan_wifi.status,
	warga.nama, warga.alamat FROM tagihan_wifi INNER JOIN warga ON
	tagihan_wifi.id_pelanggan = warga.id_warga`).Scan(&WargaWifi)

	if tagihanWifiQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanWifiQuery.Error.Error()})
		return
	}

	if tagihanWifiQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": WargaWifi})

}