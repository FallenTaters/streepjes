package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

func NewOrderRepo(db Queryable, logger *zap.Logger) repo.Order {
	return &orderRepo{
		db:     db,
		logger: logger,
	}
}

type orderRepo struct {
	db     Queryable
	logger *zap.Logger
}

func (or *orderRepo) Get(id int) (orderdomain.Order, bool) {
	row := or.db.QueryRow(`SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, `+
		`O.price, O.order_time, O.status, O.status_time FROM orders O WHERE id = $1;`, id)

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

	row := or.db.QueryRow(`SELECT * FROM users WHERE id = $1;`, order.BartenderID)
	if errors.Is(row.Scan(), sql.ErrNoRows) {
		return 0, fmt.Errorf("%w with id %d", repo.ErrUserNotFound, order.BartenderID)
	}

	if order.MemberID != 0 {
		row = or.db.QueryRow(`SELECT * FROM members WHERE id = $1;`, order.MemberID)
		if errors.Is(row.Scan(), sql.ErrNoRows) {
			return 0, fmt.Errorf("%w with id %d", repo.ErrMemberNotFound, order.MemberID)
		}
	}
	memberID := sql.NullInt64{
		Valid: order.MemberID != 0,
		Int64: int64(order.MemberID),
	}

	row = or.db.QueryRow(
		`INSERT INTO orders (club, bartender_id, member_id, contents, price, order_time, status, status_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
		order.Club, order.BartenderID, memberID, order.Contents, order.Price, order.OrderTime, order.Status, order.StatusTime,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	or.logger.Info("order created",
		zap.Int("id", id),
		zap.String("club", order.Club.String()),
		zap.Int("bartender_id", order.BartenderID),
		zap.Int("member_id", order.MemberID),
		zap.Int("price", int(order.Price)),
	)

	return id, nil
}

func (ur *orderRepo) Filter(filter repo.OrderFilter) []orderdomain.Order { //nolint:funlen,gocyclo,cyclop
	q := `SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, ` +
		`O.price, O.order_time, O.status, O.status_time FROM orders O `
	var conditions []string
	var args []interface{}

	if filter.BartenderID != 0 {
		args = append(args, filter.BartenderID)
		conditions = append(conditions, fmt.Sprintf(`O.bartender_id = $%d`, len(args)))
	}

	if filter.Club != domain.ClubUnknown {
		args = append(args, filter.Club)
		conditions = append(conditions, fmt.Sprintf(`O.club = $%d`, len(args)))
	}

	if len(filter.StatusNot) > 0 {
		for _, statusNot := range filter.StatusNot {
			args = append(args, statusNot)
			conditions = append(conditions, fmt.Sprintf(`O.status <> $%d`, len(args)))
		}
	}

	if filter.MemberID != 0 {
		args = append(args, filter.MemberID)
		conditions = append(conditions, fmt.Sprintf(`O.member_id = $%d`, len(args)))
	}

	if filter.Start != (time.Time{}) {
		args = append(args, filter.Start)
		conditions = append(conditions, fmt.Sprintf(`O.order_time >= $%d`, len(args)))
	}

	if filter.End != (time.Time{}) {
		args = append(args, filter.End)
		conditions = append(conditions, fmt.Sprintf(`O.order_time < $%d`, len(args)))
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
		q += fmt.Sprintf(` LIMIT $%d`, len(args))
		args = append(args, filter.Limit)
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

func (or *orderRepo) Delete(id int) error {
	result, err := or.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2;`, orderdomain.StatusCancelled, id)
	if err != nil {
		panic(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrOrderNotFound
	}

	or.logger.Info("order cancelled", zap.Int("id", id))

	return nil
}
