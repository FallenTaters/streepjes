package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain"
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

func (or *orderRepo) Create(order orderdomain.Order) (int, error) {
	if order.BartenderID == 0 || order.Club == domain.ClubUnknown {
		return 0, fmt.Errorf("%w: %#v", repo.ErrOrderFieldsNotFilled, order)
	}

	row := or.db.QueryRow(`SELECT * FROM users WHERE id = ?;`, order.BartenderID)
	if errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, fmt.Errorf("%w with id %d\n", repo.ErrUserNotFound, order.BartenderID)
	}

	if order.MemberID != 0 {
		row = or.db.QueryRow(`SELECT * FROM members WHERE id = ?;`, order.MemberID)
		if errors.Is(row.Scan(), sql.ErrNoRows) {
			return 0, fmt.Errorf("%w with id %d\n", repo.ErrMemberNotFound, order.MemberID)
		}
	}

	res, err := or.db.Exec(
		`INSERT INTO orders (club, bartender_id, member_id, contents, price, order_time, status, status_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
		order.Club, order.BartenderID, order.MemberID, order.Contents, order.Price, order.OrderTime, order.Status, order.StatusTime,
	)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return int(id), nil
}

// func (ur *orderRepo) Get(id int) (orderdomain.Order, bool) {
// 	row := ur.db.QueryRow(
// `SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, O.price, O.order_time, O.status, O.status_time FROM orders O WHERE O.id = ?;`, id)

// 	var order orderdomain.Order

// 	err := row.Scan(
// &order.ID, &order.Club, &order.BartenderID, &order.MemberID, &order.Contents, &order.Price, &order.OrderTime, &order.Status, &order.StatusTime)
// 	if errors.Is(err, sql.ErrNoRows) {
// 		return orderdomain.Order{}, false
// 	}
// 	if err != nil {
// 		panic(err)
// 	}

// 	return order, true
// }

// func (ur *orderRepo) Filter(filter repo.OrderFilter) []orderdomain.Order {
// 	q := `SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, O.price, O.order_time, O.status, O.status_time FROM orders O `
// 	var conditions []string
// 	var args []interface{}

// 	if filter.BartenderID != nil {
// 		conditions = append(conditions, `O.bartender_id = ?`)
// 		args = append(args, *filter.BartenderID)
// 	}

// 	if filter.Club != nil {
// 		conditions = append(conditions, `O.club = ?`)
// 		args = append(args, *filter.Club)
// 	}

// 	if filter.MemberID != nil {
// 		conditions = append(conditions, `O.member_id = ?`)
// 		args = append(args, *filter.MemberID)
// 	}

// 	if filter.Month != nil {
// 		conditions = append(conditions, )
// 	}
// }
