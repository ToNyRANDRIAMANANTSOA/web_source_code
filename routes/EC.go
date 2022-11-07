package routes

import (
	"back/database"
	"back/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	ErrECNotFound = errors.New("ec not found")
)

type EC struct {
	CodeEC      int    `json:"codeEC"`
	Libelle     string `json:"libelle"`
	Coefficient int    `json:"coefficient"`
}

func CreateECResponse(ec *models.EC) EC {
	return EC{
		CodeEC:      ec.CodeEC,
		Libelle:     ec.Libelle,
		Coefficient: ec.Coefficient,
	}
}

func CreateEC(ctx *fiber.Ctx) error {
	var ec models.EC
	if err := ctx.BodyParser(&ec); err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	if err := database.Database.DB.Create(&ec).Error; err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	return ctx.Status(200).JSON(CreateECResponse(&ec))
}

func ByCodeEc(codeEC int) (*models.EC, error) {
	var ec models.EC
	err := database.Database.DB.Where("code_ec=?", codeEC).First(&ec).Error
	switch err {
	case nil:
		return &ec, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrECNotFound
	default:
		return nil, err
	}
}

func GetECs(ctx *fiber.Ctx) error {
	ecs := []models.EC{}

	database.Database.DB.Find(&ecs)
	ecsResponse := []EC{}
	for _, ec := range ecs {
		ecsResponse = append(ecsResponse, CreateECResponse(&ec))
	}
	return ctx.Status(200).JSON(ecsResponse)
}

func GetEC(ctx *fiber.Ctx) error {
	codeEC, err := ctx.ParamsInt("codeEC")
	if err != nil {
		return ctx.Status(400).SendString("The codeEC must be integer")
	}

	ec, err := ByCodeEc(codeEC)
	if err != nil {
		if err == ErrECNotFound {
			return ctx.Status(404).SendString(err.Error())
		} else {
			return ctx.Status(500).SendString(err.Error())
		}
	}
	return ctx.Status(200).JSON(CreateECResponse(ec))
}

func UpdateEC(ctx *fiber.Ctx) error {
	codeEC, err := ctx.ParamsInt("codeEC")
	if err != nil {
		return ctx.Status(400).SendString("The codeEC must be integer")
	}

	ec, err := ByCodeEc(codeEC)
	if err != nil {
		if err == ErrECNotFound {
			return ctx.Status(404).SendString(err.Error())
		} else {
			return ctx.Status(500).SendString(err.Error())
		}
	}

	type ECUpdate struct {
		Libelle     string `json:"libelle"`
		Coefficient int    `json:"coefficient"`
	}
	var updateData ECUpdate
	err = ctx.BodyParser(&updateData)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	ec.Libelle = updateData.Libelle
	ec.Coefficient = updateData.Coefficient
	err = database.Database.DB.Save(&ec).Error
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	return ctx.Status(200).JSON(CreateECResponse(ec))
}
