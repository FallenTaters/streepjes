package orders

import (
	"errors"

	"github.com/FallenTaters/streepjes/migratelegacydb/domain/members"
	"github.com/FallenTaters/streepjes/migratelegacydb/domain/users"
	"github.com/FallenTaters/streepjes/migratelegacydb/shared"
)

var (
	ErrUnknownClub         = errors.New(`unkown club`)
	ErrNoContents          = errors.New(`order has no content`)
	ErrStatusNotOpenOrPaid = errors.New("order status must be open or paid")
	ErrPriceTooLow         = errors.New("order price must be greater than 0")
)

func validateAddOrder(o Order) error {
	// check if user is a bartender
	_, err := users.Get(o.Bartender)
	switch err {
	case nil:
	case users.ErrUserNotFound:
		return err

	default:
		panic(err)
	}

	// check if member exists
	if o.MemberID != 0 {
		_, err = members.Get(o.MemberID)
		switch err {
		case nil:
		case members.ErrMemberNotFound:
			return err

		default:
			panic(err)
		}
	}

	if o.Status != OrderStatusOpen && o.Status != OrderStatusPaid {
		return ErrStatusNotOpenOrPaid
	}

	if o.Club == shared.ClubUnknown {
		return ErrUnknownClub
	}

	if o.Contents == `` {
		return ErrNoContents
	}

	if o.Price <= 0 {
		return ErrPriceTooLow
	}

	return nil
}
