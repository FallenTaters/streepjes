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

func (cr catalogRepo) getCategory(id int) (orderdomain.Category, bool) {
	row := cr.db.QueryRow(`SELECT id, name FROM categories WHERE id = ?;`, id)

	var cat orderdomain.Category
	if err := row.Scan(&cat.ID, &cat.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cat, false
		}

		panic(err)
	}

	return cat, true
}

func (cr catalogRepo) UpdateCategory(cat orderdomain.Category) error {
	res, err := cr.db.Exec(`UPDATE categories SET name = ? WHERE id = ?;`, cat.Name, cat.ID)
	if err != nil {
		return repo.ErrCategoryNameTaken
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrCategoryNotFound
	}

	return nil
}

func (cr catalogRepo) DeleteCategory(id int) error {
	// check for existing items
	rows, err := cr.db.Query(`SELECT * from items WHERE category_id = ?;`, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		return repo.ErrCategoryHasItems
	}

	// exec delete
	res, err := cr.db.Exec(`DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		panic(err)
	}

	// check if actually found
	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		return repo.ErrCategoryNotFound
	}

	return nil
}

func (cr catalogRepo) UpdateItem(item orderdomain.Item) error {
	if _, exists := cr.getCategory(item.CategoryID); !exists {
		return repo.ErrCategoryNotFound
	}

	res, err := cr.db.Exec(`UPDATE items SET category_id = ?, name = ?, price_parabool = ?, price_gladiators = ? WHERE id = ?;`,
		item.CategoryID, item.Name, item.PriceParabool, item.PriceGladiators, item.ID)
	if err != nil {
		return repo.ErrItemNameTaken
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrItemNotFound
	}

	return nil
}

func (cr catalogRepo) DeleteItem(id int) error {
	res, err := cr.db.Exec(`DELETE FROM items WHERE id = ?`, id)
	if err != nil {
		panic(err)
	}

	// check if actually found
	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if affected == 0 {
		return repo.ErrItemNotFound
	}

	return nil
}
