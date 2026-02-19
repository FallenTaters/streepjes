package repo

import (
	"errors"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrMemberNameTaken       = errors.New(`member name taken for club`)
	ErrMemberFieldsNotFilled = errors.New(`member fields not filled`)
	ErrMemberHasOrders       = errors.New(`member has orders`)
	ErrClubChange            = errors.New(`existing member may not change club`)
)

type Member interface {
	// Get all members in the database
	GetAll() []orderdomain.Member

	// Get a specific member by ID, returns false if not found
	Get(id int) (orderdomain.Member, bool)

	// // Update a specific member. Returns ErrMemberNotFound if the ID is not found
	Update(member orderdomain.Member) error

	// Create a new member and return the id
	// if name is taken for the club, it returns ErrMemberNameTaken
	Create(member orderdomain.Member) (int, error)

	// Delete a member by id
	Delete(id int) error
}

var ErrMemberNotFound = errors.New("member not found")
