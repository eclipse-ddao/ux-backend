package controllers

import (
	"net/http"
	"strings"

	"github.com/eclipse-ddao/eclipse-backend/models"
	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/hwnprsd/fiber-utils"
	"gorm.io/gorm/clause"
)

func (c *Controller) SetupFileRoutes() {
	files := c.FiberApp.Group("/files")
	files.Post("/add", func(ctx *fiber.Ctx) error {
		body := new(BodyAddFile)
		return fiberutils.PostRequestHandler(ctx, c.AddFile(), fiberutils.RequestBody{Data: body})
	})
	files.Post("/delete", func(ctx *fiber.Ctx) error {
		body := new(BodyDeleteFile)
		return fiberutils.PostRequestHandler(ctx, c.DeleteFile(), fiberutils.RequestBody{Data: body})
	})
}

type BodyAddFile struct {
	FileName           string `json:"file_name" validate:"required"`
	UploadedBy         string `json:"uploaded_by" validate:"required"`
	DaoContractAddress string `json:"dao_contract_address" validate:"required"`
	ImageURL           string `json:"image_url" validate:"required"`
	Cid                string `json:"cid" validate:"required"`
	FileType           string `json:"file_type" validate:"required"`
}

// File Name, Date, Upload By
func (c *Controller) AddFile() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyAddFile)
		body.DaoContractAddress = strings.ToUpper(body.DaoContractAddress)
		body.UploadedBy = strings.ToUpper(body.UploadedBy)
		// Create a new file entry
		file := models.File{
			Name:               body.FileName,
			UploadedBy:         body.UploadedBy,
			ImageURL:           body.ImageURL,
			DaoContractAddress: body.DaoContractAddress,
			Cid:                body.Cid,
			FileType:           body.FileType,
		}

		// check if the dao exists
		dao := models.Dao{}

		res := c.DB.Where("contract_address = ?", body.DaoContractAddress).First(&dao)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Cannot find DAO with the given contract address", res.Error)
		}

		res = c.DB.Clauses(clause.Returning{}).Create(&file)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error creating a new file entry", res.Error)
		}

		err := c.DB.Model(&dao).Association("Files").Append(&file)
		if err != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error creating a new Association", err)
		}

		return file, nil
	}
}

type BodyDeleteFile struct {
	FileID string `json:"file_id"`
}

func (c *Controller) DeleteFile() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyDeleteFile)
		file := models.File{}
		res := c.DB.Delete(&file, body.FileID)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error Deleting File", res.Error)
		}
		return "Deleted", nil
	}

}
