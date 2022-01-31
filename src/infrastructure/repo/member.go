package repo

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/src/domain"
)

type Member interface {
	// Get all members
	GetAll() []domain.Member

	// Get a specific member by ID, returns false if not found
	Get(id int) (domain.Member, bool)

	// Update a specific member. Returns ErrMemberNotFound if the ID is not found
	UpdateMember(member domain.Member) error

	// Add a new member and return the id.
	AddMember(member domain.Member) (int, error)

	// Delete a member by id. Returns false if member is not found
	DeleteMember(id int) bool
}

var ErrMemberNotFound = errors.New("member not found")
