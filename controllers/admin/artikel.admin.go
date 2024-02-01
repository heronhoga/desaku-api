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

func CreateArtikel(c *gin.Context) {

	var input struct {
		Judul string `json:"judul"`
		Isi string `json:"isi"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := databases.DB.Exec(`INSERT INTO artikel (judul, isi) VALUES (?, ?)`, input.Judul, input.Isi)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel created successfully"})
}

func UpdateArtikel(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Judul *string `json:"judul"`
		Isi   *string `json:"isi"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Judul != nil {
		result := databases.DB.Exec(`UPDATE artikel SET judul = ? WHERE id_artikel = ?`, *input.Judul, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	if input.Isi != nil {
		result := databases.DB.Exec(`UPDATE artikel SET isi = ? WHERE id_artikel = ?`, *input.Isi, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel updated successfully"})
}

func DeleteArtikel(c *gin.Context) {
	id := c.Param("id")

	result := databases.DB.Exec(`DELETE FROM artikel WHERE id_artikel = ?`, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel deleted successfully"})
}