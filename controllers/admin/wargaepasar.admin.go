package admin

import (
	"net/http"
	"desaku-api/databases"

	"github.com/gin-gonic/gin"
)

func GetAllToko(c *gin.Context) {
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

func GetSpecificToko(c *gin.Context) {

	id := c.Param("id")

	var epasar struct {
		IdToko string `json:"id_toko"`
		NamaToko string `json:"nama_toko"`
		NamaPedagang string `json:"nama_pedagang"`
		JenisDagangan string `json:"jenis_dagangan"`
		Status string `json:"status"`
	}

	tokoQuery := databases.DB.Raw(`SELECT id_toko, nama_toko, nama_pedagang, jenis_dagangan, status FROM epasar WHERE id_toko = ?`, id).Scan(&epasar)

	if tokoQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tokoQuery.Error.Error()})
		return
	}

	if tokoQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": epasar})
}

func EditTokoData(c *gin.Context) {
	id := c.Param("id")

	var requestBody struct {
		NamaToko      *string `json:"nama_toko"`
		NamaPedagang  *string `json:"nama_pedagang"`
		JenisDagangan *string `json:"jenis_dagangan"`
		Status        *string `json:"status"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := `UPDATE epasar SET `
	var updateValues []interface{}

	if requestBody.NamaToko != nil {
		updateQuery += "nama_toko = ?, "
		updateValues = append(updateValues, *requestBody.NamaToko)
	}

	if requestBody.NamaPedagang != nil {
		updateQuery += "nama_pedagang = ?, "
		updateValues = append(updateValues, *requestBody.NamaPedagang)
	}

	if requestBody.JenisDagangan != nil {
		updateQuery += "jenis_dagangan = ?, "
		updateValues = append(updateValues, *requestBody.JenisDagangan)
	}

	if requestBody.Status != nil {
		updateQuery += "status = ?, "
		updateValues = append(updateValues, *requestBody.Status)
	}

	updateQuery = updateQuery[:len(updateQuery)-2] + " WHERE id_toko = ?"
	updateValues = append(updateValues, id)

	tokoQuery := databases.DB.Exec(updateQuery, updateValues...)

	if tokoQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tokoQuery.Error.Error()})
		return
	}

	if tokoQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func DeleteTokoData(c *gin.Context) {
	id := c.Param("id")

	deleteQuery := databases.DB.Exec("DELETE FROM epasar WHERE id_toko = ?", id)
	if deleteQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteQuery.Error.Error()})
		return
	}

	if deleteQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
