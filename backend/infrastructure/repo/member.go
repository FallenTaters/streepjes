package repo

import (
	"errors"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrMemberNotFound        = errors.New("member not found")
	ErrMemberNameTaken       = errors.New(`member name taken for club`)
	ErrMemberFieldsNotFilled = errors.New(`member fields not filled`)
	ErrMemberHasOrders       = errors.New(`member has orders`)
	ErrClubChange            = errors.New(`existing member may not change club`)
)

type Member interface {
	GetAll() ([]orderdomain.Member, error)
	Get(id int) (orderdomain.Member, error)
	Update(member orderdomain.Member) error
	Create(member orderdomain.Member) (int, error)
	Delete(id int) error
}
