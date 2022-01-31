package bbucket

import (
	"errors"

	"github.com/FallenTaters/bbucket"
	"github.com/PotatoesFall/vecty-test/src/domain"
	"github.com/PotatoesFall/vecty-test/src/infrastructure/repo"

	"go.etcd.io/bbolt"
)

func NewMemberRepo(db *bbolt.DB) repo.Member {
	return memberRepo{
		bbucket.New(db, memberBucket),
	}
}

type memberRepo struct {
	bucket bbucket.Bucket
}

func (mr memberRepo) GetAll() []domain.Member {
	members := []domain.Member{}

	err := mr.bucket.GetAll(&domain.Member{}, func(ptr interface{}) error {
		members = append(members, *ptr.(*domain.Member))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return members
}

func (mr memberRepo) Get(id int) (domain.Member, bool) {
	var member domain.Member
	err := mr.bucket.Get(itob(id), &member)
	if errors.Is(err, bbucket.ErrObjectNotFound) {
		return domain.Member{}, false
	}
	if err != nil {
		panic(err)
	}

	return member, true
}

func (mr memberRepo) UpdateMember(member domain.Member) error {
	var m domain.Member
	err := mr.bucket.Update(memberKey(member), &m, func(ptr interface{}) (object interface{}, err error) {
		return member, nil
	})
	if errors.Is(err, bbucket.ErrObjectNotFound) {
		return repo.ErrMemberNotFound
	}

	return err
}

func (mr memberRepo) AddMember(member domain.Member) (int, error) {
	member.ID = mr.bucket.NextSequence()
	return member.ID, mr.bucket.Create(memberKey(member), member)
}

func (mr memberRepo) DeleteMember(id int) bool {
	return mr.bucket.Delete(itob(id)) != nil
}

func memberKey(member domain.Member) []byte {
	return itob(member.ID)
}
