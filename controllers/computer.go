package controllers

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// ListComputers
// @Summary List computers
// @Description List of computers filtered or not
// @Accept  json
// @Produce  json
// @Param   _start  query    int     false        "Offset"
// @Success 200 {object} []models.Computer
// @Router /computers [get]
func ListComputers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var total int64
	db = ListFilter(db, c)
	db.Model(&models.Computer{}).Count(&total)
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	computers := make([]models.Computer, 0)
	db = PaginationFilter(db, c)
	db.Preload("Tags").Find(&computers)
	c.JSON(http.StatusOK, computers)
}

// GetComputer
// @Summary Get computer
// @Description Get a computer by Id
// @Accept  json
// @Produce  json
// @Param   id  path     string     true        "Computer UUID" minlength(36) maxlength(36)
// @Success 200 {object} models.Computer
// @Failure 404 {object} models.Error "Computer with uuid %s not found"
// @Router /computers/{id} [get]
func GetComputer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	computer := models.Computer{}
	result := db.Preload("Tags").Where("uuid = ?", id).First(&computer)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("Computer with uuid %s not found", id),
		})
		return
	}
	c.JSON(http.StatusOK, computer)
}

// UpdateComputer
// @Summary Update computer
// @Description Update a computer
// @Accept  json
// @Produce  json
// @Param   id  path     string     true        "Computer UUID" minlength(36) maxlength(36)
// @Param   computer body models.Computer true  "json format computer"
// @Success 200 {object} models.Computer
// @Failure 500 {object} models.Error "Unmarshall error"
// @Failure 400 {object} models.Error "Query ID and UUID miss match"
// @Failure 404 {object} models.Error "Can not find ID"
// @Router /computers/{id} [put]
func UpdateComputer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	computerUpdate := models.Computer{}
	err := c.Bind(&computerUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: err.Error(),
		})
		return
	}
	if computerUpdate.Uuid.String() != id {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "IDs missmatch"})
		return
	}

	result := db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&computerUpdate).Updates(map[string]interface{}{
		"Name":          computerUpdate.Name,
		"Tags":          computerUpdate.Tags,
		"BootentryUUID": computerUpdate.BootentryUUID,
	})

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("Computer with uuid %s not found", id),
		})
		return
	}

	// clean tags not present in updated object
	computer := models.Computer{}
	db.Preload("Tags").First(&computer, "uuid = ?", computerUpdate.Uuid)
	for _, tagInDB := range computer.Tags {
		toDelete := true
		for _, tagToKeep := range computerUpdate.Tags {
			if tagInDB.Key == tagToKeep.Key {
				toDelete = false
			}
		}
		if toDelete {
			result := db.Delete(tagInDB)
			if result.RowsAffected == 0 {
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
					Error: fmt.Sprintf("failed to delete tag %s for computer %s", tagInDB.Key, id),
				})
				return
			}
		}
	}

	// refresh data from DB before return it
	computer = models.Computer{}
	db.Preload("Tags").First(&computer, "uuid = ?", computerUpdate.Uuid)

	c.JSON(http.StatusOK, computer)
}

// DeleteComputer
// @Summary Delete computer
// @Description Delete a computer
// @Accept  json
// @Produce  json
// @Param   id  path     string     true        "Computer UUID" minlength(36) maxlength(36)
// @Success 200
// @Failure 400 {object} models.Error "Failed to parse UUID"
// @Failure 404 {object} models.Error "Can not find ID"
// @Router /computers/{id} [delete]
func DeleteComputer(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := uuid.MustParse(c.Param("id"))

	computer := models.Computer{
		Uuid: id,
	}
	result := db.Delete(&computer)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("Computer with uuid %s not found", id),
		})
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}
