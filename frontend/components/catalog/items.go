package catalog

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/vugu/vugu"
)

type Items struct {
	Items          []orderdomain.Item
	SelectedItemID int
	OnClick        func(orderdomain.Item)
	OnClickNew     func()
	HidePrice      bool
}

func (i *Items) Compute(ctx vugu.ComputeCtx) {
	sort.Slice(i.Items, func(x, y int) bool {
		return strings.Compare(i.Items[x].Name, i.Items[y].Name) < 0
	})
}

func (i *Items) classes(item orderdomain.Item) string {
	if i.SelectedItemID == item.ID {
		return `secondary`
	}

	return ``
}
