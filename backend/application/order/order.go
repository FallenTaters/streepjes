package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var ErrInvalidOrder = errors.New(`invalid order`)

type Service interface {
	// GetAllMembers gets all the members
	GetAllMembers() []orderdomain.Member

	// GetMemberDetails gets a member and fills in details
	GetMemberDetails(id int) (api.MemberDetails, bool)

	// PlaceOrder places the order for the bartender
	PlaceOrder(order orderdomain.Order, bartender authdomain.User) error
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

func (s *service) PlaceOrder(order orderdomain.Order, bartender authdomain.User) error {
	if order.Club == domain.ClubUnknown {
		return fmt.Errorf(`%w: club is %s`, ErrInvalidOrder, order.Club)
	}
	if order.Price < 0 {
		return fmt.Errorf(`%w: price is %s`, ErrInvalidOrder, order.Price)
	}

	member, ok := s.members.Get(order.MemberID)
	if !ok {
		return repo.ErrMemberNotFound
	}

	order.BartenderID = bartender.ID
	order.Status = orderdomain.StatusOpen
	order.OrderTime = time.Now()
	order.StatusTime = order.OrderTime

	_, err := s.orders.Create(order)
	if err != nil {
		return err
	}

	member.LastOrder = time.Now()

	// ignore error to avoid successful order being reported as failed
	_ = s.members.Update(member)

	return nil
}
