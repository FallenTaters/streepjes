package postgres_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/postgres"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

var testDB postgres.Queryable

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct pool: %s", err)
	}

	if err := pool.Client.Ping(); err != nil {
		log.Fatalf("could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=test",
			"POSTGRES_USER=test",
			"POSTGRES_DB=streepjes_test",
			"listen_addresses='*'",
		},
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	_ = resource.Expire(120)

	connStr := fmt.Sprintf("postgres://test:test@%s/streepjes_test?sslmode=disable",
		resource.GetHostPort("5432/tcp"))

	pool.MaxWait = 30 * time.Second
	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to postgres: %s", err)
	}

	testDB, err = postgres.OpenDB(connStr)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}

	if err := postgres.Migrate(testDB, zap.NewNop()); err != nil {
		log.Fatalf("migration failed: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func cleanup(t *testing.T) {
	t.Helper()
	for _, table := range []string{"orders", "items", "categories", "members", "users"} {
		if _, err := testDB.(interface {
			Exec(string, ...any) (sql.Result, error)
		}).Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("cleanup %s: %v", table, err)
		}
	}
}

func newLogger() *zap.Logger { return zap.NewNop() }

// --- User Repo Tests ---

func TestUserRepo_CreateAndGet(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	user := authdomain.User{
		Username:     "alice",
		PasswordHash: "hash123",
		Club:         domain.ClubGladiators,
		Name:         "Alice",
		Role:         authdomain.RoleBartender,
	}

	id, err := r.Create(user)
	assert.NoError(err)
	assert.True(id > 0)

	got, err := r.Get(id)
	assert.NoError(err)
	assert.Eq("alice", got.Username)
	assert.Eq("Alice", got.Name)
	assert.Eq(domain.ClubGladiators, got.Club)
	assert.Eq(authdomain.RoleBartender, got.Role)
}

func TestUserRepo_GetByUsername(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	user := authdomain.User{
		Username:     "bob",
		PasswordHash: "hash",
		Club:         domain.ClubParabool,
		Name:         "Bob",
		Role:         authdomain.RoleAdmin,
	}

	_, err := r.Create(user)
	assert.NoError(err)

	got, err := r.GetByUsername("bob")
	assert.NoError(err)
	assert.Eq("Bob", got.Name)

	_, err = r.GetByUsername("nonexistent")
	assert.True(errors.Is(err, repo.ErrUserNotFound))
}

func TestUserRepo_GetAll(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	users := []authdomain.User{
		{Username: "u1", PasswordHash: "h", Club: domain.ClubGladiators, Name: "User1", Role: authdomain.RoleBartender},
		{Username: "u2", PasswordHash: "h", Club: domain.ClubParabool, Name: "User2", Role: authdomain.RoleAdmin},
	}
	for _, u := range users {
		_, err := r.Create(u)
		assert.NoError(err)
	}

	all, err := r.GetAll()
	assert.NoError(err)
	assert.Eq(2, len(all))
}

func TestUserRepo_Update(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	id, err := r.Create(authdomain.User{
		Username:     "carol",
		PasswordHash: "hash",
		Club:         domain.ClubCalamari,
		Name:         "Carol",
		Role:         authdomain.RoleBartender,
	})
	assert.NoError(err)

	got, err := r.Get(id)
	assert.NoError(err)
	got.Name = "Carol Updated"

	assert.NoError(r.Update(got))

	got2, err := r.Get(id)
	assert.NoError(err)
	assert.Eq("Carol Updated", got2.Name)
}

func TestUserRepo_UpdateActivity(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	id, err := r.Create(authdomain.User{
		Username:     "dan",
		PasswordHash: "hash",
		Club:         domain.ClubGladiators,
		Name:         "Dan",
		Role:         authdomain.RoleBartender,
	})
	assert.NoError(err)

	got, err := r.Get(id)
	assert.NoError(err)
	got.AuthToken = "new-token"
	got.AuthTime = time.Now().UTC().Truncate(time.Microsecond)

	assert.NoError(r.UpdateActivity(got))

	got2, err := r.GetByToken("new-token")
	assert.NoError(err)
	assert.Eq(id, got2.ID)
}

func TestUserRepo_Create_DuplicateUsername(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	user := authdomain.User{
		Username:     "dup",
		PasswordHash: "hash",
		Club:         domain.ClubGladiators,
		Name:         "Dup1",
		Role:         authdomain.RoleBartender,
	}
	_, err := r.Create(user)
	assert.NoError(err)

	user.Name = "Dup2"
	_, err = r.Create(user)
	assert.True(errors.Is(err, repo.ErrUsernameTaken))
}

func TestUserRepo_Create_DuplicateName(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	user := authdomain.User{
		Username:     "a",
		PasswordHash: "hash",
		Club:         domain.ClubGladiators,
		Name:         "SameName",
		Role:         authdomain.RoleBartender,
	}
	_, err := r.Create(user)
	assert.NoError(err)

	user.Username = "b"
	_, err = r.Create(user)
	assert.True(errors.Is(err, repo.ErrUsernameTaken))
}

func TestUserRepo_Create_MissingFields(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	_, err := r.Create(authdomain.User{})
	assert.True(errors.Is(err, repo.ErrUserMissingFields))
}

func TestUserRepo_Delete(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	id, err := r.Create(authdomain.User{
		Username:     "del",
		PasswordHash: "hash",
		Club:         domain.ClubGladiators,
		Name:         "Del",
		Role:         authdomain.RoleBartender,
	})
	assert.NoError(err)

	assert.NoError(r.Delete(id))

	_, err = r.Get(id)
	assert.True(errors.Is(err, repo.ErrUserNotFound))
}

func TestUserRepo_DeleteNotFound(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewUserRepo(testDB, newLogger())

	err := r.Delete(99999)
	assert.True(errors.Is(err, repo.ErrUserNotFound))
}

// --- Member Repo Tests ---

func TestMemberRepo_CreateAndGet(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewMemberRepo(testDB, newLogger())

	id, err := r.Create(orderdomain.Member{
		Club: domain.ClubGladiators,
		Name: "Member1",
	})
	assert.NoError(err)
	assert.True(id > 0)

	got, err := r.Get(id)
	assert.NoError(err)
	assert.Eq("Member1", got.Name)
	assert.Eq(domain.ClubGladiators, got.Club)
}

func TestMemberRepo_GetAll(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewMemberRepo(testDB, newLogger())

	_, err := r.Create(orderdomain.Member{Club: domain.ClubGladiators, Name: "M1"})
	assert.NoError(err)
	_, err = r.Create(orderdomain.Member{Club: domain.ClubParabool, Name: "M2"})
	assert.NoError(err)

	all, err := r.GetAll()
	assert.NoError(err)
	assert.Eq(2, len(all))
}

func TestMemberRepo_Update(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewMemberRepo(testDB, newLogger())

	id, err := r.Create(orderdomain.Member{Club: domain.ClubGladiators, Name: "UpdM"})
	assert.NoError(err)

	got, err := r.Get(id)
	assert.NoError(err)
	got.Name = "UpdatedMember"

	assert.NoError(r.Update(got))

	got2, err := r.Get(id)
	assert.NoError(err)
	assert.Eq("UpdatedMember", got2.Name)
}

func TestMemberRepo_Delete(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewMemberRepo(testDB, newLogger())

	id, err := r.Create(orderdomain.Member{Club: domain.ClubGladiators, Name: "DelM"})
	assert.NoError(err)

	assert.NoError(r.Delete(id))

	_, err = r.Get(id)
	assert.True(errors.Is(err, repo.ErrMemberNotFound))
}

func TestMemberRepo_DeleteNotFound(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewMemberRepo(testDB, newLogger())

	err := r.Delete(99999)
	assert.True(errors.Is(err, repo.ErrMemberNotFound))
}

func TestMemberRepo_Create_MissingFields(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewMemberRepo(testDB, newLogger())

	_, err := r.Create(orderdomain.Member{})
	assert.True(errors.Is(err, repo.ErrMemberFieldsNotFilled))
}

// --- Catalog Repo Tests ---

func TestCatalogRepo_CreateCategory(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	id, err := r.CreateCategory(orderdomain.Category{Name: "Drinks"})
	assert.NoError(err)
	assert.True(id > 0)

	cats, err := r.GetCategories()
	assert.NoError(err)
	assert.Eq(1, len(cats))
	assert.Eq("Drinks", cats[0].Name)
}

func TestCatalogRepo_CreateCategory_Empty(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	_, err := r.CreateCategory(orderdomain.Category{})
	assert.True(errors.Is(err, repo.ErrCategoryNameEmpty))
}

func TestCatalogRepo_CreateCategory_Duplicate(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	_, err := r.CreateCategory(orderdomain.Category{Name: "Dup"})
	assert.NoError(err)

	_, err = r.CreateCategory(orderdomain.Category{Name: "Dup"})
	assert.True(errors.Is(err, repo.ErrCategoryNameTaken))
}

func TestCatalogRepo_UpdateCategory(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	id, err := r.CreateCategory(orderdomain.Category{Name: "OldName"})
	assert.NoError(err)

	assert.NoError(r.UpdateCategory(orderdomain.Category{ID: id, Name: "NewName"}))

	cats, err := r.GetCategories()
	assert.NoError(err)
	assert.Eq(1, len(cats))
	assert.Eq("NewName", cats[0].Name)
}

func TestCatalogRepo_DeleteCategory(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	id, err := r.CreateCategory(orderdomain.Category{Name: "ToDel"})
	assert.NoError(err)

	assert.NoError(r.DeleteCategory(id))

	cats, err := r.GetCategories()
	assert.NoError(err)
	assert.Eq(0, len(cats))
}

func TestCatalogRepo_DeleteCategory_NotFound(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	err := r.DeleteCategory(99999)
	assert.True(errors.Is(err, repo.ErrCategoryNotFound))
}

func TestCatalogRepo_DeleteCategory_HasItems(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	catID, err := r.CreateCategory(orderdomain.Category{Name: "WithItems"})
	assert.NoError(err)

	_, err = r.CreateItem(orderdomain.Item{CategoryID: catID, Name: "Beer"})
	assert.NoError(err)

	err = r.DeleteCategory(catID)
	assert.True(errors.Is(err, repo.ErrCategoryHasItems))
}

func TestCatalogRepo_CreateItem(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	catID, err := r.CreateCategory(orderdomain.Category{Name: "Beers"})
	assert.NoError(err)

	item := orderdomain.Item{
		CategoryID:      catID,
		Name:            "Pilsner",
		PriceGladiators: 200,
		PriceParabool:   250,
		PriceCalamari:   300,
	}
	id, err := r.CreateItem(item)
	assert.NoError(err)
	assert.True(id > 0)

	items, err := r.GetItems()
	assert.NoError(err)
	assert.Eq(1, len(items))
	assert.Eq("Pilsner", items[0].Name)
	assert.Eq(orderdomain.Price(200), items[0].PriceGladiators)
	assert.Eq(orderdomain.Price(300), items[0].PriceCalamari)
}

func TestCatalogRepo_CreateItem_EmptyName(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	_, err := r.CreateItem(orderdomain.Item{})
	assert.True(errors.Is(err, repo.ErrItemNameEmpty))
}

func TestCatalogRepo_CreateItem_DuplicateName(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	catID, err := r.CreateCategory(orderdomain.Category{Name: "Cat"})
	assert.NoError(err)

	_, err = r.CreateItem(orderdomain.Item{CategoryID: catID, Name: "Same"})
	assert.NoError(err)

	_, err = r.CreateItem(orderdomain.Item{CategoryID: catID, Name: "Same"})
	assert.True(errors.Is(err, repo.ErrItemNameTaken))
}

func TestCatalogRepo_UpdateItem(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	catID, err := r.CreateCategory(orderdomain.Category{Name: "Food"})
	assert.NoError(err)

	itemID, err := r.CreateItem(orderdomain.Item{CategoryID: catID, Name: "Chips", PriceGladiators: 100})
	assert.NoError(err)

	assert.NoError(r.UpdateItem(orderdomain.Item{ID: itemID, CategoryID: catID, Name: "Crisps", PriceGladiators: 150}))

	items, err := r.GetItems()
	assert.NoError(err)
	assert.Eq("Crisps", items[0].Name)
	assert.Eq(orderdomain.Price(150), items[0].PriceGladiators)
}

func TestCatalogRepo_DeleteItem(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	catID, err := r.CreateCategory(orderdomain.Category{Name: "Snacks"})
	assert.NoError(err)

	itemID, err := r.CreateItem(orderdomain.Item{CategoryID: catID, Name: "Nuts"})
	assert.NoError(err)

	assert.NoError(r.DeleteItem(itemID))

	items, err := r.GetItems()
	assert.NoError(err)
	assert.Eq(0, len(items))
}

func TestCatalogRepo_DeleteItem_NotFound(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewCatalogRepo(testDB, newLogger())

	err := r.DeleteItem(99999)
	assert.True(errors.Is(err, repo.ErrItemNotFound))
}

// --- Order Repo Tests ---

func createTestUser(t *testing.T) int {
	t.Helper()
	r := postgres.NewUserRepo(testDB, newLogger())
	id, err := r.Create(authdomain.User{
		Username:     fmt.Sprintf("bartender_%d", time.Now().UnixNano()),
		PasswordHash: "hash",
		Club:         domain.ClubGladiators,
		Name:         fmt.Sprintf("Bartender_%d", time.Now().UnixNano()),
		Role:         authdomain.RoleBartender,
	})
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return id
}

func createTestMember(t *testing.T) int {
	t.Helper()
	r := postgres.NewMemberRepo(testDB, newLogger())
	id, err := r.Create(orderdomain.Member{
		Club: domain.ClubGladiators,
		Name: fmt.Sprintf("Member_%d", time.Now().UnixNano()),
	})
	if err != nil {
		t.Fatalf("create test member: %v", err)
	}
	return id
}

func TestOrderRepo_CreateAndGet(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	bartenderID := createTestUser(t)
	memberID := createTestMember(t)

	r := postgres.NewOrderRepo(testDB, newLogger())
	now := time.Now().UTC().Truncate(time.Microsecond)

	order := orderdomain.Order{
		Club:        domain.ClubGladiators,
		BartenderID: bartenderID,
		MemberID:    memberID,
		Contents:    `[{"name":"Beer","amount":2}]`,
		Price:       400,
		OrderTime:   now,
		Status:      orderdomain.StatusOpen,
		StatusTime:  now,
	}

	id, err := r.Create(order)
	assert.NoError(err)
	assert.True(id > 0)

	got, err := r.Get(id)
	assert.NoError(err)
	assert.Eq(domain.ClubGladiators, got.Club)
	assert.Eq(bartenderID, got.BartenderID)
	assert.Eq(memberID, got.MemberID)
	assert.Eq(orderdomain.Price(400), got.Price)
	assert.Eq(orderdomain.StatusOpen, got.Status)
}

func TestOrderRepo_Create_MissingFields(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewOrderRepo(testDB, newLogger())

	_, err := r.Create(orderdomain.Order{})
	assert.True(errors.Is(err, repo.ErrOrderFieldsNotFilled))
}

func TestOrderRepo_Create_NonexistentBartender(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewOrderRepo(testDB, newLogger())

	_, err := r.Create(orderdomain.Order{
		Club:        domain.ClubGladiators,
		BartenderID: 99999,
		OrderTime:   time.Now(),
		Status:      orderdomain.StatusOpen,
		StatusTime:  time.Now(),
	})
	assert.True(errors.Is(err, repo.ErrUserNotFound))
}

func TestOrderRepo_Create_NonexistentMember(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	bartenderID := createTestUser(t)
	r := postgres.NewOrderRepo(testDB, newLogger())

	_, err := r.Create(orderdomain.Order{
		Club:        domain.ClubGladiators,
		BartenderID: bartenderID,
		MemberID:    99999,
		OrderTime:   time.Now(),
		Status:      orderdomain.StatusOpen,
		StatusTime:  time.Now(),
	})
	assert.True(errors.Is(err, repo.ErrMemberNotFound))
}

func TestOrderRepo_Filter(t *testing.T) {
	cleanup(t)
	a := assert.New(t)

	bartenderID := createTestUser(t)
	memberID := createTestMember(t)
	r := postgres.NewOrderRepo(testDB, newLogger())

	now := time.Now().UTC().Truncate(time.Microsecond)
	for i := range 5 {
		_, err := r.Create(orderdomain.Order{
			Club:        domain.ClubGladiators,
			BartenderID: bartenderID,
			MemberID:    memberID,
			Price:       orderdomain.Price(100 * (i + 1)),
			OrderTime:   now.Add(time.Duration(i) * time.Hour),
			Status:      orderdomain.StatusOpen,
			StatusTime:  now,
		})
		a.NoError(err)
	}

	t.Run("filter by club", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{Club: domain.ClubGladiators})
		assert.NoError(err)
		assert.Eq(5, len(orders))
	})

	t.Run("filter by bartender", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{BartenderID: bartenderID})
		assert.NoError(err)
		assert.Eq(5, len(orders))
	})

	t.Run("filter with limit", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{Club: domain.ClubGladiators, Limit: 3})
		assert.NoError(err)
		assert.Eq(3, len(orders))
	})

	t.Run("filter by time range", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{
			Start: now.Add(1 * time.Hour),
			End:   now.Add(4 * time.Hour),
		})
		assert.NoError(err)
		assert.Eq(3, len(orders))
	})

	t.Run("filter by member", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{MemberID: memberID})
		assert.NoError(err)
		assert.Eq(5, len(orders))
	})

	t.Run("filter excludes status", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{
			Club:      domain.ClubGladiators,
			StatusNot: []orderdomain.Status{orderdomain.StatusCancelled},
		})
		assert.NoError(err)
		assert.Eq(5, len(orders))
	})

	t.Run("empty result for other club", func(t *testing.T) {
		assert := assert.New(t)
		orders, err := r.Filter(repo.OrderFilter{Club: domain.ClubParabool})
		assert.NoError(err)
		assert.Eq(0, len(orders))
	})
}

func TestOrderRepo_Delete(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	bartenderID := createTestUser(t)
	r := postgres.NewOrderRepo(testDB, newLogger())

	now := time.Now().UTC().Truncate(time.Microsecond)
	id, err := r.Create(orderdomain.Order{
		Club:        domain.ClubGladiators,
		BartenderID: bartenderID,
		Price:       100,
		OrderTime:   now,
		Status:      orderdomain.StatusOpen,
		StatusTime:  now,
	})
	assert.NoError(err)

	assert.NoError(r.Delete(id))

	got, err := r.Get(id)
	assert.NoError(err)
	assert.Eq(orderdomain.StatusCancelled, got.Status)
}

func TestOrderRepo_DeleteNotFound(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewOrderRepo(testDB, newLogger())

	err := r.Delete(99999)
	assert.True(errors.Is(err, repo.ErrOrderNotFound))
}

func TestOrderRepo_GetNotFound(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	r := postgres.NewOrderRepo(testDB, newLogger())

	_, err := r.Get(99999)
	assert.True(errors.Is(err, repo.ErrOrderNotFound))
}

func TestOrderRepo_Create_WithoutMember(t *testing.T) {
	cleanup(t)
	assert := assert.New(t)

	bartenderID := createTestUser(t)
	r := postgres.NewOrderRepo(testDB, newLogger())

	now := time.Now().UTC().Truncate(time.Microsecond)
	id, err := r.Create(orderdomain.Order{
		Club:        domain.ClubGladiators,
		BartenderID: bartenderID,
		MemberID:    0,
		Price:       200,
		OrderTime:   now,
		Status:      orderdomain.StatusOpen,
		StatusTime:  now,
	})
	assert.NoError(err)

	got, err := r.Get(id)
	assert.NoError(err)
	assert.Eq(0, got.MemberID)
}

// --- Migration Tests ---

func TestMigrate_Idempotent(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(postgres.Migrate(testDB, newLogger()))
}
