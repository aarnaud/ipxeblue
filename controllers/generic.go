package controllers

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/aarnaud/ipxeblue/utils/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ListFilter(db *gorm.DB, c *gin.Context) *gorm.DB {
	var order string = "ASC"
	if value, existe := c.GetQuery("_order"); existe {
		order = value
	}
	if value, exist := c.GetQuery("_sort"); exist {
		// react-admin use id as primary key, so we convert the key depends of object
		value = ConvertReactAdminID(c, value)
		db = db.Order(fmt.Sprintf("%s %s", ToSnakeCase(value), order))
	}

	if value, exist := c.GetQuery("q"); exist {
		query := strings.Builder{}
		fields := SearchFields(c)
		queryvalues := make([]interface{}, len(fields))
		for i, field := range fields {
			if query.Len() != 0 {
				query.WriteString(" OR ")
			}
			queryvalues[i] = fmt.Sprintf("%%%s%%", value)
			query.WriteString(fmt.Sprintf("%s ILIKE ?", field))
		}
		db = db.Where(query.String(), queryvalues...)

	} else {
		for q, v := range c.Request.URL.Query() {
			if strings.HasPrefix(q, "_") {
				continue
			}
			if len(v) == 1 {
				value := v[0]
				// react-admin use id as primary key, so we convert the key depends of object
				q = ConvertReactAdminID(c, q)
				if helpers.StringToType(value) != helpers.TYPE_BOOL && helpers.StringToType(value) != helpers.TYPE_UUID &&
					q != "ip" && q != "mac" {
					db = db.Where(fmt.Sprintf("%s ILIKE ?", q), fmt.Sprintf("%%%s%%", value))
				} else {
					db = db.Where(fmt.Sprintf("%s = ?", q), value)
				}
			}
			if len(v) > 1 {
				db = db.Where(fmt.Sprintf("%s IN ?", q), v)
			}
		}
	}

	return db
}

func PaginationFilter(db *gorm.DB, c *gin.Context) *gorm.DB {
	var err error
	var start int = 0
	if value, exist := c.GetQuery("_start"); exist {
		start, err = strconv.Atoi(value)
		if err == nil {
			db = db.Offset(start)
		}
		//todo: log of failed
	}
	if value, exist := c.GetQuery("_end"); exist {
		if end, err := strconv.Atoi(value); err == nil {
			db = db.Limit(end - start)
		}
		//todo: log of failed
	}
	return db
}

func ConvertReactAdminID(c *gin.Context, key string) string {
	if key == "id" && strings.Contains(c.FullPath(), "/computers") {
		return "uuid"
	}
	if key == "id" && strings.Contains(c.FullPath(), "/bootentries") {
		return "uuid"
	}
	if key == "id" && strings.Contains(c.FullPath(), "/ipxeaccounts") {
		return "username"
	}
	return key
}

func SearchFields(c *gin.Context) []string {
	if strings.Contains(c.FullPath(), "/computers") {
		return models.ComputerSearchFields
	}
	if strings.Contains(c.FullPath(), "/bootentries") {
		return models.BootentrySearchFields
	}
	if strings.Contains(c.FullPath(), "/ipxeaccounts") {
		return models.IpxeAccountSearchFields
	}
	return []string{
		"id",
	}
}
