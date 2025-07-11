package model

import (
	"errors"
	"strings"
	"todolist/common"
)

var (
	ErrTitleCannotBeEmpty = errors.New("Title can not be empty")
	ErrItemIsDeleted      = errors.New("This Item is deleted")
)

type TodoItem struct {
	common.SQLModel
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Id          int    `json:"id" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	if i.Title == "" {
		return ErrTitleCannotBeEmpty
	}
	return nil
}
func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"` // Chuyen thanh con tro thi moi update truong nay voi gia tri rong duoc
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

//

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()

}
