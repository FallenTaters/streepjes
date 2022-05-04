package auth

import (
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

type Service interface {
	// Login returns the user if credentials are correct
	// otherwise, it returns false
	Login(user, pass string) (authdomain.User, bool)

	// Check gets the user with the correct token
	// if the token is expired or unknown, it returns false
	Check(token string) (authdomain.User, bool)

	// Active refreshes a users token, setting AuthTime to now
	// if the user is not found, it is a no-op
	Active(id int)

	// Logout deletes the users AuthToken, if found
	// if the user is not found, it is a no-op
	Logout(id int)

	// Register registers a new user. It sets the passwordHash and ID
	// if the username is taken, it return repo.ErrUsernameTaken
	// if mandatory fields are missing, it returns repo.ErrUserMissingFields
	Register(user authdomain.User, password string) error

	// Update updates a user
	// it can return repo.ErrUserMissingFields and repo.ErrUsernameTaken
	// if password == ``, it doesn't edit password
	Update(user authdomain.User, password string) error

	// ChangePassword verifies the original password and changes it to the new password
	// if anything goes wrong, it returns false
	ChangePassword(user authdomain.User, changePassword api.ChangePassword) bool

	// ChangeName attempts to change the name of the user
	// if anything goes wrong, it returns false
	ChangeName(user authdomain.User, name string) bool

	// GetUsers gets all users
	GetUsers() []authdomain.User

	// Delete deletes a user. If the user doesn't exist or has orders,
	// it does not delete the user and return false.
	Delete(id int) bool
}

func New(userRepo repo.User, orderRepo repo.Order) Service {
	return &service{userRepo, orderRepo}
}

type service struct {
	users  repo.User
	orders repo.Order
}

func (s *service) Login(username, pass string) (authdomain.User, bool) {
	user, ok := s.users.GetByUsername(username)
	if !ok {
		return authdomain.User{}, false
	}

	if !CheckPassword(user.PasswordHash, pass) {
		return authdomain.User{}, false
	}

	user.AuthToken = generateToken()
	user.AuthTime = time.Now()

	err := s.users.Update(user)
	if err != nil {
		panic(err)
	}

	return user, true
}

func (s *service) Check(token string) (authdomain.User, bool) {
	if token == `` {
		return authdomain.User{}, false
	}

	user, ok := s.users.GetByToken(token)
	if !ok {
		return authdomain.User{}, false
	}

	if time.Since(user.AuthTime) > authdomain.TokenDuration {
		return authdomain.User{}, false
	}

	user.AuthTime = time.Now()
	err := s.users.Update(user)
	if err != nil {
		return authdomain.User{}, false
	}

	return user, true
}

func (s *service) Active(id int) {
	user, ok := s.users.Get(id)
	if !ok {
		return
	}

	user.AuthTime = time.Now()

	_ = s.users.Update(user)
}

func (s *service) Logout(id int) {
	user, ok := s.users.Get(id)
	if !ok {
		return
	}

	user.AuthToken = ``
	user.AuthTime = time.Now().Add(-authdomain.TokenDuration)

	_ = s.users.Update(user)
}

func (s *service) Register(user authdomain.User, password string) error {
	user.PasswordHash = HashPassword(password)

	_, err := s.users.Create(user)
	return err
}

func (s *service) Update(userChanges authdomain.User, password string) error {
	user, ok := s.users.Get(userChanges.ID)
	if !ok {
		return repo.ErrUserNotFound
	}

	if password != `` {
		user.PasswordHash = HashPassword(password)
	}

	// only update the following fields
	user.Name = userChanges.Name
	user.Username = userChanges.Username
	user.Club = userChanges.Club
	user.Role = userChanges.Role

	return s.users.Update(user)
}

func (s *service) ChangePassword(user authdomain.User, changePassword api.ChangePassword) bool {
	if changePassword.New == `` {
		return false
	}

	if !CheckPassword(user.PasswordHash, changePassword.Original) {
		return false
	}

	user.PasswordHash = HashPassword(changePassword.New)

	return s.users.Update(user) == nil
}

func (s *service) ChangeName(user authdomain.User, name string) bool {
	if name == `` {
		return false
	}

	user.Name = name

	return s.users.Update(user) == nil
}

func (s *service) GetUsers() []authdomain.User {
	return s.users.GetAll()
}

func (s *service) Delete(id int) bool {
	orders := s.orders.Filter(repo.OrderFilter{ //nolint:exhaustivestruct
		BartenderID: id,
	})

	if len(orders) > 0 {
		return false
	}

	return s.users.Delete(id) == nil
}
