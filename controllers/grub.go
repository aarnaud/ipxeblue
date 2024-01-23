package controllers

import (
	"bytes"
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/aarnaud/ipxeblue/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"net/http"
	"text/template"
)

func GrubScript(c *gin.Context) {
	config := c.MustGet("config").(*utils.Config)

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
		c.HTML(http.StatusOK, "grub_index.gohtml", gin.H{
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
		c.HTML(http.StatusOK, "grub_empty.gohtml", gin.H{})
		return
	}
	bootentry := bootorder.Bootentry

	// Create template name by the uuid
	tpl := template.New(bootentry.Uuid.String())
	// provide a func in the FuncMap which can access tpl to be able to look up templates
	tpl.Funcs(utils.GetCustomFunctions(c, tpl))

	tpl, err = tpl.Parse(bootentry.GrupScript)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	writer := bytes.NewBuffer([]byte{})
	writer.Write([]byte("set timeout=2\n"))
	writer.Write([]byte(fmt.Sprintf("set prefix=(http,%s)\n", config.BaseURL.Host)))
	writer.Write([]byte(fmt.Sprintf("echo 'Booting %s'\n", bootentry.Description)))

	// if bootentry selected is menu load all bootentries as template
	if bootentry.Name == "menu" {
		bootentries := make([]models.Bootentry, 0)
		db.Preload("Files").Where("name != 'menu'").Find(&bootentries)
		for _, be := range bootentries {
			// test if empty
			tpl.New(be.Uuid.String()).Parse(be.GrupScript)
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
	}

	// reset bootentry if not persistent
	if !*bootentry.Persistent {
		db.Model(&bootorder).Delete(&bootorder)
	}

	c.Data(http.StatusOK, "text/plain", writer.Bytes())
}
