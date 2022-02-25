package catalog

import (
	"sort"
	"strings"

	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
	"github.com/vugu/vugu"
)

type Categories struct {
	Categories         []orderdomain.Category     `vugu:"data"`
	SelectedCategoryID int                        `vugu:"data"`
	OnClick            func(orderdomain.Category) `vugu:"data"`
}

func (c *Categories) Compute(ctx vugu.ComputeCtx) {
	sort.Slice(c.Categories, func(i, j int) bool {
		return strings.Compare(c.Categories[i].Name, c.Categories[j].Name) < 0
	})
}

func (c *Categories) classes(category orderdomain.Category) string {
	classes := `responsive extra small-margin`

	if c.SelectedCategoryID == category.ID {
		return classes + ` secondary`
	}

	return classes
}
