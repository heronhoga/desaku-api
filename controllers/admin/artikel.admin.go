package admin

import (
	"net/http"
	"desaku-api/databases"
	"github.com/gin-gonic/gin"

)

func GetAllArtikel(c *gin.Context) {
	var artikel[] struct {
		IdArtikel string `json:"id_artikel"`
		Judul string `json:"judul"`
		Isi string `json:"isi"`
	}

	artikelQuery := databases.DB.Raw(`SELECT id_artikel, judul, isi FROM artikel`).Scan(&artikel)

	if artikelQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": artikelQuery.Error.Error()})
		return
	}

	if artikelQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": artikel})
}

func GetSpecificArtikel(c *gin.Context) {

	id := c.Param("id")

	var artikel struct {
		IdArtikel string `json:"id_artikel"`
		Judul string `json:"judul"`
		Isi string `json:"isi"`
	}

	artikelQuery := databases.DB.Raw(`SELECT id_artikel, judul, isi FROM artikel WHERE id_artikel = ?`, id).Scan(&artikel)

	if artikelQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": artikelQuery.Error.Error()})
		return
	}

	if artikelQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": artikel})
}