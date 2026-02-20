package router_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/router"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

var (
	errTest     = errors.New("test error")
	errNotFound = errors.New("not found")
)

func newTestHandler(authSvc auth.Service, orderSvc order.Service) http.Handler {
	static := func(filename string) ([]byte, error) {
		if filename == "test.js" {
			return []byte("// test"), nil
		}
		return nil, errNotFound
	}
	return router.New(static, authSvc, orderSvc, false, zap.NewNop())
}

func TestGetVersion(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	handler := newTestHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Eq(http.StatusOK, rec.Code)
	assert.True(strings.Contains(rec.Body.String(), "Version:"))
}

func TestGetStatic(t *testing.T) {
	t.Parallel()

	t.Run("existing file", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/static/test.js", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusOK, rec.Code)
		assert.Eq("// test", rec.Body.String())
	})

	t.Run("missing file returns 404", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/static/missing.js", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusNotFound, rec.Code)
	})
}

func TestGetRoot(t *testing.T) {
	t.Parallel()

	t.Run("no cookie redirects to login", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(&authServiceMock{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/login", rec.Header().Get("Location"))
	})

	t.Run("invalid token redirects to login", func(t *testing.T) {
		assert := assert.New(t)

		authMock := &authServiceMock{
			CheckFunc: func(_ string) (authdomain.User, error) {
				return authdomain.User{}, auth.ErrInvalidToken
			},
		}
		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "bad"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/login", rec.Header().Get("Location"))
	})

	t.Run("admin redirects to billing", func(t *testing.T) {
		assert := assert.New(t)

		authMock := &authServiceMock{
			CheckFunc: func(_ string) (authdomain.User, error) {
				return authdomain.User{Role: authdomain.RoleAdmin}, nil
			},
		}
		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/admin/billing", rec.Header().Get("Location"))
	})

	t.Run("bartender redirects to order", func(t *testing.T) {
		assert := assert.New(t)

		authMock := &authServiceMock{
			CheckFunc: func(_ string) (authdomain.User, error) {
				return authdomain.User{Role: authdomain.RoleBartender}, nil
			},
		}
		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/order", rec.Header().Get("Location"))
	})
}

func TestPostLogin(t *testing.T) {
	t.Parallel()

	t.Run("successful login sets cookie", func(t *testing.T) {
		assert := assert.New(t)

		authMock := &authServiceMock{
			LoginFunc: func(_, _ string) (authdomain.User, error) {
				return authdomain.User{
					Role:      authdomain.RoleBartender,
					AuthToken: "test-token",
				}, nil
			},
		}

		handler := newTestHandler(authMock, nil)
		body := strings.NewReader("username=admin&password=pass")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)

		var foundCookie bool
		for _, c := range rec.Result().Cookies() {
			if c.Name == api.AuthTokenCookieName {
				assert.Eq("test-token", c.Value)
				assert.True(c.HttpOnly)
				foundCookie = true
			}
		}
		assert.True(foundCookie)
	})

	t.Run("bad credentials redirects with error", func(t *testing.T) {
		assert := assert.New(t)

		authMock := &authServiceMock{
			LoginFunc: func(_, _ string) (authdomain.User, error) {
				return authdomain.User{}, auth.ErrInvalidCredentials
			},
		}

		handler := newTestHandler(authMock, nil)
		body := strings.NewReader("username=bad&password=bad")
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/login?error=1", rec.Header().Get("Location"))
	})
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("unauthenticated request redirects to login", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(&authServiceMock{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/login", rec.Header().Get("Location"))
	})

	t.Run("bartender cannot access admin", func(t *testing.T) {
		assert := assert.New(t)

		authMock := &authServiceMock{
			CheckFunc: func(_ string) (authdomain.User, error) {
				return authdomain.User{Role: authdomain.RoleBartender}, nil
			},
		}
		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodGet, "/admin/billing", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusSeeOther, rec.Code)
		assert.Eq("/", rec.Header().Get("Location"))
	})
}

func TestGetLogout(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var logoutCalled bool
	authMock := &authServiceMock{
		CheckFunc: func(_ string) (authdomain.User, error) {
			return authdomain.User{ID: 1, Role: authdomain.RoleBartender}, nil
		},
		LogoutFunc: func(_ int) error {
			logoutCalled = true
			return nil
		},
	}

	handler := newTestHandler(authMock, nil)
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Eq(http.StatusSeeOther, rec.Code)
	assert.Eq("/login", rec.Header().Get("Location"))
	assert.True(logoutCalled)
}

func TestPostActive(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var activeCalled bool
	authMock := &authServiceMock{
		CheckFunc: func(_ string) (authdomain.User, error) {
			return authdomain.User{ID: 1, Role: authdomain.RoleBartender}, nil
		},
		ActiveFunc: func(_ int) error {
			activeCalled = true
			return nil
		},
	}

	handler := newTestHandler(authMock, nil)
	req := httptest.NewRequest(http.MethodPost, "/active", nil)
	req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Eq(http.StatusNoContent, rec.Code)
	assert.True(activeCalled)
}

func TestGetMember(t *testing.T) {
	t.Parallel()

	authMock := &authServiceMock{
		CheckFunc: func(_ string) (authdomain.User, error) {
			return authdomain.User{ID: 1, Role: authdomain.RoleBartender}, nil
		},
	}

	t.Run("returns member JSON", func(t *testing.T) {
		assert := assert.New(t)

		orderMock := &orderServiceMock{
			GetMemberDetailsFunc: func(id int) (api.MemberDetails, error) {
				return api.MemberDetails{
					Member: orderdomain.Member{ID: id, Name: "Alice"},
					Debt:   500,
				}, nil
			},
		}

		handler := newTestHandler(authMock, orderMock)
		req := httptest.NewRequest(http.MethodGet, "/api/member/1", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusOK, rec.Code)
		assert.Eq("application/json", rec.Header().Get("Content-Type"))

		var result api.MemberDetails
		assert.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
		assert.Eq("Alice", result.Name)
	})

	t.Run("bad id returns 400", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/member/abc", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusBadRequest, rec.Code)
	})

	t.Run("member not found returns 404", func(t *testing.T) {
		assert := assert.New(t)

		orderMock := &orderServiceMock{
			GetMemberDetailsFunc: func(_ int) (api.MemberDetails, error) {
				return api.MemberDetails{}, repo.ErrMemberNotFound
			},
		}

		handler := newTestHandler(authMock, orderMock)
		req := httptest.NewRequest(http.MethodGet, "/api/member/99", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusNotFound, rec.Code)
	})
}

func TestPostOrder(t *testing.T) {
	t.Parallel()

	authMock := &authServiceMock{
		CheckFunc: func(_ string) (authdomain.User, error) {
			return authdomain.User{ID: 1, Role: authdomain.RoleBartender, Club: domain.ClubGladiators}, nil
		},
	}

	t.Run("valid order", func(t *testing.T) {
		assert := assert.New(t)

		var placedOrder orderdomain.Order
		orderMock := &orderServiceMock{
			PlaceOrderFunc: func(o orderdomain.Order, _ authdomain.User) error {
				placedOrder = o
				return nil
			},
		}

		handler := newTestHandler(authMock, orderMock)
		body := `{"club":"Gladiators","memberId":1,"price":300}`
		req := httptest.NewRequest(http.MethodPost, "/api/order", strings.NewReader(body))
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusOK, rec.Code)
		assert.Eq(domain.ClubGladiators, placedOrder.Club)
	})

	t.Run("invalid body returns 400", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodPost, "/api/order", strings.NewReader("not json"))
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusBadRequest, rec.Code)
	})
}

func TestGetDownload(t *testing.T) {
	t.Parallel()

	authMock := &authServiceMock{
		CheckFunc: func(_ string) (authdomain.User, error) {
			return authdomain.User{ID: 1, Role: authdomain.RoleAdmin, Club: domain.ClubGladiators}, nil
		},
	}

	t.Run("returns CSV", func(t *testing.T) {
		assert := assert.New(t)

		orderMock := &orderServiceMock{
			BillingCSVFunc: func(_ domain.Club, _ orderdomain.Month) ([]byte, error) {
				return []byte("Member,Total\n"), nil
			},
		}

		handler := newTestHandler(authMock, orderMock)
		req := httptest.NewRequest(http.MethodGet, "/admin/download?month=2025-06", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusOK, rec.Code)
		assert.Eq("text/csv", rec.Header().Get("Content-Type"))
		assert.True(strings.Contains(rec.Header().Get("Content-Disposition"), ".csv"))
	})

	t.Run("bad month returns 400", func(t *testing.T) {
		assert := assert.New(t)

		handler := newTestHandler(authMock, nil)
		req := httptest.NewRequest(http.MethodGet, "/admin/download?month=bad", nil)
		req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Eq(http.StatusBadRequest, rec.Code)
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	authMock := &authServiceMock{
		CheckFunc: func(_ string) (authdomain.User, error) {
			panic("test panic")
		},
	}

	handler := newTestHandler(authMock, nil)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: "valid"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Eq(http.StatusInternalServerError, rec.Code)
}

// --- mock implementations ---

// authServiceMock implements auth.Service for testing.
type authServiceMock struct {
	LoginFunc          func(user, pass string) (authdomain.User, error)
	CheckFunc          func(token string) (authdomain.User, error)
	ActiveFunc         func(id int) error
	LogoutFunc         func(id int) error
	RegisterFunc       func(user authdomain.User, password string) error
	UpdateFunc         func(user authdomain.User, password string) error
	ChangePasswordFunc func(user authdomain.User, cp api.ChangePassword) error
	ChangeNameFunc     func(user authdomain.User, name string) error
	GetUsersFunc       func() ([]authdomain.User, error)
	DeleteFunc         func(id int) error
}

func (m *authServiceMock) Login(user, pass string) (authdomain.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(user, pass)
	}
	return authdomain.User{}, auth.ErrInvalidCredentials
}

func (m *authServiceMock) Check(token string) (authdomain.User, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(token)
	}
	return authdomain.User{}, auth.ErrInvalidToken
}

func (m *authServiceMock) Active(id int) error {
	if m.ActiveFunc != nil {
		return m.ActiveFunc(id)
	}
	return nil
}

func (m *authServiceMock) Logout(id int) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(id)
	}
	return nil
}

func (m *authServiceMock) Register(user authdomain.User, password string) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(user, password)
	}
	return nil
}

func (m *authServiceMock) Update(user authdomain.User, password string) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user, password)
	}
	return nil
}

func (m *authServiceMock) ChangePassword(user authdomain.User, cp api.ChangePassword) error {
	if m.ChangePasswordFunc != nil {
		return m.ChangePasswordFunc(user, cp)
	}
	return nil
}

func (m *authServiceMock) ChangeName(user authdomain.User, name string) error {
	if m.ChangeNameFunc != nil {
		return m.ChangeNameFunc(user, name)
	}
	return nil
}

func (m *authServiceMock) GetUsers() ([]authdomain.User, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc()
	}
	return nil, nil
}

func (m *authServiceMock) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// orderServiceMock implements order.Service for testing.
type orderServiceMock struct {
	GetAllMembersFunc         func() ([]orderdomain.Member, error)
	GetMemberDetailsFunc      func(id int) (api.MemberDetails, error)
	GetCatalogFunc            func() (api.Catalog, error)
	GetOrdersForBartenderFunc func(id int) ([]orderdomain.Order, error)
	GetOrdersByClubFunc       func(club domain.Club, month orderdomain.Month) ([]orderdomain.Order, error)
	BillingCSVFunc            func(club domain.Club, month orderdomain.Month) ([]byte, error)
	PlaceOrderFunc            func(order orderdomain.Order, bartender authdomain.User) error
	BartenderDeleteOrderFunc  func(bartenderID, orderID int) error
	NewCategoryFunc           func(orderdomain.Category) error
	UpdateCategoryFunc        func(orderdomain.Category) error
	DeleteCategoryFunc        func(id int) error
	NewItemFunc               func(orderdomain.Item) error
	UpdateItemFunc            func(orderdomain.Item) error
	DeleteItemFunc            func(id int) error
	GetLeaderboardFunc        func(api.LeaderboardFilter) (api.Leaderboard, error)
	NewMemberFunc             func(orderdomain.Member) error
	EditMemberFunc            func(orderdomain.Member) error
	DeleteMemberFunc          func(id int) error
}

func (m *orderServiceMock) GetAllMembers() ([]orderdomain.Member, error) {
	if m.GetAllMembersFunc != nil {
		return m.GetAllMembersFunc()
	}
	return nil, nil
}

func (m *orderServiceMock) GetMemberDetails(id int) (api.MemberDetails, error) {
	if m.GetMemberDetailsFunc != nil {
		return m.GetMemberDetailsFunc(id)
	}
	return api.MemberDetails{}, errTest
}

func (m *orderServiceMock) GetCatalog() (api.Catalog, error) {
	if m.GetCatalogFunc != nil {
		return m.GetCatalogFunc()
	}
	return api.Catalog{}, nil
}

func (m *orderServiceMock) GetOrdersForBartender(id int) ([]orderdomain.Order, error) {
	if m.GetOrdersForBartenderFunc != nil {
		return m.GetOrdersForBartenderFunc(id)
	}
	return nil, nil
}

func (m *orderServiceMock) GetOrdersByClub(club domain.Club, month orderdomain.Month) ([]orderdomain.Order, error) {
	if m.GetOrdersByClubFunc != nil {
		return m.GetOrdersByClubFunc(club, month)
	}
	return nil, nil
}

func (m *orderServiceMock) BillingCSV(club domain.Club, month orderdomain.Month) ([]byte, error) {
	if m.BillingCSVFunc != nil {
		return m.BillingCSVFunc(club, month)
	}
	return nil, nil
}

func (m *orderServiceMock) PlaceOrder(o orderdomain.Order, u authdomain.User) error {
	if m.PlaceOrderFunc != nil {
		return m.PlaceOrderFunc(o, u)
	}
	return nil
}

func (m *orderServiceMock) BartenderDeleteOrder(bartenderID, orderID int) error {
	if m.BartenderDeleteOrderFunc != nil {
		return m.BartenderDeleteOrderFunc(bartenderID, orderID)
	}
	return nil
}

func (m *orderServiceMock) NewCategory(c orderdomain.Category) error {
	if m.NewCategoryFunc != nil {
		return m.NewCategoryFunc(c)
	}
	return nil
}

func (m *orderServiceMock) UpdateCategory(c orderdomain.Category) error {
	if m.UpdateCategoryFunc != nil {
		return m.UpdateCategoryFunc(c)
	}
	return nil
}

func (m *orderServiceMock) DeleteCategory(id int) error {
	if m.DeleteCategoryFunc != nil {
		return m.DeleteCategoryFunc(id)
	}
	return nil
}

func (m *orderServiceMock) NewItem(i orderdomain.Item) error {
	if m.NewItemFunc != nil {
		return m.NewItemFunc(i)
	}
	return nil
}

func (m *orderServiceMock) UpdateItem(i orderdomain.Item) error {
	if m.UpdateItemFunc != nil {
		return m.UpdateItemFunc(i)
	}
	return nil
}

func (m *orderServiceMock) DeleteItem(id int) error {
	if m.DeleteItemFunc != nil {
		return m.DeleteItemFunc(id)
	}
	return nil
}

func (m *orderServiceMock) GetLeaderboard(f api.LeaderboardFilter) (api.Leaderboard, error) {
	if m.GetLeaderboardFunc != nil {
		return m.GetLeaderboardFunc(f)
	}
	return api.Leaderboard{}, nil
}

func (m *orderServiceMock) NewMember(member orderdomain.Member) error {
	if m.NewMemberFunc != nil {
		return m.NewMemberFunc(member)
	}
	return nil
}

func (m *orderServiceMock) EditMember(member orderdomain.Member) error {
	if m.EditMemberFunc != nil {
		return m.EditMemberFunc(member)
	}
	return nil
}

func (m *orderServiceMock) DeleteMember(id int) error {
	if m.DeleteMemberFunc != nil {
		return m.DeleteMemberFunc(id)
	}
	return nil
}
