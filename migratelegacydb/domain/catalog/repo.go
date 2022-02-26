package catalog

import (
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared/buckets"
)

func getCatalog() (Catalog, error) {
	catalog := Catalog{}

	categories, err := getCategories()
	if err != nil {
		return catalog, err
	}
	catalog.Categories = categories

	products, err := getProducts()
	if err != nil {
		return catalog, err
	}
	catalog.Products = products

	return catalog, nil
}

func getCategories() ([]Category, error) {
	categories := []Category{}
	return categories, buckets.Categories.GetAll(&Category{}, func(ptr interface{}) error {
		categories = append(categories, *ptr.(*Category))
		return nil
	})
}

func getCategory(id int) (Category, error) {
	var category Category
	return category, buckets.Categories.Get(buckets.Itob(id), &category)
}

func addCategory(category Category) error {
	category.ID = buckets.Categories.NextSequence()
	return buckets.Categories.Create(category.Key(), category)
}

func updateCategory(category Category) error {
	return buckets.Categories.Update(category.Key(), &Category{}, func(ptr interface{}) (object interface{}, err error) {
		return category, nil
	})
}

func deleteCategory(id int) error {
	return buckets.Categories.Delete(buckets.Itob(id))
}

func getProducts() ([]Product, error) {
	products := []Product{}
	return products, buckets.Products.GetAll(&Product{}, func(ptr interface{}) error {
		products = append(products, *ptr.(*Product))
		return nil
	})
}

func getProduct(id int) (Product, error) {
	var product Product
	return product, buckets.Products.Get(buckets.Itob(id), &product)
}

func addProduct(product Product) error {
	product.ID = buckets.Products.NextSequence()
	return buckets.Products.Create(product.Key(), product)
}

func updateProduct(product Product) error {
	return buckets.Products.Update(product.Key(), &Product{}, func(ptr interface{}) (object interface{}, err error) {
		return product, nil
	})
}

func deleteProduct(id int) error {
	return buckets.Products.Delete(buckets.Itob(id))
}
