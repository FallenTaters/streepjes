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
	return &orderRepo{db: db, logger: logger}
}

type orderRepo struct {
	db     Queryable
	logger *zap.Logger
}

func (or *orderRepo) Get(id int) (orderdomain.Order, error) {
	row := or.db.QueryRow(`SELECT O.id, O.club, O.bartender_id, O.member_id, O.contents, `+
		`O.price, O.order_time, O.status, O.status_time FROM orders O WHERE id = $1;`, id)

	var order orderdomain.Order
	var memberID sql.NullInt64

	err := row.Scan(&order.ID, &order.Club, &order.BartenderID, &memberID, &order.Contents,
		&order.Price, &order.OrderTime, &order.Status, &order.StatusTime)
	if errors.Is(err, sql.ErrNoRows) {
		return orderdomain.Order{}, repo.ErrOrderNotFound
	}
	if err != nil {
		return orderdomain.Order{}, fmt.Errorf("orderRepo.Get(%d): %w", id, err)
	}

	order.MemberID = int(memberID.Int64)
	return order, nil
}

func (or *orderRepo) Create(order orderdomain.Order) (int, error) {
	if order.BartenderID == 0 || order.Club == domain.ClubUnknown {
		return 0, fmt.Errorf("%w: %#v", repo.ErrOrderFieldsNotFilled, order)
	}

	var exists bool
	if err := or.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, order.BartenderID).Scan(&exists); err != nil {
		return 0, fmt.Errorf("orderRepo.Create: check bartender: %w", err)
	}
	if !exists {
		return 0, fmt.Errorf("%w with id %d", repo.ErrUserNotFound, order.BartenderID)
	}

	if order.MemberID != 0 {
		if err := or.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM members WHERE id = $1)`, order.MemberID).Scan(&exists); err != nil {
			return 0, fmt.Errorf("orderRepo.Create: check member: %w", err)
		}
		if !exists {
			return 0, fmt.Errorf("%w with id %d", repo.ErrMemberNotFound, order.MemberID)
		}
	}
	memberID := sql.NullInt64{
		Valid: order.MemberID != 0,
		Int64: int64(order.MemberID),
	}

	var id int
	if err := or.db.QueryRow(
		`INSERT INTO orders (club, bartender_id, member_id, contents, price, order_time, status, status_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
		order.Club, order.BartenderID, memberID, order.Contents, order.Price, order.OrderTime, order.Status, order.StatusTime,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("orderRepo.Create: insert: %w", err)
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

func (or *orderRepo) Filter(filter repo.OrderFilter) ([]orderdomain.Order, error) { //nolint:funlen,gocyclo,cyclop
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
		args = append(args, filter.Limit)
		q += fmt.Sprintf(` LIMIT $%d`, len(args))
	}
	q += `;`

	rows, err := or.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("orderRepo.Filter: query: %w", err)
	}
	defer rows.Close()

	var orders []orderdomain.Order
	for rows.Next() {
		var order orderdomain.Order
		var memberID sql.NullInt64

		if err := rows.Scan(
			&order.ID, &order.Club, &order.BartenderID, &memberID, &order.Contents,
			&order.Price, &order.OrderTime, &order.Status, &order.StatusTime); err != nil {
			return nil, fmt.Errorf("orderRepo.Filter: scan: %w", err)
		}

		order.MemberID = int(memberID.Int64)
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("orderRepo.Filter: rows: %w", err)
	}

	return orders, nil
}

func (or *orderRepo) Delete(id int) error {
	result, err := or.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2;`, orderdomain.StatusCancelled, id)
	if err != nil {
		return fmt.Errorf("orderRepo.Delete(%d): %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("orderRepo.Delete(%d): rows affected: %w", id, err)
	}
	if affected == 0 {
		return repo.ErrOrderNotFound
	}

	or.logger.Info("order cancelled", zap.Int("id", id))
	return nil
}
