package repo

import (
	"errors"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrCategoryNameTaken = errors.New(`category name taken`)
	ErrCategoryNameEmpty = errors.New(`category name cannot be empty`)
	ErrCategoryNotFound  = errors.New(`category not found`)
	ErrCategoryHasItems  = errors.New(`category has items`)

	ErrItemNameTaken = errors.New(`item name taken`)
	ErrItemNameEmpty = errors.New(`item name cannot be empty`)
	ErrItemNotFound  = errors.New(`item not found`)
)

type Catalog interface {
	GetCategories() ([]orderdomain.Category, error)
	GetItems() ([]orderdomain.Item, error)
	CreateItem(orderdomain.Item) (int, error)
	UpdateItem(orderdomain.Item) error
	DeleteItem(id int) error
	CreateCategory(orderdomain.Category) (int, error)
	UpdateCategory(orderdomain.Category) error
	DeleteCategory(id int) error
}
