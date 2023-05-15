package postgres

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
	defer rows.Close()

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
	defer rows.Close()

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

	row := cr.db.QueryRow(`SELECT * FROM items WHERE name = $1;`, item.Name)
	if !errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, repo.ErrItemNameTaken
	}

	row = cr.db.QueryRow(`INSERT INTO items (category_id, name, price_gladiators, price_parabool) VALUES ($1, $2, $3, $4) RETURNING id;`,
		item.CategoryID, item.Name, item.PriceGladiators, item.PriceParabool)

	var id int
	return id, row.Scan(&id)
}

func (cr catalogRepo) CreateCategory(category orderdomain.Category) (int, error) {
	if category.Name == `` {
		return 0, repo.ErrCategoryNameEmpty
	}

	row := cr.db.QueryRow(`SELECT * FROM categories WHERE name = $1;`, category.Name)
	if !errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, repo.ErrCategoryNameTaken
	}

	row = cr.db.QueryRow(`INSERT INTO categories (name) VALUES ($1) RETURNING id;`, category.Name)

	var id int
	return id, row.Scan(&id)
}

func (cr catalogRepo) getCategory(id int) (orderdomain.Category, bool) {
	row := cr.db.QueryRow(`SELECT id, name FROM categories WHERE id = $1;`, id)

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
	res, err := cr.db.Exec(`UPDATE categories SET name = $1 WHERE id = $2;`, cat.Name, cat.ID)
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
	rows, err := cr.db.Query(`SELECT * from items WHERE category_id = $1;`, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return repo.ErrCategoryHasItems
	}

	// exec delete
	res, err := cr.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
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

	res, err := cr.db.Exec(`UPDATE items SET category_id = $1, name = $2, price_parabool = $3, price_gladiators = $4 WHERE id = $5;`,
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
	res, err := cr.db.Exec(`DELETE FROM items WHERE id = $1`, id)
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
