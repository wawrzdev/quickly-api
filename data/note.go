package data

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Content    string
	NotebookID int
	Notebook   Notebook
}
