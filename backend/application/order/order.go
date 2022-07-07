package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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

	// GetLeaderboard makes a leaderboard, the members will be pre-sorted by amount
	GetLeaderboard(api.LeaderboardFilter) api.Leaderboard

	// NewMember creates a new member
	NewMember(orderdomain.Member) error

	// EditMember edits a member
	EditMember(orderdomain.Member) error

	// DeleteMember deletes a member by id
	DeleteMember(id int) error
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
		return api.MemberDetails{}, false //nolint:exhaustruct
	}
	memberDetails.Member = member

	month := orderdomain.CurrentMonth()

	orders := s.orders.Filter(repo.OrderFilter{ //nolint:exhaustivestruct,exhaustruct
		MemberID:  id,
		Start:     month.Start(),
		End:       month.End(),
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
	month := orderdomain.CurrentMonth()

	return s.orders.Filter(repo.OrderFilter{ //nolint:exhaustivestruct,exhaustruct
		BartenderID: id,
		Start:       month.Start(),
		End:         month.End(),
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

func (s *service) GetLeaderboard(filter api.LeaderboardFilter) api.Leaderboard {
	orderFilter := repo.OrderFilter{ //nolint:exhaustruct,exhaustivestruct
		Start: filter.Start,
		End:   filter.End,
	}

	members := s.members.GetAll()
	orders := s.orders.Filter(orderFilter)

	return makeLeaderboard(members, orders)
}

func makeLeaderboard(members []orderdomain.Member, orders []orderdomain.Order) api.Leaderboard {
	var total orderdomain.Price
	totals := make(map[int]orderdomain.Price)
	counts := make(map[int]map[string]int)
	totalCounts := make(map[string]int)

	for _, o := range orders {
		if o.Status == orderdomain.StatusCancelled {
			continue
		}

		total += o.Price
		totals[o.MemberID] += o.Price

		// attempt to unmarshal contents and count by item name
		var lines []orderdomain.Line
		if err := json.Unmarshal([]byte(o.Contents), &lines); err != nil {
			continue
		}

		if m, ok := counts[o.MemberID]; !ok || m == nil {
			counts[o.MemberID] = make(map[string]int)
		}

		for _, line := range lines {
			counts[o.MemberID][line.Item.Name] += line.Amount
			totalCounts[line.Item.Name] += line.Amount
		}
	}

	leaderboard := api.Leaderboard{
		TotalPrice: total,
		Members:    make([]api.LeaderboardMember, 0, len(members)),
		Items:      totalCounts,
	}

	for _, member := range members {
		leaderboard.Members = append(leaderboard.Members, api.LeaderboardMember{
			Member:  member,
			Total:   totals[member.ID],
			Amounts: counts[member.ID],
		})
	}

	sort.Slice(leaderboard.Members, func(i, j int) bool {
		return leaderboard.Members[i].Total >= leaderboard.Members[j].Total
	})

	return leaderboard
}

func (s *service) NewMember(m orderdomain.Member) error {
	_, err := s.members.Create(m)
	return err
}

func (s *service) EditMember(m orderdomain.Member) error {
	original, ok := s.members.Get(m.ID)
	if !ok {
		return repo.ErrMemberNotFound
	}

	if original.Club != m.Club {
		return repo.ErrClubChange
	}

	return s.members.Update(m)
}

func (s *service) DeleteMember(id int) error {
	_, ok := s.members.Get(id)
	if !ok {
		return repo.ErrMemberNotFound
	}

	if len(s.orders.Filter(repo.OrderFilter{MemberID: id})) > 0 { //nolint:exhaustivestruct
		return repo.ErrMemberHasOrders
	}

	if !s.members.Delete(id) {
		return repo.ErrMemberNotFound
	}

	return nil
}
