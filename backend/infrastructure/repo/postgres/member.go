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
	return &memberRepo{
		db:     db,
		logger: logger,
	}
}

type memberRepo struct {
	db     Queryable
	logger *zap.Logger
}

func (mr *memberRepo) GetAll() []orderdomain.Member {
	rows, err := mr.db.Query(`SELECT id, club, name, last_order FROM members;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var members []orderdomain.Member
	for rows.Next() {
		var member orderdomain.Member
		err := rows.Scan(&member.ID, &member.Club, &member.Name, &member.LastOrder)
		if err != nil {
			panic(err)
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return members
}

func (mr *memberRepo) Create(member orderdomain.Member) (int, error) {
	if member.Name == `` || member.Club == domain.ClubUnknown {
		return 0, fmt.Errorf(`%w: %#v`, repo.ErrMemberFieldsNotFilled, member)
	}

	row := mr.db.QueryRow(`INSERT INTO members (name, club) VALUES ($1,$2) RETURNING id;`, member.Name, member.Club)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	mr.logger.Info("member created",
		zap.Int("id", id),
		zap.String("name", member.Name),
		zap.String("club", member.Club.String()),
	)

	return id, nil
}

func (mr *memberRepo) Get(id int) (orderdomain.Member, bool) {
	row := mr.db.QueryRow(`SELECT id, club, name, last_order FROM members WHERE id = $1;`, id)

	var member orderdomain.Member

	err := row.Scan(&member.ID, &member.Club, &member.Name, &member.LastOrder)
	if errors.Is(err, sql.ErrNoRows) {
		return orderdomain.Member{}, false
	}
	if err != nil {
		panic(err)
	}

	return member, true
}

func (mr *memberRepo) Update(member orderdomain.Member) error {
	res, err := mr.db.Exec(
		`UPDATE members SET club = $1, name = $2, last_order = $3 WHERE id = $4;`,
		member.Club, member.Name, member.LastOrder, member.ID,
	)
	if err != nil {
		panic(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected == 0 {
		return repo.ErrMemberNotFound
	}

	mr.logger.Info("member updated", zap.Int("id", member.ID), zap.String("name", member.Name))

	return nil
}

func (mr *memberRepo) Delete(id int) bool {
	res, err := mr.db.Exec(`DELETE FROM members WHERE id = $1;`, id)
	if err != nil {
		panic(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if affected != 0 {
		mr.logger.Info("member deleted", zap.Int("id", id))
	}

	return affected != 0
}
