package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/aarnaud/ipxeblue/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func updateOrCreateComputer(c *gin.Context, id uuid.UUID, mac pgtype.Macaddr, ip pgtype.Inet) models.Computer {
	db := c.MustGet("db").(*gorm.DB)

	computer := models.Computer{
		Asset:        c.DefaultQuery("asset", ""),
		BuildArch:    c.DefaultQuery("buildarch", ""),
		Hostname:     c.DefaultQuery("hostname", ""),
		LastSeen:     time.Now(),
		Mac:          mac,
		IP:           ip,
		Manufacturer: c.DefaultQuery("manufacturer", ""),
		Platform:     c.DefaultQuery("platform", ""),
		Product:      c.DefaultQuery("product", ""),
		Serial:       c.DefaultQuery("serial", ""),
		Uuid:         id,
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
		computer.IP = ip
		computer.Manufacturer = c.DefaultQuery("manufacturer", "")
		computer.Platform = c.DefaultQuery("platform", "")
		computer.Product = c.DefaultQuery("product", "")
		computer.Serial = c.DefaultQuery("serial", "")
		computer.Version = c.DefaultQuery("version", "")
		db.Save(computer)
	}

	return computer
}

func IpxeScript(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, err := uuid.Parse(c.Query("uuid"))
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

	ip := pgtype.Inet{}
	err = ip.DecodeText(nil, []byte(c.Query("ip")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": err.Error(),
		})
		return
	}

	computer := updateOrCreateComputer(c, id, mac, ip)

	c.Header("Content-Type", "text/plain; charset=utf-8")
	bootentry := models.Bootentry{}
	result := db.Preload("Files").Where("uuid = ?", computer.BootentryUUID).First(&bootentry)
	if result.RowsAffected == 0 {
		c.HTML(http.StatusOK, "empty.gohtml", gin.H{})
		return
	}

	// Create template name by the uuid
	tpl := template.New(bootentry.Uuid.String())
	// provide a func in the FuncMap which can access tpl to be able to look up templates
	tpl.Funcs(utils.GetCustomFunctions(c, tpl))

	tpl, err = tpl.Parse(bootentry.IpxeScript)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	writer := bytes.NewBuffer([]byte{})
	writer.Write([]byte("#!ipxe\n"))
	writer.Write([]byte(fmt.Sprintf("echo Booting %s\n", bootentry.Description)))

	// if bootentry selected is menu load all bootentries as template
	if bootentry.Name == "menu" {
		bootentries := make([]models.Bootentry, 0)
		db.Preload("Files").Where("name != 'menu'").Find(&bootentries)
		for _, be := range bootentries {
			tpl.New(be.Uuid.String()).Parse(be.IpxeScript)
		}
		err = tpl.ExecuteTemplate(writer, bootentry.Uuid.String(), gin.H{
			"Computer":    computer,
			"Bootentry":   bootentry,
			"Bootentries": bootentries,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		err = tpl.ExecuteTemplate(writer, bootentry.Uuid.String(), bootentry)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// add failed goto that can be use in ipxescript
		writer.Write([]byte("\n\n:failed\necho Booting failed, waiting 10 sec\nsleep 10\nexit 1"))
	}

	// reset bootentry if not persistent
	if !*bootentry.Persistent {
		db.Model(&computer).Updates(map[string]interface{}{
			"BootentryUUID": nil,
		})
	}

	c.Data(http.StatusOK, "text/plain", writer.Bytes())
}

func DownloadFiles(c *gin.Context) {
	filestore := c.MustGet("filestore").(*minio.Client)
	config := c.MustGet("config").(*utils.Config)
	filename := c.Param("filename")
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: err.Error(),
		})
		return
	}

	bootentryFile := models.BootentryFile{
		Name:          filename,
		BootentryUUID: id,
	}

	getObjectOptions := minio.GetObjectOptions{}
	if byterange := c.Request.Header.Get("Range"); byterange != "" {
		rangesplit := strings.Split(byterange, "=")
		rangevalue := strings.Split(rangesplit[1], "-")
		start, err := strconv.ParseInt(rangevalue[0], 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
				Error: err.Error(),
			})
			return
		}
		end, err := strconv.ParseInt(rangevalue[1], 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
				Error: err.Error(),
			})
			return
		}
		err = getObjectOptions.SetRange(start, end)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
				Error: err.Error(),
			})
			return
		}
	}

	reader, objectFile, headers, err := minio.Core{filestore}.GetObject(context.Background(),
		config.MinioConfig.BucketName, bootentryFile.GetFileStorePath(), getObjectOptions)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: err.Error(),
		})
		return
	}

	for header, value := range headers {
		c.Header(header, value[0])
	}
	c.DataFromReader(200, objectFile.Size, objectFile.ContentType, reader, nil)
}
