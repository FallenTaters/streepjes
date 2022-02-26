package catalog

import (
	"time"
)

const cacheTime = time.Minute

var (
	cacheValid  = false
	lastCatalog struct {
		Catalog
		Time time.Time
	}
)

func Get() (Catalog, error) {
	if time.Since(lastCatalog.Time) > cacheTime || !cacheValid {
		c, err := getCatalog()
		if err != nil {
			return Catalog{}, err
		}
		lastCatalog.Catalog = c
		lastCatalog.Time = time.Now()

		cacheValid = true
	}

	return lastCatalog.Catalog, nil
}

func PutProduct(product Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}

	cacheValid = false

	if product.ID != 0 {
		return updateProduct(product)
	}
	return addProduct(product)
}

func DeleteProduct(id int) error {
	cacheValid = false
	return deleteProduct(id)
}

func PutCategory(category Category) error {
	if err := validateCategory(category); err != nil {
		return err
	}

	cacheValid = false

	if category.ID != 0 {
		return updateCategory(category)
	}
	return addCategory(category)
}

func DeleteCategory(id int) error {
	if err := validateDeleteCategory(id); err != nil {
		return err
	}

	cacheValid = false

	return deleteCategory(id)
}
