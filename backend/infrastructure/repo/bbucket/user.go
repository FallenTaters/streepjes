package bbucket

// import (
// 	"errors"

// 	"github.com/FallenTaters/bbucket"
// 	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
// 	"github.com/PotatoesFall/vecty-test/domain"

// 	"go.etcd.io/bbolt"
// )

// func NewUserRepo(db *bbolt.DB) repo.User {
// 	return userRepo{
// 		bbucket.New(db, userBucket),
// 	}
// }

// type userRepo struct {
// 	bucket bbucket.Bucket
// }

// func (ur userRepo) GetByID(id int) (domain.User, bool) {
// 	var u domain.User
// 	err := ur.bucket.Get(bbucket.Itob(id), &u)
// 	if errors.Is(err, bbucket.ErrObjectNotFound) {
// 		return domain.User{}, false
// 	}
// 	if err != nil {
// 		panic(err)
// 	}

// 	return u, true
// }

// func (ur userRepo) GetAll() []domain.User {
// 	var users []domain.User
// 	var u domain.User
// 	err := ur.bucket.GetAll(&u, func(ptr interface{}) error {
// 		users = append(users, *ptr.(*domain.User))
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	return users
// }

// func (ur userRepo) Delete(id int)
