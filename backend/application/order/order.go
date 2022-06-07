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

	// GetCatalog fetches the catalog
	GetCatalog() api.Catalog

	// GetOrdersForBartender gets all order for that bartender for the current month
	GetOrdersForBartender(id int) []orderdomain.Order

	// PlaceOrder places the order for the bartender
	PlaceOrder(order orderdomain.Order, bartender authdomain.User) error

	// BartenderDeleteOrder marks an order as deleted for a bartender
	// if the bartender does not have access or if the order is not found, it returns false
	BartenderDeleteOrder(bartenderID, orderID int) bool

	// NewCategory creates a new category
	// It can return repo.ErrCategoryNameEmty or repo.ErrCategoryNameTaken
	NewCategory(orderdomain.Category) error

	// UpdateCategory updates an existing category.
	// It can return repo.ErrCategoryNotFound, repo.ErrCategoryNameEmpty or repo.ErrCategoryNameTaken
	UpdateCategory(orderdomain.Category) error

	// DeleteCategory deletes an existing category.
	// It can return repo.ErrCategoryNotFound or repo.ErrCategoryHasItems
	DeleteCategory(id int) error

	// NewItem creates a new item.
	// It can return repo.ErrCategoryNotFound, repo.ErrItemNameEmpty or repo.ErrItemNameTaken
	NewItem(orderdomain.Item) error

	// UpdateItem updates an existing item.
	// It can return repo.ErrItemNameTaken, repo.ErrItemNameEmpty, repo.ErrCategoryNotFound, or repo.ErrItemNotFound
	UpdateItem(orderdomain.Item) error

	// DeleteItem deletes an item.
	// It can return repo.ErrItemNotFound
	DeleteItem(id int) error
}

func New(memberRepo repo.Member, orderRepo repo.Order, catalogRepo repo.Catalog) Service {
	return &service{
		members: memberRepo,
		orders:  orderRepo,
		catalog: catalogRepo,
	}
}

type service struct {
	catalog repo.Catalog
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

	orders := s.orders.Filter(repo.OrderFilter{ //nolint:exhaustivestruct
		MemberID:  id,
		Month:     orderdomain.CurrentMonth(),
		StatusNot: []orderdomain.Status{orderdomain.StatusCancelled},
	})

	for _, order := range orders {
		memberDetails.Debt += order.Price
	}

	return memberDetails, true
}

func (s *service) GetCatalog() api.Catalog {
	return api.Catalog{
		Categories: s.catalog.GetCategories(),
		Items:      s.catalog.GetItems(),
	}
}

func (s *service) GetOrdersForBartender(id int) []orderdomain.Order {
	return s.orders.Filter(repo.OrderFilter{ //nolint:exhaustivestruct
		BartenderID: id,
		Month:       orderdomain.MonthOf(time.Now()),
	})
}

func (s *service) PlaceOrder(order orderdomain.Order, bartender authdomain.User) error {
	if order.Club == domain.ClubUnknown {
		return fmt.Errorf(`%w: club is %s`, ErrInvalidOrder, order.Club)
	}
	if order.Price < 0 {
		return fmt.Errorf(`%w: price is %s`, ErrInvalidOrder, order.Price)
	}

	if order.MemberID != 0 {
		member, ok := s.members.Get(order.MemberID)
		if !ok {
			return repo.ErrMemberNotFound
		}

		member.LastOrder = time.Now()
		// ignore error to avoid successful order being reported as failed
		_ = s.members.Update(member)
	}

	order.BartenderID = bartender.ID
	order.Status = orderdomain.StatusOpen
	order.OrderTime = time.Now()
	order.StatusTime = order.OrderTime

	_, err := s.orders.Create(order)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) BartenderDeleteOrder(bartenderID, orderID int) bool {
	order, ok := s.orders.Get(orderID)
	if !ok || order.BartenderID != bartenderID {
		return false
	}

	return s.orders.Delete(order.ID)
}

func (s *service) NewCategory(cat orderdomain.Category) error {
	if cat.Name == `` {
		return repo.ErrCategoryNameEmpty
	}

	_, err := s.catalog.CreateCategory(cat)
	return err
}

func (s *service) UpdateCategory(update orderdomain.Category) error {
	if update.Name == `` {
		return repo.ErrCategoryNameEmpty
	}

	return s.catalog.UpdateCategory(update)
}

func (s *service) DeleteCategory(id int) error {
	return s.catalog.DeleteCategory(id)
}

func (s *service) NewItem(item orderdomain.Item) error {
	if item.Name == `` {
		return repo.ErrItemNameEmpty
	}

	_, err := s.catalog.CreateItem(item)
	return err
}

func (s *service) UpdateItem(update orderdomain.Item) error {
	if update.Name == `` {
		return repo.ErrItemNameEmpty
	}

	return s.catalog.UpdateItem(update)
}

func (s *service) DeleteItem(id int) error {
	return s.catalog.DeleteItem(id)
}
