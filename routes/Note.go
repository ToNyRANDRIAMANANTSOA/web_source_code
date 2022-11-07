package routes

import (
	"back/database"
	"back/models"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Note struct {
	ID      uint `json:"id"`
	Student Student
	EC      EC
	Valeur  float64 `json:"note"`
}

var ErrNoteNotFound = errors.New("note not found")

func CreateNoteResponse(note *models.Note, student Student, ec EC) Note {
	return Note{
		ID:      note.ID,
		Student: student,
		EC:      ec,
		Valeur:  note.Note,
	}
}

func GetForeign(note *models.Note) (Student, EC, error) {
	student, err := ByMatricule(int(note.Matricule))
	if err != nil {
		return Student{}, EC{}, err
	}
	ec, err := ByCodeEc(note.CodeEC)
	if err != nil {
		return Student{}, EC{}, err
	}

	return CreateStudentResponse(student), CreateECResponse(ec), nil
}

func CreateNote(ctx *fiber.Ctx) error {
	var note models.Note
	err := ctx.BodyParser(&note)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	student, ec, err := GetForeign(&note)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}
	var eNote models.Note
	err = database.Database.DB.Where("matricule=?", student.Matricule).Where("code_ec=?", ec.CodeEC).First(&eNote).Error
	log.Println(err)
	if err == nil {
		eNote.Note = note.Note
		database.Database.DB.Save(&eNote)
		return ctx.Status(200).JSON(CreateNoteResponse(&eNote, student, ec))
	} else if err == gorm.ErrRecordNotFound {
		if err = database.Database.DB.Create(&note).Error; err != nil {
			return ctx.Status(400).SendString(err.Error())
		}
	}
	return ctx.Status(200).JSON(CreateNoteResponse(&note, student, ec))
}

func NoteByID(id int) (*models.Note, error) {
	var note models.Note
	err := database.Database.DB.Where("id=?", id).First(&note).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func DeleteNote(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(200).SendString(err.Error())
	}

	_, err = NoteByID(id)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	err = database.Database.DB.Where("id=?", id).Delete(&models.Note{}).Error
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).SendString("Note successfully deleted")
}

func GetNotesByMatricule(ctx *fiber.Ctx) error {
	matricule, err := ctx.ParamsInt("matricule")
	if err != nil {
		return ctx.Status(400).SendString("invalid matricule")
	}
	notes := []models.Note{}

	database.Database.DB.Where("matricule=?", matricule).Find(&notes)
	notesResponse := []Note{}
	for _, note := range notes {
		student, ec, err := GetForeign(&note)
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		notesResponse = append(notesResponse, CreateNoteResponse(&note, student, ec))
	}
	return ctx.Status(200).JSON(notesResponse)
}

func GetNoteByID(id uint) (*models.Note, error) {
	var note models.Note
	err := database.Database.DB.Where("id=?", id).First(&note).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func UpdateNote(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).SendString("invalid ID")
	}

	note, err := GetNoteByID(uint(id))
	if err != nil {
		return ctx.Status(400).SendString(err.Error())
	}

	type UpdateData struct {
		Note float64 `json:"note"`
	}
	var noteData UpdateData
	err = ctx.BodyParser(&noteData)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	note.Note = noteData.Note
	err = database.Database.DB.Save(note).Error
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	student, ec, _ := GetForeign(note)
	return ctx.Status(200).JSON(CreateNoteResponse(note, student, ec))
}
