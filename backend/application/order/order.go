package order

import (
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Service interface {
	// GetAllMembers gets all the members
	GetAllMembers() []orderdomain.Member

	// GetMemberDetails gets a member and fills in details
	GetMemberDetails(id int) (api.MemberDetails, bool)
}

func New(memberRepo repo.Member, orderRepo repo.Order) Service {
	return &service{
		members: memberRepo,
		orders:  orderRepo,
	}
}

type service struct {
	members repo.Member
	orders  repo.Order
}

func (s *service) GetAllMembers() []orderdomain.Member {
	return s.members.GetAll()
}

func (s *service) GetMemberDetails(id int) (api.MemberDetails, bool) {
	var memberDetails api.MemberDetails

	member, ok := s.members.Get(id)
	if !ok {
		return api.MemberDetails{}, false
	}
	memberDetails.Member = member

	orders := s.orders.Filter(repo.OrderFilter{MemberID: id, Month: orderdomain.CurrentMonth()})

	for _, order := range orders {
		memberDetails.Debt += order.Price
	}

	return memberDetails, true
}
