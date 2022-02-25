package bbucket

import (
	"errors"

	"github.com/FallenTaters/bbucket"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"

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

func (mr memberRepo) GetAll() []orderdomain.Member {
	members := []orderdomain.Member{}

	err := mr.bucket.GetAll(&orderdomain.Member{}, func(ptr interface{}) error {
		members = append(members, *ptr.(*orderdomain.Member))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return members
}

func (mr memberRepo) Get(id int) (orderdomain.Member, bool) {
	var member orderdomain.Member

	err := mr.bucket.Get(itob(id), &member)
	if errors.Is(err, bbucket.ErrObjectNotFound) {
		return orderdomain.Member{}, false
	}
	if err != nil {
		panic(err)
	}

	return member, true
}

func (mr memberRepo) UpdateMember(member orderdomain.Member) error {
	var m orderdomain.Member

	err := mr.bucket.Update(memberKey(member), &m, func(ptr interface{}) (object interface{}, err error) {
		return member, nil
	})
	if errors.Is(err, bbucket.ErrObjectNotFound) {
		return repo.ErrMemberNotFound
	}

	return err
}

func (mr memberRepo) AddMember(member orderdomain.Member) (int, error) {
	member.ID = mr.bucket.NextSequence()
	return member.ID, mr.bucket.Create(memberKey(member), member)
}

func (mr memberRepo) DeleteMember(id int) bool {
	return mr.bucket.Delete(itob(id)) != nil
}

func memberKey(member orderdomain.Member) []byte {
	return itob(member.ID)
}
