package controllers

import (
	"net/http"
	"strings"

	"github.com/eclipse-ddao/eclipse-backend/models"
	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/hwnprsd/fiber-utils"
	"gorm.io/gorm/clause"
)

func (c *Controller) SetupStorageProviderRoutes() {

	sp := c.FiberApp.Group("/storage-provider")
	sp.Post("/create", func(ctx *fiber.Ctx) error {
		data := new(BodyCreateStorageProvider)
		data.Address = strings.ToUpper(data.Address)
		return fiberutils.PostRequestHandler(ctx, c.CreateStorageProvider(), fiberutils.RequestBody{Data: data})
	})
	sp.Get("/", func(ctx *fiber.Ctx) error {
		query := new(QueryGetStroageProviderInfo)
		query.StorageProviderAddress = ctx.Query("address")
		query.StorageProviderAddress = strings.ToUpper(query.StorageProviderAddress)
		return fiberutils.GetRequestHandler(ctx, c.GetStorageProvderInfo(), fiberutils.RequestBody{Query: query})
	})
}

type BodyCreateStorageProvider struct {
	Address        string `json:"address" validate:"required"`
	ContactInfo    string `json:"contact_info" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
	ReputationLink string `json:"reputation_link"`
	AvatarURL      string `json:"avatar_url"`
}

func (c *Controller) CreateStorageProvider() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyCreateStorageProvider)
		body.Address = strings.ToUpper(body.Address)
		storageProvider := models.StorageProvider{
			Address:        body.Address,
			ContactInfo:    body.ContactInfo,
			Name:           body.Name,
			Description:    body.Description,
			ReputationLink: body.ReputationLink,
			AvatarURL:      body.AvatarURL,
		}
		res := c.DB.Clauses(clause.Returning{}).Create(&storageProvider)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error creating storage provider", res.Error)
		}
		return storageProvider, nil
	}
}

type QueryGetStroageProviderInfo struct {
	StorageProviderAddress string
}

func (c *Controller) GetStorageProvderInfo() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		query := data.Query.(*QueryGetStroageProviderInfo)
		sp := models.StorageProvider{}
		res := c.DB.Preload("Proposals.BFile").Where("address = ?", query.StorageProviderAddress).First(&sp)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid sp address", res.Error)
		}
		return sp, nil
	}
}

func (c *Controller) IsStorageProviderValid(address string) (*models.StorageProvider, error) {
	sp := models.StorageProvider{}
	res := c.DB.Where("address = ?", address).First(&sp)
	if res.Error != nil {
		return nil, res.Error
	}
	return &sp, nil
}
