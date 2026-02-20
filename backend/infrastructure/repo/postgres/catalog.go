package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

func NewCatalogRepo(db Queryable, logger *zap.Logger) repo.Catalog {
	return &catalogRepo{db: db, logger: logger}
}

type catalogRepo struct {
	db     Queryable
	logger *zap.Logger
}

func (cr catalogRepo) GetCategories() ([]orderdomain.Category, error) {
	rows, err := cr.db.Query(`SELECT id, name FROM categories;`)
	if err != nil {
		return nil, fmt.Errorf("catalogRepo.GetCategories: query: %w", err)
	}
	defer rows.Close()

	var categories []orderdomain.Category
	for rows.Next() {
		var c orderdomain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("catalogRepo.GetCategories: scan: %w", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("catalogRepo.GetCategories: rows: %w", err)
	}

	return categories, nil
}

func (cr catalogRepo) GetItems() ([]orderdomain.Item, error) {
	rows, err := cr.db.Query(`SELECT id, category_id, name, price_gladiators, price_parabool, price_calamari FROM items;`)
	if err != nil {
		return nil, fmt.Errorf("catalogRepo.GetItems: query: %w", err)
	}
	defer rows.Close()

	var items []orderdomain.Item
	for rows.Next() {
		var item orderdomain.Item
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.PriceGladiators, &item.PriceParabool, &item.PriceCalamari); err != nil {
			return nil, fmt.Errorf("catalogRepo.GetItems: scan: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("catalogRepo.GetItems: rows: %w", err)
	}

	return items, nil
}

func (cr catalogRepo) CreateItem(item orderdomain.Item) (int, error) {
	if item.Name == `` {
		return 0, repo.ErrItemNameEmpty
	}

	var exists bool
	if err := cr.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM items WHERE name = $1)`, item.Name).Scan(&exists); err != nil {
		return 0, fmt.Errorf("catalogRepo.CreateItem: check name: %w", err)
	}
	if exists {
		return 0, repo.ErrItemNameTaken
	}

	var id int
	if err := cr.db.QueryRow(`INSERT INTO items (category_id, name, price_gladiators, price_parabool, price_calamari) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		item.CategoryID, item.Name, item.PriceGladiators, item.PriceParabool, item.PriceCalamari,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("catalogRepo.CreateItem: insert: %w", err)
	}

	cr.logger.Info("item created", zap.Int("id", id), zap.String("name", item.Name), zap.Int("category_id", item.CategoryID))
	return id, nil
}

func (cr catalogRepo) CreateCategory(category orderdomain.Category) (int, error) {
	if category.Name == `` {
		return 0, repo.ErrCategoryNameEmpty
	}

	var exists bool
	if err := cr.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1)`, category.Name).Scan(&exists); err != nil {
		return 0, fmt.Errorf("catalogRepo.CreateCategory: check name: %w", err)
	}
	if exists {
		return 0, repo.ErrCategoryNameTaken
	}

	var id int
	if err := cr.db.QueryRow(`INSERT INTO categories (name) VALUES ($1) RETURNING id;`, category.Name,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("catalogRepo.CreateCategory: insert: %w", err)
	}

	cr.logger.Info("category created", zap.Int("id", id), zap.String("name", category.Name))
	return id, nil
}

func (cr catalogRepo) getCategory(id int) (orderdomain.Category, error) {
	row := cr.db.QueryRow(`SELECT id, name FROM categories WHERE id = $1;`, id)

	var cat orderdomain.Category
	if err := row.Scan(&cat.ID, &cat.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cat, repo.ErrCategoryNotFound
		}
		return cat, fmt.Errorf("catalogRepo.getCategory(%d): %w", id, err)
	}

	return cat, nil
}

func (cr catalogRepo) UpdateCategory(cat orderdomain.Category) error {
	res, err := cr.db.Exec(`UPDATE categories SET name = $1 WHERE id = $2;`, cat.Name, cat.ID)
	if err != nil {
		return fmt.Errorf("catalogRepo.UpdateCategory(%d): %w", cat.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("catalogRepo.UpdateCategory(%d): rows affected: %w", cat.ID, err)
	}
	if affected == 0 {
		return repo.ErrCategoryNotFound
	}

	cr.logger.Info("category updated", zap.Int("id", cat.ID), zap.String("name", cat.Name))
	return nil
}

func (cr catalogRepo) DeleteCategory(id int) error {
	var hasItems bool
	if err := cr.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM items WHERE category_id = $1)`, id).Scan(&hasItems); err != nil {
		return fmt.Errorf("catalogRepo.DeleteCategory(%d): check items: %w", id, err)
	}
	if hasItems {
		return repo.ErrCategoryHasItems
	}

	res, err := cr.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("catalogRepo.DeleteCategory(%d): %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("catalogRepo.DeleteCategory(%d): rows affected: %w", id, err)
	}
	if affected == 0 {
		return repo.ErrCategoryNotFound
	}

	cr.logger.Info("category deleted", zap.Int("id", id))
	return nil
}

func (cr catalogRepo) UpdateItem(item orderdomain.Item) error {
	if _, err := cr.getCategory(item.CategoryID); err != nil {
		if errors.Is(err, repo.ErrCategoryNotFound) {
			return repo.ErrCategoryNotFound
		}
		return fmt.Errorf("catalogRepo.UpdateItem: check category: %w", err)
	}

	res, err := cr.db.Exec(`UPDATE items SET category_id = $1, name = $2, price_parabool = $3, price_gladiators = $4, price_calamari = $5 WHERE id = $6;`,
		item.CategoryID, item.Name, item.PriceParabool, item.PriceGladiators, item.PriceCalamari, item.ID)
	if err != nil {
		return fmt.Errorf("catalogRepo.UpdateItem(%d): %w", item.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("catalogRepo.UpdateItem(%d): rows affected: %w", item.ID, err)
	}
	if affected == 0 {
		return repo.ErrItemNotFound
	}

	cr.logger.Info("item updated", zap.Int("id", item.ID), zap.String("name", item.Name))
	return nil
}

func (cr catalogRepo) DeleteItem(id int) error {
	res, err := cr.db.Exec(`DELETE FROM items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("catalogRepo.DeleteItem(%d): %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("catalogRepo.DeleteItem(%d): rows affected: %w", id, err)
	}
	if affected == 0 {
		return repo.ErrItemNotFound
	}

	cr.logger.Info("item deleted", zap.Int("id", id))
	return nil
}
