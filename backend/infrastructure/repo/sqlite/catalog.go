package sqlite

import (
	"database/sql"
	"errors"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func NewCatalogRepo(db Queryable) repo.Catalog {
	return &catalogRepo{
		db: db,
	}
}

type catalogRepo struct {
	db Queryable
}

func (cr catalogRepo) GetCategories() []orderdomain.Category {
	rows, err := cr.db.Query(`SELECT id, name FROM categories;`)
	if err != nil {
		panic(err)
	}

	var categories []orderdomain.Category
	for rows.Next() {
		var category orderdomain.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			panic(err)
		}

		categories = append(categories, category)
	}

	return categories
}

func (cr catalogRepo) GetItems() []orderdomain.Item {
	rows, err := cr.db.Query(`SELECT id, category_id, name, price_gladiators, price_parabool FROM items;`)
	if err != nil {
		panic(err)
	}

	var items []orderdomain.Item
	for rows.Next() {
		var item orderdomain.Item
		err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.PriceGladiators, &item.PriceParabool)
		if err != nil {
			panic(err)
		}

		items = append(items, item)
	}

	return items
}

func (cr catalogRepo) CreateItem(item orderdomain.Item) (int, error) {
	if item.Name == `` {
		return 0, repo.ErrItemNameEmpty
	}

	row := cr.db.QueryRow(`SELECT * FROM items WHERE name = ?;`, item.Name)
	if !errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, repo.ErrItemNameTaken
	}

	res, err := cr.db.Exec(`INSERT INTO items (category_id, name, price_gladiators, price_parabool) VALUES (?, ?, ?, ?);`,
		item.CategoryID, item.Name, item.PriceGladiators, item.PriceParabool)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(id), nil
}

func (cr catalogRepo) CreateCategory(category orderdomain.Category) (int, error) {
	if category.Name == `` {
		return 0, repo.ErrCategoryNameEmpty
	}

	row := cr.db.QueryRow(`SELECT * FROM categories WHERE name = ?;`, category.Name)
	if !errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, repo.ErrCategoryNameTaken
	}

	res, err := cr.db.Exec(`INSERT INTO categories (name) VALUES (?);`, category.Name)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(id), nil
}
