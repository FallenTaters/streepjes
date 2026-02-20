package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

func NewMemberRepo(db Queryable, logger *zap.Logger) repo.Member {
	return &memberRepo{db: db, logger: logger}
}

type memberRepo struct {
	db     Queryable
	logger *zap.Logger
}

func (mr *memberRepo) GetAll() ([]orderdomain.Member, error) {
	rows, err := mr.db.Query(`SELECT id, club, name, last_order FROM members;`)
	if err != nil {
		return nil, fmt.Errorf("memberRepo.GetAll: query: %w", err)
	}
	defer rows.Close()

	var members []orderdomain.Member
	for rows.Next() {
		var m orderdomain.Member
		if err := rows.Scan(&m.ID, &m.Club, &m.Name, &m.LastOrder); err != nil {
			return nil, fmt.Errorf("memberRepo.GetAll: scan: %w", err)
		}
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("memberRepo.GetAll: rows: %w", err)
	}

	return members, nil
}

func (mr *memberRepo) Get(id int) (orderdomain.Member, error) {
	row := mr.db.QueryRow(`SELECT id, club, name, last_order FROM members WHERE id = $1;`, id)

	var m orderdomain.Member
	err := row.Scan(&m.ID, &m.Club, &m.Name, &m.LastOrder)
	if errors.Is(err, sql.ErrNoRows) {
		return orderdomain.Member{}, repo.ErrMemberNotFound
	}
	if err != nil {
		return orderdomain.Member{}, fmt.Errorf("memberRepo.Get(%d): %w", id, err)
	}

	return m, nil
}

func (mr *memberRepo) Create(member orderdomain.Member) (int, error) {
	if member.Name == `` || member.Club == domain.ClubUnknown {
		return 0, fmt.Errorf(`%w: %#v`, repo.ErrMemberFieldsNotFilled, member)
	}

	var id int
	if err := mr.db.QueryRow(
		`INSERT INTO members (name, club) VALUES ($1,$2) RETURNING id;`, member.Name, member.Club,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("memberRepo.Create: %w", err)
	}

	mr.logger.Info("member created",
		zap.Int("id", id),
		zap.String("name", member.Name),
		zap.String("club", member.Club.String()),
	)

	return id, nil
}

func (mr *memberRepo) Update(member orderdomain.Member) error {
	res, err := mr.db.Exec(
		`UPDATE members SET club = $1, name = $2, last_order = $3 WHERE id = $4;`,
		member.Club, member.Name, member.LastOrder, member.ID,
	)
	if err != nil {
		return fmt.Errorf("memberRepo.Update(%d): %w", member.ID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("memberRepo.Update(%d): rows affected: %w", member.ID, err)
	}
	if affected == 0 {
		return repo.ErrMemberNotFound
	}

	mr.logger.Info("member updated", zap.Int("id", member.ID), zap.String("name", member.Name))
	return nil
}

func (mr *memberRepo) Delete(id int) error {
	res, err := mr.db.Exec(`DELETE FROM members WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("memberRepo.Delete(%d): %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("memberRepo.Delete(%d): rows affected: %w", id, err)
	}
	if affected == 0 {
		return repo.ErrMemberNotFound
	}

	mr.logger.Info("member deleted", zap.Int("id", id))
	return nil
}
