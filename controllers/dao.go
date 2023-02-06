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

func (c *Controller) SetupDaoRoutes() {
	log.Println("Setting up dao routes")
	dao := c.FiberApp.Group("/daos")
	dao.Post("/create", func(ctx *fiber.Ctx) error {
		body := new(BodyCreateDao)
		return fiberutils.PostRequestHandler(ctx, c.CreateDao(), fiberutils.RequestBody{Data: body})
	})

	dao.Get("/", func(ctx *fiber.Ctx) error {
		query := new(QueryGetDaoInfo)
		address := ctx.Query("address")
		query.Address = address
		return fiberutils.GetRequestHandler(ctx, c.GetDaoInfo(), fiberutils.RequestBody{Query: query})

	})

	dao.Post("/add/member", func(ctx *fiber.Ctx) error {
		body := new(BodyAddMember)
		return fiberutils.PostRequestHandler(ctx, c.AddMember(), fiberutils.RequestBody{Data: body})
	})
}

type BodyCreateDao struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	AvatarURL       string `json:"avatar_url" validate:"required"`
	ContractAddress string `json:"contract_address" validate:"required"`
	CreatedBy       string `json:"created_by" validate:"required"`
}

// Create a DAO
// Name, Description, Members, Image URL, ID
// Smart contract deployment as well
func (c *Controller) CreateDao() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyCreateDao)
		body.ContractAddress = strings.ToUpper(body.ContractAddress)
		body.CreatedBy = strings.ToUpper(body.CreatedBy)
		dao := models.Dao{
			ContractAddress: body.ContractAddress,
			Description:     body.Description,
			AvatarURL:       body.AvatarURL,
			Name:            body.Name,
		}
		result := c.DB.Clauses(clause.Returning{}).Create(&dao)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Cannot create Data DAO", result.Error)
		}
		_, err := c.AddMember()(fiberutils.RequestBody{Data: &BodyAddMember{DaoContractAddress: body.ContractAddress, MemberAddress: body.CreatedBy}})
		if err != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error adding user as member", err)
		}
		return dao, nil
	}
}

type BodyAddMember struct {
	DaoContractAddress string `json:"dao_contract_address" validate:"required"`
	MemberAddress      string `json:"member_address" validate:"required"`
}

func (c *Controller) AddMember() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyAddMember)
		body.MemberAddress = strings.ToUpper(body.MemberAddress)
		body.DaoContractAddress = strings.ToUpper(body.DaoContractAddress)
		// Check if a DAO Exists with the address
		var dao models.Dao
		result1 := c.DB.Where("contract_address = ?", body.DaoContractAddress).First(&dao)
		if result1.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Inavlid Data Dao Contract Address", result1.Error)
		}

		// Create user if does not exist
		count := int64(0)
		result2 := c.DB.Model(&models.User{}).Where("Address = ?", body.MemberAddress).Count(&count)
		if result2.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error fetching user", result2.Error)
		}
		if count == 0 {
			user := models.User{
				Address: body.MemberAddress,
			}
			result3 := c.DB.Create(&user)
			if result3.Error != nil {
				return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error creating user", result3.Error)
			}
		}
		// Fetch the user (if newly created) - Can optimize this
		user := models.User{}
		result4 := c.DB.Model(&models.User{}).Where("address = ?", body.MemberAddress).First(&user)
		if result4.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error fetching user", result4.Error)
		}

		err := c.DB.Model(&user).Association("Daos").Append(&dao)

		if err != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error fetching user dao association", err)
		}

		return "Member Added", nil
	}
}

type QueryGetDaoInfo struct {
	Address string
}

func (c *Controller) GetDaoInfo() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		query := data.Query.(*QueryGetDaoInfo)
		query.Address = strings.ToUpper(query.Address)
		dao := models.Dao{}
		result := c.DB.Preload("Members").Preload("Files").Where("contract_address = ?", query.Address).First(&dao)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Cannot find DAO for the given contract address", result.Error)
		}
		return dao, nil
	}
}

// Should include the members of the DAO along with their information
func (c *Controller) GetDaosForAddress() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		return nil, nil
	}
}

func (c *Controller) GetFiles() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		return nil, nil
	}
}
