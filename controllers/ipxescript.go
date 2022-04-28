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
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func updateOrCreateComputer(c *gin.Context, id uuid.UUID, mac pgtype.Macaddr, ip pgtype.Inet) models.Computer {
	var computer models.Computer
	var err error
	db := c.MustGet("db").(*gorm.DB)

	// auto set name based on hostname or asset for new computer
	name := c.DefaultQuery("hostname", "")
	if name == "" {
		name = c.DefaultQuery("asset", "")
	}

	computer, err = searchComputer(db, id, mac)

	if err != nil {
		computer = models.Computer{
			Name:              name,
			Asset:             c.DefaultQuery("asset", ""),
			BuildArch:         c.DefaultQuery("buildarch", ""),
			Hostname:          c.DefaultQuery("hostname", ""),
			LastSeen:          time.Now(),
			Mac:               mac,
			IP:                ip,
			Manufacturer:      c.DefaultQuery("manufacturer", ""),
			Platform:          c.DefaultQuery("platform", ""),
			Product:           c.DefaultQuery("product", ""),
			Serial:            c.DefaultQuery("serial", ""),
			Uuid:              id,
			Version:           c.DefaultQuery("version", ""),
			LastIpxeaccountID: c.MustGet("account").(*models.Ipxeaccount).Username,
		}
		db.FirstOrCreate(&computer)
	}

	// Uuid may change with Virtual Machine like VMware
	if computer.Uuid != id {
		db.Model(&computer).Where("uuid = ?", computer.Uuid).Update("uuid", id)
	}

	if time.Now().Sub(computer.LastSeen).Seconds() > 10 {
		computer.Asset = c.DefaultQuery("asset", "")
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
		computer.LastIpxeaccountID = c.MustGet("account").(*models.Ipxeaccount).Username
		db.Save(computer)
	}

	return computer
}

func searchComputer(db *gorm.DB, id uuid.UUID, mac pgtype.Macaddr) (models.Computer, error) {
	computer := models.Computer{}
	result := db.Where("uuid = ?", id).First(&computer)
	if result.RowsAffected > 0 {
		return computer, nil
	}
	result = db.Where("mac = ?", mac.Addr.String()).First(&computer)
	if result.RowsAffected > 0 {
		return computer, nil
	}
	return computer, fmt.Errorf("computer not found")
}

func IpxeScript(c *gin.Context) {
	config := c.MustGet("config").(*utils.Config)

	// redirect to admin if it's a browser
	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.Redirect(http.StatusMultipleChoices, "/admin/")
		return
	}

	// basic check or reply with ipxe chain
	_, uuidExist := c.GetQuery("uuid")
	_, macExist := c.GetQuery("mac")
	_, ipExist := c.GetQuery("ip")
	if !uuidExist || !macExist || !ipExist {
		baseURL := *config.BaseURL
		// use the same scheme from request to generate URL
		if schem := c.Request.Header.Get("X-Forwarded-Proto"); schem != "" {
			baseURL.Scheme = schem
		}
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"BaseURL": config.BaseURL.String(),
			"Scheme":  config.BaseURL.Scheme,
			"Host":    config.BaseURL.Host,
		})
		return
	}

	// process query params
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
	// Add computer in gin context to use it in template function
	c.Set("computer", &computer)

	c.Header("Content-Type", "text/plain; charset=utf-8")
	bootorder := models.Bootorder{}
	result := db.Preload("Bootentry").Preload("Bootentry.Files").
		Where("computer_uuid = ?", computer.Uuid).Order("bootorders.order").First(&bootorder)
	if result.RowsAffected == 0 {
		c.HTML(http.StatusOK, "empty.gohtml", gin.H{})
		return
	}
	bootentry := bootorder.Bootentry

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
		db.Model(&bootorder).Delete(&bootorder)
	}

	c.Data(http.StatusOK, "text/plain", writer.Bytes())
}

func DownloadPublicFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	filepath := c.Param("filepath")
	filename := path.Base(filepath)
	subpath := strings.TrimLeft(path.Dir(filepath), "/")
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: err.Error(),
		})
		return
	}

	bootentryFile := models.BootentryFile{
		Name:          filename,
		SubPath:       subpath,
		BootentryUUID: id,
	}

	db.Model(&models.BootentryFile{}).Where("bootentry_uuid = ? AND name = ?", id, filename).First(&bootentryFile)

	if *bootentryFile.Protected {
		c.AbortWithStatusJSON(http.StatusForbidden, models.Error{
			Error: fmt.Sprintf("protected file, you need to use a token URL"),
		})
		return
	}

	Downloadfile(c, &bootentryFile, nil)

}

func DownloadProtectedFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	filepath := c.Param("filepath")
	filename := path.Base(filepath)
	subpath := strings.TrimLeft(path.Dir(filepath), "/")
	tokenString := c.Param("token")
	token := models.Token{}

	db.Preload("Computer").
		Preload("Computer.Tags").
		Preload("Bootentry").
		Preload("Bootentry.Files").
		Where("token = ?", tokenString).First(&token)

	if token.BootentryFile != nil {
		if *token.Filename != filename {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
				Error: fmt.Sprintf("filename not allow with this token"),
			})
			return
		}
	} else {
		for _, file := range token.Bootentry.Files {
			if file.Name == filename && file.SubPath == subpath {
				token.BootentryFile = &file
				token.BootentryFile.Bootentry = &token.Bootentry
				c.Set("computer", &token.Computer)
				break
			}
		}
		// if BootentryFile still null, return not found
		if token.BootentryFile == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
				Error: fmt.Sprintf("file not found"),
			})
			return
		}
	}

	Downloadfile(c, token.BootentryFile, &token.Computer)
}

func Downloadfile(c *gin.Context, bootentryFile *models.BootentryFile, computer *models.Computer) {
	filestore := c.MustGet("filestore").(*minio.Client)
	config := c.MustGet("config").(*utils.Config)

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

	isTemplate := *bootentryFile.Templatized && objectFile.Size < 2*1024*1024

	for header, value := range headers {
		if header == "Content-Length" && isTemplate {
			continue
		}
		c.Header(header, value[0])
	}

	if c.Request.Method == "HEAD" {
		c.Done()
		return
	}

	if isTemplate {
		// Create template name by the uuid
		tpl := template.New(bootentryFile.Name)
		// provide a func in the FuncMap which can access tpl to be able to look up templates
		tpl.Funcs(utils.GetCustomFunctions(c, tpl))
		buf := new(strings.Builder)
		_, err := io.Copy(buf, reader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
				Error: err.Error(),
			})
			return
		}
		tpl, err = tpl.Parse(buf.String())
		writer := bytes.NewBuffer([]byte{})
		err = tpl.ExecuteTemplate(writer, bootentryFile.Name, gin.H{
			"Bootentry": bootentryFile.Bootentry,
			"Computer":  computer,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
				Error: err.Error(),
			})
		}
		// override Content-Length with size after template is render
		c.Header("Content-Length", fmt.Sprintf("%d", writer.Len()))
		c.Data(http.StatusOK, objectFile.ContentType, writer.Bytes())
	} else {
		c.DataFromReader(http.StatusOK, objectFile.Size, objectFile.ContentType, reader, nil)
	}
}
