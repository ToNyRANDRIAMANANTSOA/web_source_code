package routes

import (
	"back/database"
	"back/models"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	ErrStudentNotFound = errors.New("student : no student found")
)

type Student struct {
	Matricule uint   `json:"matricule"`
	Nom       string `json:"nom"`
	Adresse   string `json:"adresse"`
	Sexe      string `json:"sexe"`
	Niveau    string `json:"niveau"`
	Annee     int    `json:"annee"`
}

func CreateStudentResponse(s *models.Student) Student {
	return Student{
		Matricule: s.Matricule,
		Nom:       s.Nom,
		Adresse:   s.Adresse,
		Sexe:      s.Sexe,
		Niveau:    s.Niveau,
		Annee:     s.Annee,
	}
}

func CreateStudent(ctx *fiber.Ctx) error {
	var student models.Student
	if err := ctx.BodyParser(&student); err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	if err := database.Database.DB.Create(&student).Error; err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	return ctx.Status(200).JSON(CreateStudentResponse(&student))
}

func ByMatricule(id int) (*models.Student, error) {
	var student models.Student
	err := database.Database.DB.Where("matricule=?", id).First(&student).Error
	switch err {
	case nil:
		return &student, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrStudentNotFound
	default:
		return nil, err
	}
}

func ByNiveau(ctx *fiber.Ctx) error {
	students := []models.Student{}
	niveau := ctx.Params("niveau")
	if niveau == "" {
		return ctx.Status(400).SendString("The niveau can't be empty")
	}

	database.Database.DB.Where("niveau = ?", niveau).Find(&students)
	studentsResponse := []Student{}
	for _, s := range students {
		studentsResponse = append(studentsResponse, CreateStudentResponse(&s))
	}
	return ctx.Status(200).JSON(studentsResponse)
}

func GetStudents(ctx *fiber.Ctx) error {
	students := []models.Student{}
	database.Database.DB.Find(&students)
	studentsResponse := []Student{}
	for _, s := range students {
		studentsResponse = append(studentsResponse, CreateStudentResponse(&s))
	}
	return ctx.Status(200).JSON(studentsResponse)
}

func GetStudent(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("matricule")
	if err != nil {
		return ctx.Status(400).SendString("The matricule must be an integer")
	}

	student, err := ByMatricule(id)
	if err != nil {
		if err == ErrStudentNotFound {
			return ctx.Status(404).SendString("The student you're searching for doesn't exist")
		} else {
			return ctx.Status(400).SendString(err.Error())
		}
	}
	return ctx.Status(200).JSON(CreateStudentResponse(student))
}

func DeleteStudent(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("matricule")
	if err != nil {
		return ctx.Status(400).SendString("The matricule must be an integer")
	}

	student, err := ByMatricule(id)
	if err != nil {
		if err == ErrStudentNotFound {
			return ctx.Status(404).SendString("The student you're trying to delete for doesn't exist")
		} else {
			return ctx.Status(400).SendString(err.Error())
		}
	}
	err = database.Database.DB.Delete(&student).Error
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	database.Database.DB.Where("matricule=?", id).Delete(&models.Note{})
	return ctx.Status(200).SendString("Student Deleted successfully")
}

func UpdateStudent(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("matricule")
	if err != nil {
		return ctx.Status(400).SendString("The matricule must be an integer")
	}

	std, err := ByMatricule(id)
	if err != nil {
		if err == ErrStudentNotFound {
			return ctx.Status(404).SendString("The student you're trying to update doesn't exist")
		} else {
			return ctx.Status(400).SendString(err.Error())
		}
	}
	type StudentUpdate struct {
		Nom     string `json:"nom"`
		Adresse string `json:"adresse"`
		Niveau  string `json:"niveau"`
		Annee   int    `json:"annee"`
	}
	updateData := StudentUpdate{}
	fmt.Println(string(ctx.Body()))
	err = ctx.BodyParser(&updateData)
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	if updateData.Niveau == "" {
		updateData.Niveau = std.Niveau
	}
	student := &models.Student{
		Matricule: uint(id),
		Nom:       updateData.Nom,
		Adresse:   updateData.Adresse,
		Sexe:      std.Sexe,
		Niveau:    updateData.Niveau,
		Annee:     updateData.Annee,
	}

	err = database.Database.DB.Save(student).Error
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}
	return ctx.Status(200).JSON(CreateStudentResponse(student))
}
