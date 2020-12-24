package controllers

import (
	"fmt"
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
	var err error
	var start int = 0
	var order string = "ASC"
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
	if value, existe := c.GetQuery("_order"); existe {
		order = value
	}
	if value, exist := c.GetQuery("_sort"); exist {
		if value == "id" {
			value = "uuid"
		}
		db = db.Order(fmt.Sprintf("%s %s", ToSnakeCase(value), order))
	}

	for q, v := range c.Request.URL.Query() {
		if strings.HasPrefix(q, "_") {
			continue
		}
		if len(v) > 0 {
			value := v[0]
			if strings.Contains(value, "%") {
				db = db.Where(fmt.Sprintf("%s ~~ ?", q), value)
			} else {
				db = db.Where(fmt.Sprintf("%s = ?", q), value)
			}
		}
	}

	return db
}
