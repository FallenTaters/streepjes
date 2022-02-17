package catalog

import (
	"sort"
	"strings"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/vugu/vugu"
)

type Categories struct {
	Categories         []domain.Category     `vugu:"data"`
	SelectedCategoryID int                   `vugu:"data"`
	OnClick            func(domain.Category) `vugu:"data"`
}

func (c *Categories) Compute(ctx vugu.ComputeCtx) {
	sort.Slice(c.Categories, func(i, j int) bool {
		return strings.Compare(c.Categories[i].Name, c.Categories[j].Name) < 0
	})
}

func (c *Categories) classes(category domain.Category) string {
	classes := `responsive extra small-margin`

	if c.SelectedCategoryID == category.ID {
		return classes + ` secondary`
	}

	return classes
}
