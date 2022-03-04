package auth

import (
	"errors"
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo/mockdb"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Username:     `username`,
		PasswordHash: hashPassword(`password`),
	}

	mock.GetByUsernameFunc = func(username string) (authdomain.User, bool) {
		if username == `username` {
			return testUser, true
		}

		return authdomain.User{}, false
	}

	var updateCalled bool
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalled = true
		return nil
	}

	t.Run(`correct login`, func(t *testing.T) {
		assert := assert.New(t)
		updateCalled = false

		user, ok := s.Login(`username`, `password`)

		assert.True(updateCalled)
		assert.True(ok)
		assert.Cmp(testUser, user, cmpopts.IgnoreFields(authdomain.User{}, `AuthTime`, `AuthToken`))
	})

	t.Run(`wrong username or password`, func(t *testing.T) {
		assert := assert.New(t)

		_, ok := s.Login(`emanresu`, `password`)
		assert.False(ok)

		_, ok = s.Login(`username`, `drowssap`)
		assert.False(ok)
	})

	t.Run(`panic on repo error`, func(t *testing.T) {
		assert := assert.New(t)
		err := errors.New(`bla`) //nolint:goerr113

		mock.UpdateFunc = func(user authdomain.User) error {
			return err
		}

		defer func() {
			v := recover()
			assert.Eq(err, v)
		}()

		s.Login(`username`, `password`)
	})
}

func TestCheck(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Username:     `username`,
		PasswordHash: hashPassword(`password`),
		AuthToken:    `abcdefg`,
		AuthTime:     time.Now().Add(-time.Minute),
	}

	mock.GetByTokenFunc = func(token string) (authdomain.User, bool) {
		if token == `abcdefg` {
			return testUser, true
		}

		return authdomain.User{}, false
	}

	mock.UpdateFunc = func(user authdomain.User) error {
		return nil
	}

	t.Run(`valid token`, func(t *testing.T) {
		assert := assert.New(t)

		user, ok := s.Check(`abcdefg`)
		assert.True(ok)
		assert.Cmp(testUser, user, cmpopts.IgnoreFields(authdomain.User{}, `AuthTime`))
		assert.True(user.AuthTime.After(time.Now().Add(-time.Second)))
	})

	t.Run(`expired token`, func(t *testing.T) {
		assert := assert.New(t)

		testUser.AuthTime = time.Now().Add(-authdomain.TokenDuration - time.Second)

		_, ok := s.Check(`abcdefg`)
		assert.False(ok)

		testUser.AuthTime = time.Now()
	})

	t.Run(`no or unknown token`, func(t *testing.T) {
		assert := assert.New(t)

		_, ok := s.Check(``)
		assert.False(ok)

		_, ok = s.Check(`gfedcba`)
		assert.False(ok)
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		_, ok := s.Check(`abcdefg`)
		assert.True(ok)

		mock.UpdateFunc = func(user authdomain.User) error { return errors.New(`bla`) } //nolint:goerr113

		_, ok = s.Check(`abcdefg`)
		assert.False(ok)
	})
}

func TestActive(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID: 1,
	}

	mock.GetFunc = func(i int) (authdomain.User, bool) {
		if i == 1 {
			return testUser, true
		}

		return authdomain.User{}, false
	}

	var updateCalledWith authdomain.User
	var updateErr error
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalledWith = user
		return updateErr
	}

	t.Run(`set user active --> update authtime`, func(t *testing.T) {
		assert := assert.New(t)

		s.Active(1)
		assert.Eq(testUser.ID, updateCalledWith.ID)
		assert.True(updateCalledWith.AuthTime.After(time.Now().Add(-time.Second)))
		updateCalledWith = authdomain.User{}
	})

	t.Run(`wrong id`, func(t *testing.T) {
		assert := assert.New(t)

		s.Active(2)
		assert.Eq(0, updateCalledWith.ID)
	})
}

func TestLogout(t *testing.T) {
	// TODO
}

func TestRegister(t *testing.T) {
	// TODO
}

func TestChangePassword(t *testing.T) {
	// TODO
}

func TestChangeName(t *testing.T) {
	// TODO
}
