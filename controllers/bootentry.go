package controllers

import (
	"context"
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/aarnaud/ipxeblue/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
)

//
// @Summary List Bootentries
// @Description List of Bootentry filtered or not
// @Accept  json
// @Produce json
// @Param   _start  query    int     false        "Offset"
// @Success 200 {object} []models.Bootentry
// @Router /bootentries [get]
func ListBootentries(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	bootentries := make([]models.Bootentry, 0)
	db = ListFilter(db, c)
	db.Preload("Files").Find(&bootentries)
	c.Header("X-Total-Count", strconv.Itoa(len(bootentries)))
	c.JSON(http.StatusOK, bootentries)
}

//
// @Summary Get Bootentry
// @Description Get a Bootentry by Id
// @Accept  json
// @Produce json
// @Param   id  path     string     true        "Bootentry UUID" minlength(36) maxlength(36)
// @Success 200 {object} models.Bootentry
// @Failure 404 {object} models.Error "Computer with uuid %s not found"
// @Router /bootentries/{id} [get]
func GetBootentry(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("uuid")

	bootentry := models.Bootentry{}
	result := db.Preload("Files").Where("uuid = ?", id).First(&bootentry)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("Bootentry with uuid %s not found", id),
		})
		return
	}
	c.JSON(http.StatusOK, bootentry)
}

//
// @Summary Create Bootentry
// @Description Create a Bootentry
// @Accept  json
// @Produce json
// @Param   bootentry body models.Bootentry true  "json format Bootentry"
// @Success 200 {object} models.Bootentry
// @Failure 400 {object} models.Error "Failed to create account in DB"
// @Failure 500 {object} models.Error "Unmarshall error"
// @Router /bootentries [post]
func CreateBootentry(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	bootentryToCreate := models.Bootentry{}
	err := c.Bind(&bootentryToCreate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: err.Error(),
		})
		return
	}

	bootentryToCreate.Uuid = uuid.New()
	result := db.Create(&bootentryToCreate)

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: fmt.Sprintf("Bootentry not create, unknown error"),
		})
		return
	}

	// refresh data from DB before return it
	bootentry := models.Bootentry{}
	db.First(&bootentry, "uuid = ?", bootentryToCreate.Uuid)

	c.JSON(http.StatusOK, bootentry)
}

//
// @Summary Update Bootentry
// @Description Update a Bootentry
// @Accept  json
// @Produce json
// @Param   uuid  path     string     true        "Bootentry UUID" minlength(36) maxlength(36)
// @Param   bootentry body models.Bootentry true  "json format of Bootentry"
// @Success 200 {object} models.Bootentry
// @Failure 500 {object} models.Error "Unmarshall error"
// @Failure 400 {object} models.Error "Query uuid and uuid miss match"
// @Failure 404 {object} models.Error "Bootentry UUID not found"
// @Router /bootentries/{username} [put]
func UpdateBootentry(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	config := c.MustGet("config").(*utils.Config)
	filestore := c.MustGet("filestore").(*minio.Client)
	id := c.Param("uuid")

	bootentryUpdate := models.Bootentry{}
	err := c.Bind(&bootentryUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: err.Error(),
		})
		return
	}

	if bootentryUpdate.Uuid.String() != id {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Usernames missmatch"})
		return
	}

	for _, file := range bootentryUpdate.Files {
		if *file.Templatized && !*file.Protected {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "templatized file has to be protected"})
			return
		}
	}

	result := db.Model(&bootentryUpdate).Preload("Files").Updates(bootentryUpdate)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("Bootentry with uuid %s not found", id),
		})
		return
	}
	for _, file := range bootentryUpdate.Files {
		result := db.Model(&file).Updates(file)
		if result.RowsAffected != 1 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
				Error: fmt.Sprintf("Failed to save file, %v", result.Error),
			})
			return
		}
	}

	// clean files not present in updated object
	bootenty := models.Bootentry{}
	db.Preload("Files").First(&bootenty, "uuid = ?", bootentryUpdate.Uuid)
	for _, fileInDB := range bootenty.Files {
		toDelete := true
		for _, fileToKeep := range bootentryUpdate.Files {
			if fileInDB.Name == fileToKeep.Name {
				toDelete = false
			}
		}
		if toDelete {
			if err = filestore.RemoveObject(
				context.Background(), config.MinioConfig.BucketName,
				fileInDB.GetFileStorePath(), minio.RemoveObjectOptions{}); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
					Error: fmt.Sprintf("failed to delete object file %s, %v", fileInDB.GetFileStorePath(), err),
				})
				return
			}
			result := db.Delete(fileInDB)
			if result.RowsAffected == 0 {
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
					Error: fmt.Sprintf("failed to delete file %s for bootentry %s, %v", fileInDB.Name, id, err),
				})
				return
			}
		}
	}

	// refresh data from DB before return it
	bootentry := models.Bootentry{}
	db.Preload("Files").First(&bootentry, "uuid = ?", bootentryUpdate.Uuid)

	c.JSON(http.StatusOK, bootentry)
}

//
// @Summary Delete Bootentry
// @Description Delete Bootentry
// @Accept  json
// @Produce json
// @Param   uuid  path     string     true        "Bootentry UUID" minlength(36) maxlength(36)
// @Success 200
// @Failure 400 {object} models.Error "Failed to parse UUID"
// @Failure 404 {object} models.Error "Bootentry UUID not found"
// @Router /Bootentries/{username} [delete]
func DeleteBootentry(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	config := c.MustGet("config").(*utils.Config)
	filestore := c.MustGet("filestore").(*minio.Client)
	id := uuid.MustParse(c.Param("uuid"))

	err := utils.RemoveRecursive(filestore, config.MinioConfig.BucketName, id.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{
			Error: fmt.Sprintf("failed to delete object file %s, %v", id.String(), err),
		})
		return
	}

	bootentry := models.Bootentry{
		Uuid: id,
	}
	result := db.Preload("Files").Delete(&bootentry)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, models.Error{
			Error: fmt.Sprintf("Bootentry with uuid %s not found", id),
		})
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}

func UploadBootentryFile(c *gin.Context) {
	filestore := c.MustGet("filestore").(*minio.Client)
	config := c.MustGet("config").(*utils.Config)
	id := uuid.MustParse(c.Param("uuid"))

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("get form err: %s\n", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: fmt.Sprintf("get form err: %s", err.Error()),
		})
		return
	}

	filename := filepath.Base(file.Filename)
	bootentryfile := models.BootentryFile{
		Name:          filename,
		BootentryUUID: id,
	}
	filereader, err := file.Open()
	if err != nil {
		fmt.Printf("open file err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: fmt.Sprintf("open file err: %s", err.Error()),
		})
		return
	}
	_, err = filestore.PutObject(context.Background(), config.MinioConfig.BucketName, bootentryfile.GetFileStorePath(),
		filereader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		fmt.Printf("upload file err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: fmt.Sprintf("upload file err: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusAccepted, struct{}{})
}
