/* api package contains payload types to be used by backend and frontend */
package api

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Catalog struct {
	Categories []orderdomain.Category `json:"categories"`
	Items      []orderdomain.Item     `json:"items"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Original string `json:"original"`
	New      string `json:"new"`
}

type MemberDetails struct {
	orderdomain.Member

	Debt orderdomain.Price `json:"debt"`
}
