package admin

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllPajak (c *gin.Context) {
	year := c.Param("year")

	var WargaPajak []struct {
		IdPajak string `json:"id_pajak"`
		Tahun string `json:"tahun"`
		TotalTagihanPajak string `json:"total_tagihan_pajak"`
		StatusBayar string `json:"status_bayar"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

	tagihanPajakQuery := databases.DB.Raw(`SELECT pajak.id_pajak,
	pajak.tahun, 
	pajak.total_tagihan_pajak, pajak.status_bayar, warga.nama, 
	warga.alamat FROM pajak INNER JOIN warga ON pajak.id_warga = 
	warga.id_warga WHERE YEAR(pajak.tahun) = ?`, year).Scan(&WargaPajak)

	if tagihanPajakQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanPajakQuery.Error.Error()})
		return
	}

	if tagihanPajakQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": WargaPajak})
}

func GetSpecificPajak(c *gin.Context) {
	id := c.Param("id") //ID PAJAK

	var WargaPajak struct {
		IdPajak string `json:"id_pajak"`
		Tahun string `json:"tahun"`
		TotalTagihanPajak string `json:"total_tagihan_pajak"`
		StatusBayar string `json:"status_bayar"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

	tagihanPajakQuery := databases.DB.Raw(`SELECT pajak.id_pajak,
	pajak.tahun, 
	pajak.total_tagihan_pajak, pajak.status_bayar, warga.nama, 
	warga.alamat FROM pajak INNER JOIN warga ON pajak.id_warga = 
	warga.id_warga WHERE id_pajak = ?`, id).Scan(&WargaPajak)

	if tagihanPajakQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanPajakQuery.Error.Error()})
		return
	}

	if tagihanPajakQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": WargaPajak})
}

func CreateTagihanPajak (c *gin.Context) {
	year := c.Param("year")

	var WargaList []struct {
        IdWarga string `gorm:"column:id_warga"`
    }

    wargaQuery := databases.DB.Raw(`SELECT id_warga FROM warga`).Scan(&WargaList)
    if wargaQuery.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": wargaQuery.Error.Error()})
        return
    }

	for _, warga := range WargaList {
		pajak := struct {
			IdWarga string `gorm:"column:id_warga"`
			Tahun string `gorm:"column:tahun"`
			StatusBayar string `gorm:"column:status_bayar"`
		} {
			IdWarga: warga.IdWarga,
			Tahun: year,
			StatusBayar: "pending",
		}

		result := databases.DB.Exec(`INSERT INTO pajak (id_warga, tahun, status_bayar) VALUES (?, ?, ?)`, pajak.IdWarga, pajak.Tahun, pajak.StatusBayar)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pajak created successfully"})

}

func UpdateTagihanPajak (c *gin.Context) {
    id := c.Param("id")

    var reqBody struct {
        TotalTagihanPajak string `json:"total_tagihan_pajak"`
    }

    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result := databases.DB.Exec(`UPDATE pajak SET total_tagihan_pajak = ? WHERE id_pajak = ?`, reqBody.TotalTagihanPajak, id)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Pajak updated successfully"})
}

