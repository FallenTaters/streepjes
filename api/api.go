/* api package contains payload types to be used by backend and frontend */
package api

import (
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

type Catalog struct {
	Categories []orderdomain.Category `json:"categories"`
	Items      []orderdomain.Item     `json:"items"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Role  authdomain.Role `json:"role"`
	Token string          `json:"token"`
}
