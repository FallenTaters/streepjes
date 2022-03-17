package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func NewMemberRepo(db Queryable) repo.Member {
	return &memberRepo{
		db: db,
	}
}

type memberRepo struct {
	db Queryable
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
