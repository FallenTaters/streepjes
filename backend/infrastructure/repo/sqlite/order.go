package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (or *orderRepo) Get(id int) (orderdomain.Order, bool) {
	row := or.db.QueryRow(`SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, `+
		`O.price, O.order_time, O.status, O.status_time FROM orders O WHERE id = ?;`, id)

	var order orderdomain.Order
	var memberID sql.NullInt64

	err := row.Scan(&order.ID, &order.Club, &order.BartenderID, &memberID, &order.Contents,
		&order.Price, &order.OrderTime, &order.Status, &order.StatusTime)
	if errors.Is(err, sql.ErrNoRows) {
		return order, false
	}
	if err != nil {
		panic(err)
	}

	order.MemberID = int(memberID.Int64)

	return order, true
}

func (or *orderRepo) Create(order orderdomain.Order) (int, error) {
	if order.BartenderID == 0 || order.Club == domain.ClubUnknown {
		return 0, fmt.Errorf("%w: %#v", repo.ErrOrderFieldsNotFilled, order)
	}

	row := or.db.QueryRow(`SELECT * FROM users WHERE id = ?;`, order.BartenderID)
	if errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, fmt.Errorf("%w with id %d", repo.ErrUserNotFound, order.BartenderID)
	}

	if order.MemberID != 0 {
		row = or.db.QueryRow(`SELECT * FROM members WHERE id = ?;`, order.MemberID)
		if errors.Is(row.Scan(), sql.ErrNoRows) {
			return 0, fmt.Errorf("%w with id %d", repo.ErrMemberNotFound, order.MemberID)
		}
	}
	memberID := sql.NullInt64{
		Valid: order.MemberID != 0,
		Int64: int64(order.MemberID),
	}

	res, err := or.db.Exec(
		`INSERT INTO orders (club, bartender_id, member_id, contents, price, order_time, status, status_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
		order.Club, order.BartenderID, memberID, order.Contents, order.Price, order.OrderTime, order.Status, order.StatusTime,
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

func (ur *orderRepo) Filter(filter repo.OrderFilter) []orderdomain.Order { //nolint:funlen,gocyclo,cyclop
	q := `SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, ` +
		`O.price, O.order_time, O.status, O.status_time FROM orders O `
	var conditions []string
	var args []interface{}

	if filter.BartenderID != 0 {
		conditions = append(conditions, `O.bartender_id = ?`)
		args = append(args, filter.BartenderID)
	}

	if filter.Club != domain.ClubUnknown {
		conditions = append(conditions, `O.club = ?`)
		args = append(args, filter.Club)
	}

	if len(filter.StatusNot) > 0 {
		conditions = append(conditions, `O.status NOT IN (?`+strings.Repeat(`,?`, len(filter.StatusNot)-1)+`)`)
		for _, statusNot := range filter.StatusNot {
			args = append(args, statusNot)
		}
	}

	if filter.MemberID != 0 {
		conditions = append(conditions, `O.member_id = ?`)
		args = append(args, filter.MemberID)
	}

	if filter.Start != (time.Time{}) {
		conditions = append(conditions, `O.order_time >= ?`)
		args = append(args, filter.Start)
	}

	if filter.End != (time.Time{}) {
		conditions = append(conditions, `O.order_time < ?`)
		args = append(args, filter.End)
	}

	if len(conditions) > 0 {
		q += `WHERE `
	}
	for i, condition := range conditions {
		q += condition
		if i < len(conditions)-1 {
			q += ` AND `
		}
	}
	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		q += ` LIMIT ?`
	}
	q += `;`

	rows, err := ur.db.Query(q, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var orders []orderdomain.Order
	for rows.Next() {
		var order orderdomain.Order
		var memberID sql.NullInt64

		err := rows.Scan(
			&order.ID, &order.Club, &order.BartenderID, &memberID, &order.Contents,
			&order.Price, &order.OrderTime, &order.Status, &order.StatusTime)
		if err != nil {
			panic(err)
		}

		order.MemberID = int(memberID.Int64)
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return orders
}

func (or *orderRepo) Delete(id int) bool {
	result, err := or.db.Exec(`UPDATE orders SET status = ? WHERE id = ?;`, orderdomain.StatusCancelled, id)
	if err != nil {
		panic(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return affected == 1
}
