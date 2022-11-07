package models

type Student struct {
	Matricule uint   `json:"matricule" gorm:"primarykey"`
	Nom       string `json:"nom"`
	Adresse   string `json:"adresse"`
	Sexe      string `json:"sexe"`
	Niveau    string `json:"niveau"`
	Annee     int    `json:"annee"`
}
