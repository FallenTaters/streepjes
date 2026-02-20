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
	s := auth.New(mock, nil)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Username:     `username`,
		PasswordHash: auth.HashPassword(`password`),
	}

	mock.GetByUsernameFunc = func(username string) (authdomain.User, error) {
		if username == `username` {
			return testUser, nil
		}

		return authdomain.User{}, repo.ErrUserNotFound
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

		user, err := s.Login(`username`, `password`)

		assert.True(updateCalled)
		assert.NoError(err)
		assert.Cmp(testUser, user, cmpopts.IgnoreFields(authdomain.User{}, `AuthTime`, `AuthToken`))

		cleanup()
	})

	t.Run(`wrong username or password`, func(t *testing.T) {
		assert := assert.New(t)

		_, err := s.Login(`emanresu`, `password`)
		assert.Error(err)

		_, err = s.Login(`username`, `drowssap`)
		assert.Error(err)

		cleanup()
	})

	t.Run(`repo error on update`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound
		}

		_, err := s.Login(`username`, `password`)
		assert.Error(err)

		mock.UpdateFunc = func(user authdomain.User) error {
			updateCalled = true
			return nil
		}

		cleanup()
	})
}

func TestCheck(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock, nil)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:           1,
		Username:     `username`,
		PasswordHash: auth.HashPassword(`password`),
		AuthToken:    `abcdefg`,
		AuthTime:     time.Now().Add(-time.Minute),
	}

	mock.GetByTokenFunc = func(token string) (authdomain.User, error) {
		if token == `abcdefg` {
			return testUser, nil
		}

		return authdomain.User{}, repo.ErrUserNotFound
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

		user, err := s.Check(`abcdefg`)
		assert.NoError(err)
		assert.Cmp(testUser, user, cmpopts.IgnoreFields(authdomain.User{}, `AuthTime`))
		assert.True(user.AuthTime.After(time.Now().Add(-time.Second)))
		assert.True(updateCalled)

		cleanup()
	})

	t.Run(`expired token`, func(t *testing.T) {
		assert := assert.New(t)

		testUser.AuthTime = time.Now().Add(-authdomain.TokenDuration - time.Second)

		_, err := s.Check(`abcdefg`)
		assert.Error(err)
		assert.False(updateCalled)

		cleanup()
	})

	t.Run(`no or unknown token`, func(t *testing.T) {
		assert := assert.New(t)

		_, err := s.Check(``)
		assert.Error(err)
		assert.False(updateCalled)

		_, err = s.Check(`gfedcba`)
		assert.Error(err)

		cleanup()
	})

	t.Run(`repo error on update`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound
		}

		_, err := s.Check(`abcdefg`)
		assert.Error(err)

		mock.UpdateFunc = func(user authdomain.User) error {
			updateCalled = true
			return nil
		}

		cleanup()
	})
}

func TestActive(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock, nil)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID: 1,
	}

	mock.GetFunc = func(i int) (authdomain.User, error) {
		if i == 1 {
			return testUser, nil
		}

		return authdomain.User{}, repo.ErrUserNotFound
	}

	var updateCalledWith authdomain.User
	mock.UpdateFunc = func(user authdomain.User) error {
		updateCalledWith = user
		return nil
	}
	cleanup := func() {
		updateCalledWith = authdomain.User{}
	}

	t.Run(`set user active --> update authtime`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.Active(1)
		assert.NoError(err)
		assert.Eq(testUser.ID, updateCalledWith.ID)
		assert.True(updateCalledWith.AuthTime.After(time.Now().Add(-time.Second)))

		cleanup()
	})

	t.Run(`wrong id`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.Active(2)
		assert.NoError(err)
		assert.Eq(0, updateCalledWith.ID)

		cleanup()
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock, nil)
	testUser := authdomain.User{ //nolint:exhaustivestruct
		ID:        1,
		AuthToken: `abcdefg`,
		AuthTime:  time.Now().Add(-time.Second),
	}

	mock.GetFunc = func(i int) (authdomain.User, error) {
		if i == 1 {
			return testUser, nil
		}

		return authdomain.User{}, repo.ErrUserNotFound
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

		err := s.Logout(1)
		assert.NoError(err)
		assert.Eq(1, updateCalledWith.ID)
		assert.Eq(``, updateCalledWith.AuthToken)
		assert.True(updateCalledWith.AuthTime.Before(time.Now().Add(-authdomain.TokenDuration + time.Second)))

		cleanup()
	})

	t.Run(`log out non-existent user --> no error, no update`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.Logout(69)
		assert.NoError(err)
		assert.Eq(authdomain.User{}, updateCalledWith)

		mock.GetFunc = func(i int) (authdomain.User, error) { return authdomain.User{}, repo.ErrUserNotFound }
		err = s.Logout(1)
		assert.NoError(err)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})
}

func TestRegister(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock, nil)
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
	s := auth.New(mock, nil)
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

		err := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abc`,
			New:      `cba`,
		})

		assert.NoError(err)
		assert.Eq(1, updateCalledWith.ID)
		assert.True(auth.CheckPassword(updateCalledWith.PasswordHash, `cba`))

		cleanup()
	})

	t.Run(`wrong password`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abcasdfasdfasdf`,
			New:      `cba`,
		})

		assert.Error(err)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})

	t.Run(`empty password`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abc`,
			New:      ``,
		})

		assert.Error(err)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound
		}

		err := s.ChangePassword(testUser, api.ChangePassword{
			Original: `abc`,
			New:      `cba`,
		})

		assert.Error(err)

		cleanup()
	})
}

func TestChangeName(t *testing.T) {
	t.Parallel()

	mock := &mockdb.User{}
	s := auth.New(mock, nil)
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

		err := s.ChangeName(testUser, `Dory`)

		assert.NoError(err)
		assert.Eq(`Dory`, updateCalledWith.Name)

		cleanup()
	})

	t.Run(`empty name`, func(t *testing.T) {
		assert := assert.New(t)

		err := s.ChangeName(testUser, ``)

		assert.Error(err)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})

	t.Run(`repo error`, func(t *testing.T) {
		assert := assert.New(t)

		mock.UpdateFunc = func(user authdomain.User) error {
			return repo.ErrUserNotFound
		}

		err := s.ChangeName(testUser, `Dory`)

		assert.Error(err)
		assert.Eq(authdomain.User{}, updateCalledWith)

		cleanup()
	})
}
