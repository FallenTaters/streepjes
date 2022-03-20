package repo

import (
	"errors"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrMemberNameTaken       = errors.New(`member name taken for club`)
	ErrMemberFieldsNotFilled = errors.New(`member fields not filled`)
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

	// // Delete a member by id. Returns false if member is not found
	// DeleteMember(id int) bool
}

var ErrMemberNotFound = errors.New("member not found")
