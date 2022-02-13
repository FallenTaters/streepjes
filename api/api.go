/* api package contains payload types to be used by backend and frontend */
package api

import "github.com/PotatoesFall/vecty-test/domain"

type Catalog struct {
	Categories []domain.Category `json:"categories"`
	Items      []domain.Item     `json:"items"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
