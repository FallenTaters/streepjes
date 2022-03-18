package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func NewOrderRepo(db Queryable) repo.Order {
	return &orderRepo{
		db: db,
	}
}

type orderRepo struct {
	db Queryable
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

func (ur *orderRepo) Filter(filter repo.OrderFilter) []orderdomain.Order { //nolint:funlen
	q := `SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, ` +
		`O.price, O.order_time, O.status, O.status_time FROM orders O `
	var conditions []string
	var args []interface{}

	// if filter.BartenderID != nil {
	// 	conditions = append(conditions, `O.bartender_id = ?`)
	// 	args = append(args, *filter.BartenderID)
	// }

	// if filter.Club != nil {
	// 	conditions = append(conditions, `O.club = ?`)
	// 	args = append(args, *filter.Club)
	// }

	if filter.MemberID != 0 {
		conditions = append(conditions, `O.member_id = ?`)
		args = append(args, filter.MemberID)
	}

	if filter.Month != (orderdomain.Month{}) {
		conditions = append(conditions, `O.order_time >= ?`)
		args = append(args, filter.Month.Time())

		conditions = append(conditions, `O.order_time < ?`)
		args = append(args, filter.Month.Time().AddDate(0, 1, 0))
	}

	if len(conditions) > 0 {
		q += `WHERE `
	}

	for _, condition := range conditions {
		q += condition + ` `
	}

	rows, err := ur.db.Query(q, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var orders []orderdomain.Order
	for rows.Next() {
		var order orderdomain.Order

		err := rows.Scan(
			&order.ID, &order.Club, &order.BartenderID, &order.MemberID, &order.Contents,
			&order.Price, &order.OrderTime, &order.Status, &order.StatusTime)
		if err != nil {
			panic(err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return orders
}
