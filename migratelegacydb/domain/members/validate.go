package members

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared"
)

var (
	ErrUnknownClub  = errors.New("member must have club")
	ErrEmptyName    = errors.New("name may not be empty")
	ErrNameTaken    = errors.New("this name is taken for this club")
	ErrUnpaidOrders = errors.New("member still has unpaid orders")
)

func validateMember(member Member) error {
	original, err := repo.get(member.ID)
	switch err {
	case nil:
		if (original.Name != member.Name || original.Club != member.Club) && memberClubExists(member) {
			return ErrNameTaken
		}
	case ErrMemberNotFound:
		if member.ID != 0 {
			return ErrMemberNotFound
		}
		if memberClubExists(member) {
			return ErrNameTaken
		}
	default:
		panic(err)
	}

	if member.Club == shared.ClubUnknown {
		return ErrUnknownClub
	}

	if member.Name == `` {
		return ErrEmptyName
	}

	return nil
}

func memberClubExists(member Member) bool {
	members, err := repo.getAll()
	if err != nil {
		panic(err)
	}

	for _, m := range members {
		if m.Name == member.Name && m.Club == member.Club {
			return true
		}
	}

	return false
}
