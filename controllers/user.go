package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/eclipse-ddao/eclipse-backend/models"
	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/hwnprsd/fiber-utils"
	"gorm.io/gorm/clause"
)

func (c *Controller) SetupUserRoutes() {
	log.Println("Setting up user routes")
	dao := c.FiberApp.Group("/users")
	dao.Post("/create", func(ctx *fiber.Ctx) error {
		body := new(BodyCreateUser)
		return fiberutils.PostRequestHandler(ctx, c.CreateUser(), fiberutils.RequestBody{Data: body})
	})
	dao.Get("/", func(ctx *fiber.Ctx) error {
		address := ctx.Query("address")
		address = strings.ToUpper(address)
		query := new(QueryGetOneUserInfo)
		query.Address = address
		return fiberutils.GetRequestHandler(ctx, c.GetOneUserInfo(), fiberutils.RequestBody{Query: query})
	})
	dao.Post("/addresses", func(ctx *fiber.Ctx) error {
		body := new(BodyGetManyUserInfo)
		return fiberutils.PostRequestHandler(ctx, c.GetManyUserInfo(), fiberutils.RequestBody{Data: body})
	})
}

type BodyCreateUser struct {
	Address   string `json:"address" validate:"required"`
	Username  string `json:"username" validate:"required"`
	AvatarURL string `json:"avatar_url" validate:"required"`
}

// Username, address, avatar url
func (c *Controller) CreateUser() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyCreateUser)
		body.Address = strings.ToUpper(body.Address)
		user := models.User{
			Address:   body.Address,
			Username:  body.Username,
			AvatarURL: body.AvatarURL,
		}
		result := c.DB.Clauses(clause.Returning{}).Create(&user)
		if result.Error != nil {

			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error creating user", result.Error)
		}

		return user, nil
	}
}

type BodyGetManyUserInfo struct {
	Addresses []string `json:"addresses" validate:"required"`
}

// Get a list of user addresses and return the info for all of them
func (c *Controller) GetManyUserInfo() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyGetManyUserInfo)
		var users []models.User
		result := c.DB.Preload("Daos").Where("address IN ?", body.Addresses).Find(&users)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error Fetching Users", result.Error)
		}
		return users, nil
	}
}

type QueryGetOneUserInfo struct {
	Address string `json:"address" validate:"required"`
}

func (c *Controller) GetOneUserInfo() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Query.(*QueryGetOneUserInfo)
		var users models.User
		result := c.DB.Preload("Daos").Where("address = ?", body.Address).First(&users)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error Fetching Users", result.Error)
		}
		return users, nil
	}
}
