package repo

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

var (
	ErrCategoryNameTaken = errors.New(`category name taken`)
	ErrCategoryNotFound  = errors.New(`category not found`)

	ErrItemNameTaken = errors.New(`item name taken`)
	ErrItemNotFound  = errors.New(`item not found`)
)

type Catalog interface {
	// Get entire catalog
	Get() (api.Catalog, error)

	// GetItem gets a single item. If not found, it returns false
	GetItem(id int) (orderdomain.Item, bool)

	// CreateItem makes a new item. If the name is already taken, it returns ErrItemNameTaken
	CreateItem(orderdomain.Item) error

	// UpdateItem updates a item. If not found, it returns ErrItemNotFound
	UpdateItem(orderdomain.Item) error

	// DeleteItem deletes a item by id. If not found, it returns ErrItemNotFound
	DeleteItem(id int) error

	// GetCategory gets a single category. If not found, it returns false
	GetCategory(id int) (orderdomain.Category, bool)

	// CreateCategory makes a new category. If the name is already taken, it returns ErrCategoryNameTaken
	CreateCategory(orderdomain.Category) error

	// UpdateCategory updates a category. If not found, it returns ErrCategoryNotFound
	UpdateCategory(orderdomain.Category) error

	// DeleteCategory deletes a category by id. If not found, it returns ErrCategoryNotFound
	DeleteCategory(id int) error
}
