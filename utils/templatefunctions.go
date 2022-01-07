package utils

import (
	"bytes"
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"text/template"
)

func GetCustomFunctions(c *gin.Context, tpl *template.Template) template.FuncMap {
	config := c.MustGet("config").(*Config)
	db := c.MustGet("db").(*gorm.DB)
	baseURL := *config.BaseURL

	// use the same scheme from request to generate URL
	if schem := c.Request.Header.Get("X-Forwarded-Proto"); schem != "" {
		baseURL.Scheme = schem
	}
	return map[string]interface{}{
		"BootentryTemplate": func(name uuid.UUID, data interface{}) (ret string, err error) {
			buf := bytes.NewBuffer([]byte{})
			err = tpl.ExecuteTemplate(buf, name.String(), data)
			ret = buf.String()
			return
		},
		"GetBaseURL": func() (ret string, err error) {
			return baseURL.String(), nil
		},
		"GetDownloadURL": func(bootentry models.Bootentry, filename string) (ret string, err error) {
			file := bootentry.GetFile(filename)
			if file == nil {
				return fmt.Sprintf("%s not found in bootentry %s", filename, bootentry.Uuid), err
			}
			path, token := file.GetDownloadPath()
			if token != nil {
				// Get computer in gin context to add it in token, to used it in file template.
				token.Computer = *c.MustGet("computer").(*models.Computer)
				db.Create(&token)
			}
			return fmt.Sprintf("%s%s", baseURL.String(), path), err
		},
		"GetDownloadBaseURL": func(bootentry models.Bootentry) (ret string, err error) {
			path, token := bootentry.GetDownloadBasePath()
			if token != nil {
				// Get computer in gin context to add it in token, to used it in file template.
				token.Computer = *c.MustGet("computer").(*models.Computer)
				db.Create(&token)
			}
			return fmt.Sprintf("%s%s", baseURL.String(), path), err
		},
		"GetTagValue": func(name string) (ret string, err error) {
			computer := c.MustGet("computer").(*models.Computer)
			if computer == nil {
				return "", fmt.Errorf("failed to GetTagValue because computer is nil")
			}
			for _, tag := range computer.Tags {
				if tag.Key == name {
					return tag.Value, err
				}
			}
			return "", err
		},
		"GetComputerName": func() (ret string, err error) {
			computer := c.MustGet("computer").(*models.Computer)
			if computer == nil {
				return "", fmt.Errorf("failed to GetTagValue because computer is nil")
			}
			return computer.Name, nil
		},
	}
}
