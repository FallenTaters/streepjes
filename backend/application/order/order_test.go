package order

import (
	"errors"
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/mockdb"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestMakeLeaderboard(t *testing.T) {
	t.Parallel()

	members := []orderdomain.Member{
		{ID: 1, Name: "Alice", Club: domain.ClubGladiators},
		{ID: 2, Name: "Bob", Club: domain.ClubGladiators},
		{ID: 3, Name: "Carol", Club: domain.ClubGladiators},
	}

	t.Run("sorts members by total descending", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 2, Price: 300, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 3, Price: 200, Status: orderdomain.StatusOpen, Contents: `[]`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(3, len(lb.Members))
		assert.Eq("Bob", lb.Members[0].Name)
		assert.Eq("Carol", lb.Members[1].Name)
		assert.Eq("Alice", lb.Members[2].Name)
	})

	t.Run("stable sort for equal totals", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 2, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 3, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(3, len(lb.Members))
		assert.Eq(orderdomain.Price(100), lb.Members[0].Total)
		assert.Eq(orderdomain.Price(100), lb.Members[1].Total)
		assert.Eq(orderdomain.Price(100), lb.Members[2].Total)
	})

	t.Run("excludes cancelled orders", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 1, Price: 200, Status: orderdomain.StatusCancelled, Contents: `[]`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(orderdomain.Price(100), lb.TotalPrice)
		assert.Eq(orderdomain.Price(100), lb.Members[0].Total)
	})

	t.Run("counts items from order contents", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{
				MemberID: 1,
				Price:    300,
				Status:   orderdomain.StatusOpen,
				Contents: `[{"product":{"name":"Beer"},"amount":2},{"product":{"name":"Wine"},"amount":1}]`,
			},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(2, lb.Items["Beer"])
		assert.Eq(1, lb.Items["Wine"])
		assert.Eq(2, lb.Members[0].Amounts["Beer"])
		assert.Eq(1, lb.Members[0].Amounts["Wine"])
	})

	t.Run("handles malformed contents gracefully", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `not json`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(orderdomain.Price(100), lb.TotalPrice)
		assert.Eq(orderdomain.Price(100), lb.Members[0].Total)
	})

	t.Run("empty orders", func(t *testing.T) {
		assert := assert.New(t)

		lb := makeLeaderboard(members, nil)

		assert.Eq(orderdomain.Price(0), lb.TotalPrice)
		assert.Eq(3, len(lb.Members))
	})
}

func newTestService(memberMock *mockdb.Member, orderMock *mockdb.Order, catalogMock *mockdb.Catalog) Service {
	return New(memberMock, orderMock, catalogMock, time.UTC)
}

func TestPlaceOrder(t *testing.T) {
	t.Parallel()

	t.Run("valid order with member", func(t *testing.T) {
		assert := assert.New(t)

		var memberUpdated bool
		var orderCreated orderdomain.Order
		memberMock := &mockdb.Member{
			GetFunc: func(id int) (orderdomain.Member, error) {
				return orderdomain.Member{ID: id, Club: domain.ClubGladiators}, nil
			},
			UpdateFunc: func(_ orderdomain.Member) error {
				memberUpdated = true
				return nil
			},
		}
		orderMock := &mockdb.Order{
			CreateFunc: func(o orderdomain.Order) (int, error) {
				orderCreated = o
				return 1, nil
			},
		}

		svc := newTestService(memberMock, orderMock, nil)
		bartender := authdomain.User{ID: 5}

		err := svc.PlaceOrder(orderdomain.Order{
			Club:     domain.ClubGladiators,
			MemberID: 1,
			Price:    200,
		}, bartender)

		assert.NoError(err)
		assert.True(memberUpdated)
		assert.Eq(5, orderCreated.BartenderID)
		assert.Eq(orderdomain.StatusOpen, orderCreated.Status)
	})

	t.Run("rejects unknown club", func(t *testing.T) {
		assert := assert.New(t)

		svc := newTestService(nil, nil, nil)

		err := svc.PlaceOrder(orderdomain.Order{
			Club:  domain.ClubUnknown,
			Price: 200,
		}, authdomain.User{})

		assert.True(errors.Is(err, ErrInvalidOrder))
	})

	t.Run("rejects negative price", func(t *testing.T) {
		assert := assert.New(t)

		svc := newTestService(nil, nil, nil)

		err := svc.PlaceOrder(orderdomain.Order{
			Club:  domain.ClubGladiators,
			Price: -100,
		}, authdomain.User{})

		assert.True(errors.Is(err, ErrInvalidOrder))
	})

	t.Run("order without member skips member update", func(t *testing.T) {
		assert := assert.New(t)

		orderMock := &mockdb.Order{
			CreateFunc: func(_ orderdomain.Order) (int, error) { return 1, nil },
		}

		svc := newTestService(nil, orderMock, nil)

		err := svc.PlaceOrder(orderdomain.Order{
			Club:  domain.ClubGladiators,
			Price: 200,
		}, authdomain.User{ID: 1})

		assert.NoError(err)
	})
}

func TestBartenderDeleteOrder(t *testing.T) {
	t.Parallel()

	t.Run("deletes own order", func(t *testing.T) {
		assert := assert.New(t)

		var deleted bool
		orderMock := &mockdb.Order{
			GetFunc: func(id int) (orderdomain.Order, error) {
				return orderdomain.Order{ID: id, BartenderID: 5}, nil
			},
			DeleteFunc: func(_ int) error {
				deleted = true
				return nil
			},
		}

		svc := newTestService(nil, orderMock, nil)
		err := svc.BartenderDeleteOrder(5, 10)

		assert.NoError(err)
		assert.True(deleted)
	})

	t.Run("rejects deleting other bartender's order", func(t *testing.T) {
		assert := assert.New(t)

		orderMock := &mockdb.Order{
			GetFunc: func(id int) (orderdomain.Order, error) {
				return orderdomain.Order{ID: id, BartenderID: 99}, nil
			},
		}

		svc := newTestService(nil, orderMock, nil)
		err := svc.BartenderDeleteOrder(5, 10)

		assert.True(errors.Is(err, ErrOrderNotAllowed))
	})

	t.Run("order not found", func(t *testing.T) {
		assert := assert.New(t)

		orderMock := &mockdb.Order{
			GetFunc: func(_ int) (orderdomain.Order, error) {
				return orderdomain.Order{}, repo.ErrOrderNotFound
			},
		}

		svc := newTestService(nil, orderMock, nil)
		err := svc.BartenderDeleteOrder(5, 10)

		assert.True(errors.Is(err, ErrOrderNotFound))
	})
}

func TestEditMember(t *testing.T) {
	t.Parallel()

	t.Run("allows edit without club change", func(t *testing.T) {
		assert := assert.New(t)

		var updated orderdomain.Member
		memberMock := &mockdb.Member{
			GetFunc: func(_ int) (orderdomain.Member, error) {
				return orderdomain.Member{ID: 1, Club: domain.ClubGladiators, Name: "Old"}, nil
			},
			UpdateFunc: func(m orderdomain.Member) error {
				updated = m
				return nil
			},
		}

		svc := newTestService(memberMock, nil, nil)
		err := svc.EditMember(orderdomain.Member{ID: 1, Club: domain.ClubGladiators, Name: "New"})

		assert.NoError(err)
		assert.Eq("New", updated.Name)
	})

	t.Run("rejects club change", func(t *testing.T) {
		assert := assert.New(t)

		memberMock := &mockdb.Member{
			GetFunc: func(_ int) (orderdomain.Member, error) {
				return orderdomain.Member{ID: 1, Club: domain.ClubGladiators}, nil
			},
		}

		svc := newTestService(memberMock, nil, nil)
		err := svc.EditMember(orderdomain.Member{ID: 1, Club: domain.ClubParabool})

		assert.True(errors.Is(err, repo.ErrClubChange))
	})

	t.Run("member not found", func(t *testing.T) {
		assert := assert.New(t)

		memberMock := &mockdb.Member{
			GetFunc: func(_ int) (orderdomain.Member, error) {
				return orderdomain.Member{}, repo.ErrMemberNotFound
			},
		}

		svc := newTestService(memberMock, nil, nil)
		err := svc.EditMember(orderdomain.Member{ID: 99})

		assert.Error(err)
	})
}

func TestDeleteMember(t *testing.T) {
	t.Parallel()

	t.Run("deletes member without orders", func(t *testing.T) {
		assert := assert.New(t)

		var deleted bool
		memberMock := &mockdb.Member{
			GetFunc: func(_ int) (orderdomain.Member, error) {
				return orderdomain.Member{ID: 1}, nil
			},
			DeleteFunc: func(_ int) error {
				deleted = true
				return nil
			},
		}
		orderMock := &mockdb.Order{
			FilterFunc: func(_ repo.OrderFilter) ([]orderdomain.Order, error) {
				return nil, nil
			},
		}

		svc := newTestService(memberMock, orderMock, nil)
		err := svc.DeleteMember(1)

		assert.NoError(err)
		assert.True(deleted)
	})

	t.Run("rejects delete when member has orders", func(t *testing.T) {
		assert := assert.New(t)

		memberMock := &mockdb.Member{
			GetFunc: func(_ int) (orderdomain.Member, error) {
				return orderdomain.Member{ID: 1}, nil
			},
		}
		orderMock := &mockdb.Order{
			FilterFunc: func(_ repo.OrderFilter) ([]orderdomain.Order, error) {
				return []orderdomain.Order{{ID: 1}}, nil
			},
		}

		svc := newTestService(memberMock, orderMock, nil)
		err := svc.DeleteMember(1)

		assert.True(errors.Is(err, repo.ErrMemberHasOrders))
	})

	t.Run("member not found", func(t *testing.T) {
		assert := assert.New(t)

		memberMock := &mockdb.Member{
			GetFunc: func(_ int) (orderdomain.Member, error) {
				return orderdomain.Member{}, repo.ErrMemberNotFound
			},
		}

		svc := newTestService(memberMock, nil, nil)
		err := svc.DeleteMember(99)

		assert.Error(err)
	})
}

func TestNewCategory(t *testing.T) {
	t.Parallel()

	t.Run("rejects empty name", func(t *testing.T) {
		assert := assert.New(t)

		svc := newTestService(nil, nil, nil)
		err := svc.NewCategory(orderdomain.Category{Name: ""})

		assert.True(errors.Is(err, repo.ErrCategoryNameEmpty))
	})

	t.Run("creates valid category", func(t *testing.T) {
		assert := assert.New(t)

		catalogMock := &mockdb.Catalog{
			CreateCategoryFunc: func(_ orderdomain.Category) (int, error) { return 1, nil },
		}

		svc := newTestService(nil, nil, catalogMock)
		err := svc.NewCategory(orderdomain.Category{Name: "Drinks"})

		assert.NoError(err)
	})
}

func TestNewItem(t *testing.T) {
	t.Parallel()

	t.Run("rejects empty name", func(t *testing.T) {
		assert := assert.New(t)

		svc := newTestService(nil, nil, nil)
		err := svc.NewItem(orderdomain.Item{Name: ""})

		assert.True(errors.Is(err, repo.ErrItemNameEmpty))
	})

	t.Run("creates valid item", func(t *testing.T) {
		assert := assert.New(t)

		catalogMock := &mockdb.Catalog{
			CreateItemFunc: func(_ orderdomain.Item) (int, error) { return 1, nil },
		}

		svc := newTestService(nil, nil, catalogMock)
		err := svc.NewItem(orderdomain.Item{Name: "Beer"})

		assert.NoError(err)
	})
}
