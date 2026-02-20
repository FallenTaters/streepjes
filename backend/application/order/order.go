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

var (
	ErrInvalidOrder    = errors.New(`invalid order`)
	ErrOrderNotFound   = errors.New(`order not found`)
	ErrOrderNotAllowed = errors.New(`bartender does not have access to this order`)
)

type Service interface {
	GetAllMembers() ([]orderdomain.Member, error)
	GetMemberDetails(id int) (api.MemberDetails, error)
	GetCatalog() (api.Catalog, error)
	GetOrdersForBartender(id int) ([]orderdomain.Order, error)
	GetOrdersByClub(club domain.Club, month orderdomain.Month) ([]orderdomain.Order, error)
	BillingCSV(club domain.Club, month orderdomain.Month) ([]byte, error)
	PlaceOrder(order orderdomain.Order, bartender authdomain.User) error
	BartenderDeleteOrder(bartenderID, orderID int) error
	NewCategory(orderdomain.Category) error
	UpdateCategory(orderdomain.Category) error
	DeleteCategory(id int) error
	NewItem(orderdomain.Item) error
	UpdateItem(orderdomain.Item) error
	DeleteItem(id int) error
	GetLeaderboard(api.LeaderboardFilter) (api.Leaderboard, error)
	NewMember(orderdomain.Member) error
	EditMember(orderdomain.Member) error
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

func (s *service) GetAllMembers() ([]orderdomain.Member, error) {
	return s.members.GetAll()
}

func (s *service) GetMemberDetails(id int) (api.MemberDetails, error) {
	member, err := s.members.Get(id)
	if err != nil {
		return api.MemberDetails{}, fmt.Errorf("order.GetMemberDetails: %w", err)
	}

	month := orderdomain.CurrentMonth()

	orders, err := s.orders.Filter(repo.OrderFilter{
		MemberID:  id,
		Start:     month.Start(),
		End:       month.End(),
		StatusNot: []orderdomain.Status{orderdomain.StatusCancelled},
	})
	if err != nil {
		return api.MemberDetails{}, fmt.Errorf("order.GetMemberDetails: filter: %w", err)
	}

	var debt orderdomain.Price
	for _, order := range orders {
		debt += order.Price
	}

	return api.MemberDetails{
		Member: member,
		Debt:   debt,
	}, nil
}

func (s *service) GetCatalog() (api.Catalog, error) {
	categories, err := s.catalog.GetCategories()
	if err != nil {
		return api.Catalog{}, fmt.Errorf("order.GetCatalog: categories: %w", err)
	}

	items, err := s.catalog.GetItems()
	if err != nil {
		return api.Catalog{}, fmt.Errorf("order.GetCatalog: items: %w", err)
	}

	return api.Catalog{
		Categories: categories,
		Items:      items,
	}, nil
}

func (s *service) GetOrdersForBartender(id int) ([]orderdomain.Order, error) {
	month := orderdomain.CurrentMonth()

	orders, err := s.orders.Filter(repo.OrderFilter{
		BartenderID: id,
		Start:       month.Start(),
		End:         month.End(),
	})
	if err != nil {
		return nil, fmt.Errorf("order.GetOrdersForBartender: %w", err)
	}

	return orders, nil
}

func (s *service) GetOrdersByClub(club domain.Club, month orderdomain.Month) ([]orderdomain.Order, error) {
	orders, err := s.orders.Filter(repo.OrderFilter{
		Club:      club,
		Start:     month.Start(),
		End:       month.End(),
		StatusNot: []orderdomain.Status{orderdomain.StatusCancelled},
	})
	if err != nil {
		return nil, fmt.Errorf("order.GetOrdersByClub: %w", err)
	}

	for i := range orders {
		orders[i].OrderTime = orders[i].OrderTime.In(timezone)
	}

	return orders, nil
}

func (s *service) BillingCSV(club domain.Club, month orderdomain.Month) ([]byte, error) {
	orders, err := s.orders.Filter(repo.OrderFilter{
		Club:      club,
		Start:     month.Start(),
		End:       month.End(),
		StatusNot: []orderdomain.Status{orderdomain.StatusCancelled},
	})
	if err != nil {
		return nil, fmt.Errorf("order.BillingCSV: orders: %w", err)
	}

	members, err := s.members.GetAll()
	if err != nil {
		return nil, fmt.Errorf("order.BillingCSV: members: %w", err)
	}

	return writeCSV(orders, members), nil
}

func (s *service) PlaceOrder(order orderdomain.Order, bartender authdomain.User) error {
	if order.Club == domain.ClubUnknown {
		return fmt.Errorf(`%w: club is %s`, ErrInvalidOrder, order.Club)
	}
	if order.Price < 0 {
		return fmt.Errorf(`%w: price is %s`, ErrInvalidOrder, order.Price)
	}

	if order.MemberID != 0 {
		member, err := s.members.Get(order.MemberID)
		if err != nil {
			return fmt.Errorf("order.PlaceOrder: get member: %w", err)
		}

		member.LastOrder = time.Now()
		_ = s.members.Update(member)
	}

	order.BartenderID = bartender.ID
	order.Status = orderdomain.StatusOpen
	order.OrderTime = time.Now()
	order.StatusTime = order.OrderTime

	if _, err := s.orders.Create(order); err != nil {
		return fmt.Errorf("order.PlaceOrder: create: %w", err)
	}

	return nil
}

func (s *service) BartenderDeleteOrder(bartenderID, orderID int) error {
	order, err := s.orders.Get(orderID)
	if errors.Is(err, repo.ErrOrderNotFound) {
		return ErrOrderNotFound
	}
	if err != nil {
		return fmt.Errorf("order.BartenderDeleteOrder: %w", err)
	}

	if order.BartenderID != bartenderID {
		return ErrOrderNotAllowed
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

func (s *service) GetLeaderboard(filter api.LeaderboardFilter) (api.Leaderboard, error) {
	members, err := s.members.GetAll()
	if err != nil {
		return api.Leaderboard{}, fmt.Errorf("order.GetLeaderboard: members: %w", err)
	}

	orders, err := s.orders.Filter(repo.OrderFilter{
		Start: filter.Start,
		End:   filter.End,
	})
	if err != nil {
		return api.Leaderboard{}, fmt.Errorf("order.GetLeaderboard: orders: %w", err)
	}

	return makeLeaderboard(members, orders), nil
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
		return leaderboard.Members[i].Total > leaderboard.Members[j].Total
	})

	return leaderboard
}

func (s *service) NewMember(m orderdomain.Member) error {
	_, err := s.members.Create(m)
	return err
}

func (s *service) EditMember(m orderdomain.Member) error {
	original, err := s.members.Get(m.ID)
	if err != nil {
		return fmt.Errorf("order.EditMember: %w", err)
	}

	if original.Club != m.Club {
		return repo.ErrClubChange
	}

	return s.members.Update(m)
}

func (s *service) DeleteMember(id int) error {
	if _, err := s.members.Get(id); err != nil {
		return fmt.Errorf("order.DeleteMember: %w", err)
	}

	orders, err := s.orders.Filter(repo.OrderFilter{MemberID: id})
	if err != nil {
		return fmt.Errorf("order.DeleteMember: check orders: %w", err)
	}
	if len(orders) > 0 {
		return repo.ErrMemberHasOrders
	}

	return s.members.Delete(id)
}
