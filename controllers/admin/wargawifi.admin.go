package admin

import (
	"desaku-api/databases"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
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

func CreateTagihanWifi (c *gin.Context) {
	month := c.Param("month")
	var TotalIdPelanggan[] struct {
		IdPelanggan string `json:"id_pelanggan"`
	}

	pelangganWifiQuery := databases.DB.Raw(`SELECT id_pelanggan FROM daftar_pelanggan_wifi WHERE status = 'aktif'`).Scan(&TotalIdPelanggan)

	if pelangganWifiQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": pelangganWifiQuery.Error.Error()})
		return
	}

	if pelangganWifiQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	for _, id := range TotalIdPelanggan {
		var DataTagihanWifi struct {
			TotalTagihanWifi string `json:"total_tagihan_wifi"`
			TanggalTagihan string `json:"tanggal_tagihan"`
			IdPelanggan string `json:"id_pelanggan"`
		} 

		DataTagihanWifi.TotalTagihanWifi = "250000"
		DataTagihanWifi.TanggalTagihan = fmt.Sprintf("%d-%s-20", time.Now().Year(), month)
		DataTagihanWifi.IdPelanggan = id.IdPelanggan

		tagihanInsertQuery := databases.DB.Exec("INSERT INTO tagihan_wifi (total_tagihan_wifi, tanggal_tagihan, id_pelanggan) VALUES (?, ?, ?)", DataTagihanWifi.TotalTagihanWifi, DataTagihanWifi.TanggalTagihan, DataTagihanWifi.IdPelanggan)

		if tagihanInsertQuery.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanInsertQuery.Error.Error()})
			return
		}

		if tagihanInsertQuery.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		}
	
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan created successfully"})
}