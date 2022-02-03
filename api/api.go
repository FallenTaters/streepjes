/* api package contains payload types to be used by backend and frontend */
package api

import "github.com/PotatoesFall/vecty-test/domain"

type Catalog struct {
	Categories []domain.Category `json:"categories"`
	Items      []domain.Item     `json:"items"`
}
