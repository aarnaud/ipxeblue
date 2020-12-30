package controllers

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

//
// @Summary List iPXE account
// @Description List of accounts for ipxe
// @Accept  json
// @Produce  json
// @Param   _start  query    int     false        "Offset"
// @Success 200 {object} []models.Ipxeaccount
// @Router /ipxeaccount [get]
func ListIpxeaccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	logins := make([]models.Ipxeaccount, 0)
	db = ListFilter(db, c)
	db.Find(&logins)
	c.Header("X-Total-Count", strconv.Itoa(len(logins)))
	c.JSON(http.StatusOK, logins)
}

//
// @Summary Get iPXE account
// @Description Get iPXE account by username
// @Accept  json
// @Produce  json
// @Param   username  path     string     true        "Username"
// @Success 200 {object} models.Ipxeaccount
// @Failure 404 {object} models.Error "iPXE account not found"
// @Router /ipxeaccount/{username} [get]
func GetIpxeaccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	username := c.Param("username")

	account := models.Ipxeaccount{}
	result := db.Where("username = ?", username).First(&account)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("iPXE account with username %s not found", username),
		})
		return
	}
	c.JSON(http.StatusOK, account)
}

//
// @Summary Create iPXE account
// @Description Create a iPXE account
// @Accept  json
// @Produce  json
// @Param   ipxeaccount body models.Ipxeaccount true  "json format iPXE account"
// @Success 200 {object} models.Ipxeaccount
// @Failure 400 {object} models.Error "Failed to create account in DB"
// @Failure 500 {object} models.Error "Unmarshall error"
// @Router /ipxeaccount [post]
func CreateIpxeaccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	accountToCreate := models.Ipxeaccount{}
	err := c.Bind(&accountToCreate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: err.Error(),
		})
		return
	}

	result := db.Create(&accountToCreate)

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: fmt.Sprintf("iPXE account not create, unknown error"),
		})
		return
	}

	// refresh data from DB before return it
	account := models.Ipxeaccount{}
	db.First(&account, "username = ?", accountToCreate.Username)

	c.JSON(http.StatusOK, account)
}

//
// @Summary Update iPXE account
// @Description Update a iPXE account
// @Accept  json
// @Produce  json
// @Param   username  path     string     true        "Username"
// @Param   ipxeaccount body models.Ipxeaccount true  "json format iPXE account"
// @Success 200 {object} models.Ipxeaccount
// @Failure 500 {object} models.Error "Unmarshall error"
// @Failure 400 {object} models.Error "Query username and username miss match"
// @Failure 404 {object} models.Error "iPXE account not found"
// @Router /ipxeaccount/{username} [put]
func UpdateIpxeaccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")

	accountUpdate := models.Ipxeaccount{}
	err := c.Bind(&accountUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: err.Error(),
		})
		return
	}

	if accountUpdate.Username != username {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Usernames missmatch"})
		return
	}

	result := db.Model(&accountUpdate).Updates(accountUpdate)

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("iPXE account with username %s not found", username),
		})
		return
	}

	// refresh data from DB before return it
	account := models.Ipxeaccount{}
	db.First(&account, "username = ?", accountUpdate.Username)

	c.JSON(http.StatusOK, account)
}

//
// @Summary Delete iPXE account
// @Description Delete a iPXE account
// @Accept  json
// @Produce  json
// @Param   username  path     string     true        "Username"
// @Success 200 {object} models.Ipxeaccount
// @Failure 404 {object} models.Error "iPXE account not found"
// @Router /ipxeaccount/{username} [delete]
func DeleteIpxeaccount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")

	account := models.Ipxeaccount{
		Username: username,
	}
	result := db.Delete(&account)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("iPXE account with username %s not found", username),
		})
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}
