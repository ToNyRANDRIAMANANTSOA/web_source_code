package models

type Note struct {
	ID        uint    `gorm:"primarykey"`
	CodeEC    int     `json:"codeEC"`
	EC        EC      `gorm:"foreignKey:CodeEC"`
	Matricule uint    `json:"matricule"`
	Student   Student `gorm:"foreignKey:Matricule"`
	Note      float64 `json:"note"`
}
