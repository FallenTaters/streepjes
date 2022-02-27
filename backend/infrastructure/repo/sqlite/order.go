package sqlite

import (
	"database/sql"
	"errors"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

func NewOrderRepo(db *sql.DB) repo.Order {
	return &orderRepo{
		db: db,
	}
}

type orderRepo struct {
	db *sql.DB
}

func (ur *orderRepo) Get(id int) (orderdomain.Order, bool) {
	row := ur.db.QueryRow(`SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, O.price, O.order_time, O.status, O.status_time FROM orders O WHERE O.id = ?;`, id)

	var order orderdomain.Order

	err := row.Scan(&order.ID, &order.Club, &order.BartenderID, &order.MemberID, &order.Contents, &order.Price, &order.OrderTime, &order.Status, &order.StatusTime)
	if errors.Is(err, sql.ErrNoRows) {
		return orderdomain.Order{}, false
	}
	if err != nil {
		panic(err)
	}

	return order, true
}

func (ur *orderRepo) Filter(filter repo.OrderFilter) []orderdomain.Order {
	q := `SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, O.price, O.order_time, O.status, O.status_time FROM orders O `
	var conditions []string
	var args []interface{}

	if filter.BartenderID != nil {
		conditions = append(conditions, `O.bartender_id = ?`)
		args = append(args, *filter.BartenderID)
	}

	if filter.Club != nil {
		conditions = append(conditions, `O.club = ?`)
		args = append(args, *filter.Club)
	}

	if filter.MemberID != nil {
		conditions = append(conditions, `O.member_id = ?`)
		args = append(args, *filter.MemberID)
	}
}
