package catalog

import (
	"sort"
	"strings"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/vugu/vugu"
)

type Items struct {
	Items          []domain.Item
	SelectedItemID int
	OnClick        func(domain.Item)
}

func (i *Items) Compute(ctx vugu.ComputeCtx) {
	sort.Slice(i.Items, func(x, y int) bool {
		return strings.Compare(i.Items[x].Name, i.Items[y].Name) < 0
	})
}

func (i *Items) classes(item domain.Item) string {
	if i.SelectedItemID == item.ID {
		return `secondary`
	}

	return ``
}
