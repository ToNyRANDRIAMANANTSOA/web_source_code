package models

type EC struct {
	CodeEC      int    `gorm:"primarykey,not null" json:"codeEC"`
	Libelle     string `json:"libelle"`
	Coefficient int    `json:"coefficient"`
}
