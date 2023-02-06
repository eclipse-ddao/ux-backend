package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-ddao/eclipse-backend/models"
	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/hwnprsd/fiber-utils"
	"gorm.io/gorm/clause"
)

func (c *Controller) SetupBigFileRoutes() {
	bigFiles := c.FiberApp.Group("big-files")

	bigFiles.Post("/add", func(ctx *fiber.Ctx) error {
		data := new(BodyAddBigFile)
		return fiberutils.PostRequestHandler(ctx, c.AddBigFile(), fiberutils.RequestBody{Data: data})
	})

	bigFiles.Get("/dao", func(ctx *fiber.Ctx) error {
		query := new(QueryGetBigFileForDao)
		address := ctx.Query("address")
		address = strings.ToUpper(address)
		query.ContractAddress = address
		return fiberutils.GetRequestHandler(ctx, c.GetBigFilesForDao(), fiberutils.RequestBody{Query: query})
	})

	bigFiles.Get("/status/:status", func(ctx *fiber.Ctx) error {
		status := ctx.Params("status")
		query := new(QueryGetAllBigFiles)
		query.Status = status
		return fiberutils.GetRequestHandler(ctx, c.GetAllOpenBigFiles(), fiberutils.RequestBody{Query: query})
	})

	bigFiles.Post("/apply", func(ctx *fiber.Ctx) error {
		data := new(BodyApplyAsStorageProvider)
		return fiberutils.PostRequestHandler(ctx, c.ApplyAsStorageProvider(), fiberutils.RequestBody{Data: data})
	})

	bigFiles.Get("/", func(ctx *fiber.Ctx) error {
		query := new(QueryGetBigFileInfo)
		query.BigFileID = ctx.Query("id")
		return fiberutils.GetRequestHandler(ctx, c.GetBigFileInfo(), fiberutils.RequestBody{Query: query})
	})

	bigFiles.Post("/select-proposal", func(ctx *fiber.Ctx) error {
		data := new(BodySelectStorageProvider)
		return fiberutils.PostRequestHandler(ctx, c.SelectStorageProviderForBigFile(), fiberutils.RequestBody{Data: data})
	})

	bigFiles.Post("/accept-proposal", func(ctx *fiber.Ctx) error {
		data := new(BodyAcceptDealAsStorageProvider)
		return fiberutils.PostRequestHandler(ctx, c.AcceptDealAsStorageProvider(), fiberutils.RequestBody{Data: data})
	})
}

type BodyAddBigFile struct {
	Duration           uint   `json:"duration" validate:"required"`
	SizeInGb           uint   `json:"size_in_gb" validate:"required"`
	BaseBounty         uint   `json:"base_bounty" validate:"required"`
	FileType           uint   `json:"file_type" validate:"required"` // 1 = Link, 2 = Hardware Delivery, 3 = Filecoin CID
	Name               string `json:"name" validate:"required"`
	Description        string `json:"description" validate:"required"`
	DaoContractAddress string `json:"dao_contract_address" validate:"required"`
	UploadedBy         string `json:"uploaded_by" validate:"required"`
	Expiry             string `json:"expiry" validate:"required"`
}

// DAO Contract Address, File Information!!
// Deal Duration, Deal Size, Base Bounty, FileType [Link, Hardware, Filecoin CID]
func (c *Controller) AddBigFile() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyAddBigFile)
		body.UploadedBy = strings.ToUpper(body.UploadedBy)
		body.DaoContractAddress = strings.ToUpper(body.DaoContractAddress)
		expiry, err := time.Parse(time.RFC3339, body.Expiry)
		if err != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error parsing time. Please use 'DD-MM-YYYY' format", err)
		}
		bigFile := models.BigFile{
			Duration:           body.Duration,
			SizeInGb:           body.SizeInGb,
			BaseBounty:         body.BaseBounty,
			FileType:           body.FileType,
			Name:               body.Name,
			Description:        body.Description,
			DaoContractAddress: body.DaoContractAddress,
			UploadedBy:         body.UploadedBy,
			Status:             1,
			Expiry:             expiry,
		}
		res := c.DB.Clauses(clause.Returning{}).Create(&bigFile)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error creating a new big-file entry", res.Error)
		}
		return bigFile, nil
	}
}

type QueryGetBigFileForDao struct {
	ContractAddress string `json:"contract_address"`
}

func (c *Controller) GetBigFilesForDao() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		query := data.Query.(*QueryGetBigFileForDao)
		dao := models.Dao{}
		res := c.DB.Preload("BFiles").Where("contract_address = ?", query.ContractAddress).First(&dao)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error finding big files for contract address", res.Error)
		}
		return dao.BFiles, nil
	}
}

type QueryGetAllBigFiles struct {
	Status string `json:"status"`
}

func (c *Controller) GetAllOpenBigFiles() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		query := data.Query.(*QueryGetAllBigFiles)
		bigFiles := []models.BigFile{}
		status, err := strconv.Atoi(query.Status)
		if err != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid Status Passed", errors.New("Pass only 1 or 2 or 3"))
		}
		// If status is 0, then give back all big files
		if status == 0 {
			res := c.DB.Find(&bigFiles)
			if res.Error != nil {
				return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error finding big files for contract address", res.Error)
			}
			return bigFiles, nil
		}
		res := c.DB.Where("status = ?", query.Status).Find(&bigFiles)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error finding big files for contract address", res.Error)
		}
		return bigFiles, nil
	}
}

type QueryGetBigFileInfo struct {
	BigFileID string `json:"big_file_id"`
}

func (c *Controller) GetBigFileInfo() fiberutils.GetHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		query := data.Query.(*QueryGetBigFileInfo)
		bigFile := models.BigFile{}
		res := c.DB.Model(&bigFile).Preload("Proposals").Where("id = ?", query.BigFileID).First(&bigFile)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Error finding big files for contract address", res.Error)
		}
		return bigFile, nil
	}
}

type BodyApplyAsStorageProvider struct {
	RequestedBounty        uint   `json:"requested_bounty" validate:"required"`
	StorageProviderAddress string `json:"storage_provider_address" validate:"required"`
	BFileId                uint   `json:"b_file_id" validate:"required"`
}

// Proposal Bounty, Reputation Link, Contact Information
// Create a proposal
func (c *Controller) ApplyAsStorageProvider() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyApplyAsStorageProvider)
		body.StorageProviderAddress = strings.ToUpper(body.StorageProviderAddress)
		_, err := c.IsStorageProviderValid(body.StorageProviderAddress)
		if err != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid storage provider", err)
		}
		bigFileProposal := models.BigFileProposal{
			RequestedBounty:        body.RequestedBounty,
			IsSelected:             false,
			StorageProviderAddress: body.StorageProviderAddress,
			BFileID:                body.BFileId,
		}
		res := c.DB.Clauses(clause.Returning{}).Create(&bigFileProposal)
		if res.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error creating a new proposal", res.Error)
		}
		return bigFileProposal, nil
	}
}

type BodySelectStorageProvider struct {
	ProposalID uint `json:"proposal_id" validate:"required"`
}

// Make Deal - BigFile Creator accepts a SP
// TODO: Add checks to make sure only the owner can do this
func (c *Controller) SelectStorageProviderForBigFile() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodySelectStorageProvider)

		proposal := models.BigFileProposal{}
		result := c.DB.Preload("BFile").Where("id = ?", body.ProposalID).First(&proposal)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid proposal Id", result.Error)
		}
		// TODO
		if proposal.BFile == nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid BigFile Status", result.Error)
		}
		if proposal.BFile.Status != 1 {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Big File is not in OPEN status", errors.New("Big File is not in Open Status"))
		}

		result = c.DB.Model(&proposal).Clauses(clause.Returning{}).Where("id = ?", body.ProposalID).Update("is_selected", true)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error updating proposal", result.Error)
		}
		if result.RowsAffected != 1 {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid big file id", errors.New("Invalid big file id"))
		}
		// Update the BigFile Status
		bigFile := models.BigFile{}
		result = c.DB.Model(&bigFile).Where("id = ?", proposal.BFileID).Update("status", 2)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error updating Big File Status", result.Error)
		}
		return proposal, nil
	}
}

type BodyAcceptDealAsStorageProvider struct {
	DealCID                string `json:"deal_cid" validate:"required"`
	ProposalID             uint   `json:"proposal_id" validate:"required"`
	StorageProviderAddress string `json:"storage_provider_address" validate:"required"`
}

// Deal CID (needs to be verified by the smart contract)
func (c *Controller) AcceptDealAsStorageProvider() fiberutils.PostHandler {
	return func(data fiberutils.RequestBody) (interface{}, error) {
		body := data.Data.(*BodyAcceptDealAsStorageProvider)
		body.StorageProviderAddress = strings.ToUpper(body.StorageProviderAddress)
		proposal := models.BigFileProposal{}
		result := c.DB.Preload("BFile").Where("id = ?", body.ProposalID).First(&proposal)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid proposal Id", result.Error)
		}
		// TODO
		if proposal.BFile == nil {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid BigFile Status", result.Error)
		}
		if proposal.BFile.Status != 2 {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Big File is not in 'SP Selected' status", errors.New("Big File is not in 'SP Selected' Status"))
		}
		result = c.DB.Model(&proposal).Clauses(clause.Returning{}).Where("id = ?", body.ProposalID).Update("deal_c_id", body.DealCID)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error updating proposal", result.Error)
		}
		if result.RowsAffected != 1 {
			return nil, fiberutils.NewRequestError(http.StatusBadRequest, "Invalid big file id", errors.New("Invalid big file id"))
		}
		bigFile := models.BigFile{}
		result = c.DB.Model(&bigFile).Where("id = ?", proposal.BFileID).Update("status", 3)
		if result.Error != nil {
			return nil, fiberutils.NewRequestError(http.StatusInternalServerError, "Error updating Big File Status", result.Error)
		}
		return proposal, nil

	}
}
