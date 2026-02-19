package sqlite

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

	res, err := mr.db.Exec(`INSERT INTO members (name, club) VALUES (?,?);`, member.Name, member.Club)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	mr.logger.Info("member created",
		zap.Int64("id", id),
		zap.String("name", member.Name),
		zap.String("club", member.Club.String()),
	)

	return int(id), nil
}

func (mr *memberRepo) Get(id int) (orderdomain.Member, bool) {
	row := mr.db.QueryRow(`SELECT id, club, name, last_order FROM members WHERE id = ?;`, id)

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
		`UPDATE members SET club = ?, name = ?, last_order = ? WHERE id = ?;`,
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
	res, err := mr.db.Exec(`DELETE FROM members WHERE id = ?;`, id)
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
