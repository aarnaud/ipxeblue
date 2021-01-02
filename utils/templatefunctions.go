package utils

import (
	"bytes"
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"text/template"
)

func GetCustomFunctions(c *gin.Context, tpl *template.Template) template.FuncMap {
	config := c.MustGet("config").(*Config)
	return map[string]interface{}{
		"BootentryTemplate": func(name uuid.UUID, data interface{}) (ret string, err error) {
			buf := bytes.NewBuffer([]byte{})
			err = tpl.ExecuteTemplate(buf, name.String(), data)
			ret = buf.String()
			return
		},
		"GetBaseURL": func() (ret string, err error) {
			return config.BaseURL, nil
		},
		"GetDownloadURL": func(bootentry models.Bootentry, filename string) (ret string, err error) {
			return fmt.Sprintf("%s%s", config.BaseURL, bootentry.GetDownloadPath(filename)), err
		},
	}
}
