package midlewares

import (
	"encoding/base64"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func BasicAuthIpxeAccount(onlyAdmin bool) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	return func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)

		auth := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		account := models.Ipxeaccount{}
		var result *gorm.DB
		if onlyAdmin {
			result = db.First(&account, "username = ? AND is_admin = TRUE", pair[0])
		} else {
			result = db.First(&account, "username = ?", pair[0])
		}

		if result.RowsAffected == 0 {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(pair[1])); err != nil {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// update last login field if it's older than 5min
		if time.Now().Sub(account.LastLogin).Seconds() > 300 {
			account.LastLogin = time.Now()
			db.Model(&account).Updates(account)
		}
	}
}
