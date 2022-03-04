package order

import (
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Service interface {
	// GetAllMembers gets all the members
	GetAllMembers() []orderdomain.Member
}

func New(memberRepo repo.Member) Service {
	return &service{memberRepo}
}

type service struct {
	members repo.Member
}

func (s *service) GetAllMembers() []orderdomain.Member {
	return s.members.GetAll()
}
