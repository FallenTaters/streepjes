package auth_test

import (
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/mockdb"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Username:     `username`,
		PasswordHash: auth.HashPassword(`password`),
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

	cleanup := func() {
		updateCalled = false
	}

	t.Run(`correct login`, func(t *testing.T) {
		assert := assert.New(t)
		updateCalled = false

		user, ok := s.Login(`username`, `password`)

		assert.True(updateCalled)
		assert.True(ok)
		assert.Cmp(testUser, user, cmpopts.IgnoreFields(authdomain.User{}, `AuthTime`, `AuthToken`))

		cleanup()
	})

	t.Run(`wrong username or password`, func(t *testing.T) {
		assert := assert.New(t)

		_, ok := s.Login(`emanresu`, `password`)
		assert.False(ok)

		_, ok = s.Login(`username`, `drowssap`)
		assert.False(ok)

		cleanup()
	})

	t.Run(`panic on repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound // not expected, but any error will do
		}

		defer func() {
			v := recover()
			assert.Eq(repo.ErrUserNotFound, v)
		}()

		s.Login(`username`, `password`)

		cleanup()
	})
}

func TestCheck(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Username:     `username`,
		PasswordHash: auth.HashPassword(`password`),
		AuthToken:    `abcdefg`,
		AuthTime:     time.Now().Add(-time.Minute),
	}

	mock.GetByTokenFunc = func(token string) (authdomain.User, bool) {
		if token == `abcdefg` {
			return testUser, true
		}

		return authdomain.User{}, false
	}

	var updateCalled bool
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalled = true
		return nil
	}

	cleanup := func() {
		testUser.AuthTime = time.Now().Add(-time.Minute)
		updateCalled = false
	}

	t.Run(`valid token`, func(t *testing.T) {
		assert := assert.New(t)

		user, ok := s.Check(`abcdefg`)
		assert.True(ok)
		assert.Cmp(testUser, user, cmpopts.IgnoreFields(authdomain.User{}, `AuthTime`))
		assert.True(user.AuthTime.After(time.Now().Add(-time.Second)))
		assert.True(updateCalled)

		cleanup()
	})

	t.Run(`expired token`, func(t *testing.T) {
		assert := assert.New(t)

		testUser.AuthTime = time.Now().Add(-authdomain.TokenDuration - time.Second)

		_, ok := s.Check(`abcdefg`)
		assert.False(ok)
		assert.False(updateCalled)

		cleanup()
	})

	t.Run(`no or unknown token`, func(t *testing.T) {
		assert := assert.New(t)

		_, ok := s.Check(``)
		assert.False(ok)
		assert.False(updateCalled)

		_, ok = s.Check(`gfedcba`)
		assert.False(ok)

		cleanup()
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound // not expected but any error will do
		}

		_, ok := s.Check(`abcdefg`)
		assert.False(ok)
		assert.False(updateCalled)

		cleanup()
	})
}

func TestActive(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
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
	cleanup := func() {
		updateCalledWith = authdomain.User{}
	}

	t.Run(`set user active --> update authtime`, func(t *testing.T) {
		assert := assert.New(t)

		s.Active(1)
		assert.Eq(testUser.ID, updateCalledWith.ID)
		assert.True(updateCalledWith.AuthTime.After(time.Now().Add(-time.Second)))

		cleanup()
	})

	t.Run(`wrong id`, func(t *testing.T) {
		assert := assert.New(t)

		s.Active(2)
		assert.Eq(0, updateCalledWith.ID)

		cleanup()
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:        1,
		AuthToken: `abcdefg`,
		AuthTime:  time.Now().Add(-time.Second),
	}

	mock.GetFunc = func(i int) (authdomain.User, bool) {
		if i == 1 {
			return testUser, true
		}

		return authdomain.User{}, false
	}

	var updateCalledWith authdomain.User
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalledWith = user
		return nil
	}
	cleanup := func() {
		updateCalledWith = authdomain.User{}
	}

	t.Run(`log out logged in user`, func(t *testing.T) {
		assert := assert.New(t)

		s.Logout(1)
		assert.Eq(1, updateCalledWith.ID)
		assert.Eq(``, updateCalledWith.AuthToken)
		assert.True(updateCalledWith.AuthTime.Before(time.Now().Add(-authdomain.TokenDuration + time.Second)))

		cleanup()
	})

	t.Run(`log out non-existent user --> no panic, no update`, func(t *testing.T) {
		assert := assert.New(t)

		s.Logout(69)
		assert.Eq(authdomain.User{}, updateCalledWith)

		mock.GetFunc = func(i int) (authdomain.User, bool) { return authdomain.User{}, false }
		s.Logout(1)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})
}

func TestRegister(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		Username: `username`,
	}

	var createCalledWith authdomain.User
	mock.CreateFunc = func(user authdomain.User) (int, error) {
		createCalledWith = user
		return 0, nil
	}
	cleanup := func() {
		createCalledWith = authdomain.User{}
	}

	t.Run(`register`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.Register(testUser, `abc`)

		assert.NoError(err)
		assert.Eq(`username`, createCalledWith.Username)
		assert.True(auth.CheckPassword(createCalledWith.PasswordHash, `abc`))

		cleanup()
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.CreateFunc = func(user authdomain.User) (int, error) {
			return 0, repo.ErrUsernameTaken
		}

		err := s.Register(testUser, `abc`)

		assert.Error(err)
		assert.Eq(authdomain.User{}, createCalledWith)

		cleanup()
	})
}

func TestChangePassword(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		PasswordHash: auth.HashPassword(`abc`),
	}

	var updateCalledWith authdomain.User
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalledWith = user
		return nil
	}
	cleanup := func() {
		updateCalledWith = authdomain.User{}
	}

	t.Run(`change password of existing user`, func(t *testing.T) {
		assert := assert.New(t)

		ok := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abc`,
			New:      `cba`,
		})

		assert.True(ok)
		assert.Eq(1, updateCalledWith.ID)
		assert.True(auth.CheckPassword(updateCalledWith.PasswordHash, `cba`))

		cleanup()
	})

	t.Run(`wrong password`, func(t *testing.T) {
		assert := assert.New(t)

		ok := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abcasdfasdfasdf`,
			New:      `cba`,
		})

		assert.False(ok)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})

	t.Run(`empty password`, func(t *testing.T) {
		assert := assert.New(t)

		ok := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abc`,
			New:      ``,
		})

		assert.False(ok)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound
		}

		ok := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abc`,
			New:      `cba`,
		})

		assert.False(ok)

		cleanup()
	})
}

func TestChangeName(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Name:         `Hank`,
		PasswordHash: auth.HashPassword(`abc`),
	}

	var updateCalledWith authdomain.User
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalledWith = user
		return nil
	}
	cleanup := func() {
		updateCalledWith = authdomain.User{}
	}

	t.Run(`change name`, func(t *testing.T) {
		assert := assert.New(t)

		ok := s.ChangeName(testUser, `Dory`)

		assert.True(ok)
		assert.Eq(`Dory`, updateCalledWith.Name)

		cleanup()
	})

	t.Run(`empty name`, func(t *testing.T) {
		assert := assert.New(t)

		ok := s.ChangeName(testUser, ``)

		assert.False(ok)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound
		}

		ok := s.ChangeName(testUser, `Dory`)

		assert.False(ok)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})
}
