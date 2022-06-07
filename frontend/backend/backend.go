package backend

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrStatus       = errors.New(`received unexpected status code`)
	ErrUnauthorized = errors.New(`request not authorized`)
	ErrForbidden    = errors.New(`request forbidden`)
)

func Init(endpoint *url.URL) {
	settings.Endpoint = endpoint
}

type Settings struct {
	Endpoint *url.URL
}

var settings Settings

func (s Settings) URL() string {
	return s.Endpoint.String()
}

func GetCatalog() (api.Catalog, error) {
	var catalog api.Catalog
	return catalog, get(`/catalog`, &catalog)
}

func PostNewCategory(category orderdomain.Category) error {
	return post(`/admin/category/new`, category, nil)
}

func PostUpdateCategory(category orderdomain.Category) error {
	return post(`/admin/category/update`, category, nil)
}

func PostDeleteCategory(id int) error {
	return post(`/admin/category/`+strconv.Itoa(id)+`/delete`, nil, nil)
}

func PostNewItem(item orderdomain.Item) error {
	return post(`/admin/item/new`, item, nil)
}

func PostUpdateItem(item orderdomain.Item) error {
	return post(`/admin/item/update`, item, nil)
}

func PostDeleteItem(id int) error {
	return post(`/admin/item/`+strconv.Itoa(id)+`/delete`, nil, nil)
}

func GetMembers() ([]orderdomain.Member, error) {
	var members []orderdomain.Member
	return members, get(`/members`, &members)
}

func GetMember(id int) (api.MemberDetails, error) {
	var member api.MemberDetails
	return member, get(`/member/`+strconv.Itoa(id), &member)
}

func GetOrders() ([]orderdomain.Order, error) {
	var orders []orderdomain.Order
	return orders, get(`/orders`, &orders)
}

func PostLogout() error {
	resp, err := http.Post(settings.URL()+`/logout`, ``, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// don't check status because logout is often called when already logged out, just to ensure logout

	return nil
}

func PostLogin(req api.Credentials) (authdomain.User, error) {
	var user authdomain.User
	return user, post(`/login`, req, &user)
}

func PostActive() (authdomain.User, error) {
	var user authdomain.User
	return user, post(`/active`, nil, &user)
}

func PostChangePassword(changePassword api.ChangePassword) error {
	return post(`/me/password`, changePassword, nil)
}

func PostChangeName(name string) error {
	return post(`/me/name`, name, nil)
}

func PostOrder(order orderdomain.Order) error {
	return post(`/order`, order, nil)
}

func PostDeleteOrder(id int) error {
	return post(fmt.Sprintf(`/order/%d/delete`, id), nil, nil)
}

func GetUsers() ([]authdomain.User, error) {
	var users []authdomain.User
	return users, get(`/admin/users`, &users)
}

func PostNewUser(user api.UserWithPassword) error {
	return post(`/admin/users/new`, user, nil)
}

func PostEditUser(user api.UserWithPassword) error {
	return post(`/admin/users/edit`, user, nil)
}

func PostDeleteUser(id int) error {
	return post(fmt.Sprintf(`/admin/users/%d/delete`, id), nil, nil)
}
