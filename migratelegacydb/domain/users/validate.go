package users

import (
	"errors"

	"github.com/FallenTaters/bbucket"
	"github.com/FallenTaters/streepjes/migratelegacydb/shared"
)

var (
	ErrClubUnknown            = errors.New(`club unkown`)
	ErrNotAuthorized          = errors.New(`not authorized`)
	ErrEmptyPassword          = errors.New(`empty password`)
	ErrCannotChangeClub       = errors.New(`cannot change club`)
	ErrUserHasOpenOrders      = errors.New(`cannot delete user with open orders`)
	ErrCannotChangeOwnAccount = errors.New(`cannot change your own account`)
)

func validatePutUser(u User) (User, error) {
	user := User{
		Username: u.Username,
		Club:     u.Club,
		Name:     u.Name,
		Role:     u.Role,
	}

	original, err := get(u.Username)
	switch err {
	case nil:
		if u.Club != original.Club && original.Club != shared.ClubUnknown {
			return User{}, ErrCannotChangeClub
		}

		user.Password = nil

	case bbucket.ErrObjectNotFound:
		if len(u.Password) == 0 {
			return User{}, ErrEmptyPassword
		}

		user.Password = u.Password

	default:
		panic(err)
	}

	if user.Club == shared.ClubUnknown && user.Role != RoleBartender {
		return User{}, ErrClubUnknown
	}

	if user.Role == RoleNotAuthorized {
		return User{}, ErrNotAuthorized
	}

	return user, nil
}

func validateDeleteUser(username string) error {
	_, err := get(username)
	switch {
	case err == ErrUserNotFound:
		return err
	case err != nil:
		panic(err)
	}

	return nil
}
