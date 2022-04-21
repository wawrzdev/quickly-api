package data

import "gorm.io/gorm"

type Notebook struct {
	gorm.Model
	Name  string
	User  User
	Notes []Note
}
