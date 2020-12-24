package controllers

import (
	"github.com/aarnaud/ipxeblue/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// GET /computers
// Get all computer
func Index(c *gin.Context) {
	uuid, err := uuid.Parse(c.Query("uuid"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": err.Error(),
		})
		return
	}

	mac := pgtype.Macaddr{}
	err = mac.DecodeText(nil, []byte(c.Query("mac")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": err.Error(),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	computer := models.Computer{
		Asset:        c.DefaultQuery("asset", ""),
		BuildArch:    c.DefaultQuery("buildarch", ""),
		Hostname:     c.DefaultQuery("hostname", ""),
		LastSeen:     time.Now(),
		Mac:          mac,
		Manufacturer: c.DefaultQuery("manufacturer", ""),
		Platform:     c.DefaultQuery("platform", ""),
		Product:      c.DefaultQuery("product", ""),
		Serial:       c.DefaultQuery("serial", ""),
		Uuid:         uuid,
		Version:      c.DefaultQuery("version", ""),
	}
	db.FirstOrCreate(&computer)
	if time.Now().Sub(computer.LastSeen).Seconds() > 10 {
		computer.Asset = c.DefaultQuery("asset", "")
		computer.BuildArch = c.DefaultQuery("buildarch", "")
		computer.BuildArch = c.DefaultQuery("buildarch", "")
		computer.Hostname = c.DefaultQuery("hostname", "")
		computer.LastSeen = time.Now()
		computer.Mac = mac
		computer.Manufacturer = c.DefaultQuery("manufacturer", "")
		computer.Platform = c.DefaultQuery("platform", "")
		computer.Product = c.DefaultQuery("product", "")
		computer.Serial = c.DefaultQuery("serial", "")
		computer.Version = c.DefaultQuery("version", "")
		db.Save(computer)
	}
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.HTML(http.StatusOK, "index.tmpl", nil)
}
