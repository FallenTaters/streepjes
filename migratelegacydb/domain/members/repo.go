package members

import (
	"errors"

	"github.com/FallenTaters/bbucket"
	"github.com/FallenTaters/streepjes/migratelegacydb/shared/buckets"
)

var repo memberRepo = defaultRepo{}

type defaultRepo struct{}

type memberRepo interface {
	getAll() ([]Member, error)
	get(id int) (Member, error)
	updateMember(member Member) error
	addMember(member Member) error
	deleteMember(id int) error
}

var ErrMemberNotFound = errors.New("member not found")

func (defaultRepo) getAll() ([]Member, error) {
	members := []Member{}

	return members, buckets.Members.GetAll(&Member{}, func(ptr interface{}) error {
		members = append(members, *ptr.(*Member))
		return nil
	})
}

func (defaultRepo) get(id int) (Member, error) {
	var member Member
	err := buckets.Members.Get(buckets.Itob(id), &member)
	if err == bbucket.ErrObjectNotFound {
		return Member{}, ErrMemberNotFound
	}

	return member, err
}

func (defaultRepo) updateMember(member Member) error {
	var m Member
	return buckets.Members.Update(member.Key(), &m, func(ptr interface{}) (object interface{}, err error) {
		return Member{
			ID:   m.ID,
			Club: member.Club,
			Name: member.Name,
			Debt: member.Debt,
		}, nil
	})
}

func (defaultRepo) addMember(member Member) error {
	member.ID = buckets.Members.NextSequence()
	return buckets.Members.Create(member.Key(), member)
}

func (defaultRepo) deleteMember(id int) error {
	return buckets.Members.Delete(buckets.Itob(id))
}
