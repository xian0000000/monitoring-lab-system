package handlers

import (
	"lab-monitor-backend/db"
	"lab-monitor-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ===== LAB HANDLERS =====

func GetLabs(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, location, capacity FROM labs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Tampung semua lab ke slice
	labs := []models.Lab{}
	for rows.Next() {
		var lab models.Lab
		rows.Scan(&lab.ID, &lab.Name, &lab.Location, &lab.Capacity)
		labs = append(labs, lab)
	}

	c.JSON(http.StatusOK, labs)
}

func CreateLab(c *gin.Context) {
	var lab models.Lab

	// Baca data dari request body (JSON)
	if err := c.ShouldBindJSON(&lab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Simpan ke database
	err := db.DB.QueryRow(
		"INSERT INTO labs (name, location, capacity) VALUES ($1, $2, $3) RETURNING id",
		lab.Name, lab.Location, lab.Capacity,
	).Scan(&lab.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lab)
}

// ===== DEVICE HANDLERS =====

func GetDevices(c *gin.Context) {
	rows, err := db.DB.Query(
		"SELECT id, name, ip_address, lab_id, status, last_seen FROM devices",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	devices := []models.Device{}
	for rows.Next() {
		var device models.Device
		rows.Scan(
			&device.ID,
			&device.Name,
			&device.IPAddress,
			&device.LabID,
			&device.Status,
			&device.LastSeen,
		)
		devices = append(devices, device)
	}

	c.JSON(http.StatusOK, devices)
}

func CreateDevice(c *gin.Context) {
	var device models.Device

	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO devices (name, ip_address, lab_id) VALUES ($1, $2, $3) RETURNING id",
		device.Name, device.IPAddress, device.LabID,
	).Scan(&device.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}

func DeleteDevice(c *gin.Context) {
	// Ambil id dari URL, contoh: DELETE /api/devices/3
	id := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM devices WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device berhasil dihapus"})
}