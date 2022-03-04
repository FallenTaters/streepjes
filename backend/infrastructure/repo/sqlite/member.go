package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func NewMemberRepo(db *sql.DB) repo.Member {
	return &memberRepo{
		db: db,
	}
}

type memberRepo struct {
	db *sql.DB
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
