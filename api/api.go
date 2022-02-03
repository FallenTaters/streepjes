/* api package contains payload types to be used by backend and frontend */
package api

import "github.com/PotatoesFall/vecty-test/src/domain"

// from domain package
type (
	Club domain.Club

	Category domain.Category
	Item     domain.Item

	Order domain.Order

	Member domain.Member
)

type Catalog struct {
	Categories []Category `json:"categories"`
	Items      []Item     `json:"items"`
}
