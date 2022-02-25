package repo

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

type Member interface {
	// Get all members
	GetAll() []orderdomain.Member

	// Get a specific member by ID, returns false if not found
	Get(id int) (orderdomain.Member, bool)

	// Update a specific member. Returns ErrMemberNotFound if the ID is not found
	UpdateMember(member orderdomain.Member) error

	// Add a new member and return the id.
	AddMember(member orderdomain.Member) (int, error)

	// Delete a member by id. Returns false if member is not found
	DeleteMember(id int) bool
}

var ErrMemberNotFound = errors.New("member not found")
