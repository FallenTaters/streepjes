package mockdb

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Catalog struct {
	GetCategoriesFunc  func() ([]orderdomain.Category, error)
	GetItemsFunc       func() ([]orderdomain.Item, error)
	CreateItemFunc     func(item orderdomain.Item) (int, error)
	UpdateItemFunc     func(item orderdomain.Item) error
	DeleteItemFunc     func(id int) error
	CreateCategoryFunc func(cat orderdomain.Category) (int, error)
	UpdateCategoryFunc func(cat orderdomain.Category) error
	DeleteCategoryFunc func(id int) error
}

func (c Catalog) GetCategories() ([]orderdomain.Category, error) {
	return c.GetCategoriesFunc()
}

func (c Catalog) GetItems() ([]orderdomain.Item, error) {
	return c.GetItemsFunc()
}

func (c Catalog) CreateItem(item orderdomain.Item) (int, error) {
	return c.CreateItemFunc(item)
}

func (c Catalog) UpdateItem(item orderdomain.Item) error {
	return c.UpdateItemFunc(item)
}

func (c Catalog) DeleteItem(id int) error {
	return c.DeleteItemFunc(id)
}

func (c Catalog) CreateCategory(cat orderdomain.Category) (int, error) {
	return c.CreateCategoryFunc(cat)
}

func (c Catalog) UpdateCategory(cat orderdomain.Category) error {
	return c.UpdateCategoryFunc(cat)
}

func (c Catalog) DeleteCategory(id int) error {
	return c.DeleteCategoryFunc(id)
}
