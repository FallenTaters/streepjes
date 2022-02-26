package catalog

import (
	"errors"

	"github.com/FallenTaters/bbucket"
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrNameTaken          = errors.New("this name is taken")
	ErrEmptyName          = errors.New("name may not be empty")
	ErrNoPrice            = errors.New("product must have at least one price")
	ErrCategoryHasProduct = errors.New("category still has products")
	ErrProductNotFound    = errors.New("product not found")
)

func validateProduct(product Product) error {
	original, err := getProduct(product.ID)
	switch err {
	case nil:
		if original.Name != product.Name && productNameExists(product.Name) {
			return ErrNameTaken
		}
	case bbucket.ErrObjectNotFound:
		if product.ID != 0 {
			return ErrProductNotFound
		}
		if productNameExists(product.Name) {
			return ErrNameTaken
		}
	default:
		panic(err)
	}

	if !categoryExists(product.CategoryID) {
		return ErrCategoryNotFound
	}
	if product.Name == `` {
		return ErrEmptyName
	}
	if product.PriceGladiators == 0 && product.PriceParabool == 0 {
		return ErrNoPrice
	}

	return nil
}

func productNameExists(name string) bool {
	products, err := getProducts()
	if err != nil {
		panic(err)
	}

	for _, p := range products {
		if p.Name == name {
			return true
		}
	}

	return false
}

func categoryExists(id int) bool {
	categories, err := getCategories()
	if err != nil {
		panic(err)
	}

	for _, c := range categories {
		if c.ID == id {
			return true
		}
	}

	return false
}

func validateCategory(category Category) error {
	original, err := getCategory(category.ID)
	switch err {
	case nil:
		if original.Name != category.Name && categoryNameExists(category.Name) {
			return ErrNameTaken
		}
	case bbucket.ErrObjectNotFound:
		if category.ID != 0 {
			return ErrCategoryNotFound
		}
		if categoryNameExists(category.Name) {
			return ErrNameTaken
		}
	default:
		panic(err)
	}

	if category.Name == `` {
		return ErrEmptyName
	}

	return nil
}

func categoryNameExists(name string) bool {
	categories, err := getCategories()
	if err != nil {
		panic(err)
	}

	for _, c := range categories {
		if c.Name == name {
			return true
		}
	}

	return false
}

func validateDeleteCategory(id int) error {
	products, err := getProducts()
	if err != nil {
		panic(err)
	}

	for _, p := range products {
		if p.CategoryID == id {
			return ErrCategoryHasProduct
		}
	}

	return nil
}
