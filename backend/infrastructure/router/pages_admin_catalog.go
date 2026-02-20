package router

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"go.uber.org/zap"
)

type catalogData struct {
	pageData
	Categories         []orderdomain.Category
	DisplayItems       []orderdomain.Item
	SelectedCategory   *orderdomain.Category
	SelectedCategoryID int
	SelectedItemID     int
	FormTitle          string

	ShowCategoryForm bool
	EditCategory     *orderdomain.Category
	CategoryName     string

	ShowItemForm    bool
	EditItem        *orderdomain.Item
	ItemCategoryID  int
	ItemName        string
	PriceGladiators int
	PriceParabool   int
	PriceCalamari   int

	Error string
}

func (s *Server) getCatalogPage(w http.ResponseWriter, r *http.Request) {
	catalog, err := s.order.GetCatalog()
	if err != nil {
		s.internalError(w, "get catalog", err)
		return
	}

	categories := catalog.Categories
	items := catalog.Items

	sort.Slice(categories, func(i, j int) bool {
		return strings.ToLower(categories[i].Name) < strings.ToLower(categories[j].Name)
	})
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	data := catalogData{
		pageData:   newPageData(r, "catalog"),
		Categories: categories,
	}

	catID, _ := strconv.Atoi(r.URL.Query().Get("cat"))
	if catID != 0 {
		for i := range categories {
			if categories[i].ID == catID {
				c := categories[i]
				data.SelectedCategory = &c
				data.SelectedCategoryID = c.ID
				data.FormTitle = "Edit Category - " + c.Name
				data.ShowCategoryForm = true
				data.EditCategory = &c
				data.CategoryName = c.Name
				break
			}
		}

		displayItems := make([]orderdomain.Item, 0)
		for _, item := range items {
			if item.CategoryID == catID {
				displayItems = append(displayItems, item)
			}
		}
		data.DisplayItems = displayItems
	}

	action := r.URL.Query().Get("action")
	errMsg := r.URL.Query().Get("error")

	switch action {
	case "new-category":
		data.ShowCategoryForm = true
		data.ShowItemForm = false
		data.EditCategory = nil
		data.FormTitle = "New Category"
		data.CategoryName = ""
	case "edit-item":
		itemID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		for i := range items {
			if items[i].ID == itemID {
				it := items[i]
				data.ShowCategoryForm = false
				data.ShowItemForm = true
				data.EditItem = &it
				data.SelectedItemID = it.ID
				data.FormTitle = "Edit Item - " + it.Name
				data.ItemCategoryID = it.CategoryID
				data.ItemName = it.Name
				data.PriceGladiators = int(it.PriceGladiators)
				data.PriceParabool = int(it.PriceParabool)
				data.PriceCalamari = int(it.PriceCalamari)
				break
			}
		}
	case "new-item":
		data.ShowCategoryForm = false
		data.ShowItemForm = true
		data.FormTitle = "New Item"
		if data.SelectedCategory != nil {
			data.ItemCategoryID = data.SelectedCategory.ID
		}
	}

	if errMsg != "" {
		data.Error = errMsg
	}

	s.render(w, "admin/catalog.html", data)
}

func catalogRedirect(catID int, errMsg string) string {
	u := "/admin/catalog"
	if catID != 0 {
		u += fmt.Sprintf("?cat=%d", catID)
	}
	if errMsg != "" {
		sep := "?"
		if catID != 0 {
			sep = "&"
		}
		u += sep + "error=" + errMsg
	}
	return u
}

func (s *Server) postCatalogCategoryPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, catalogRedirect(0, "Invalid+form+data"), http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	name := r.FormValue("name")

	if idStr == "" {
		cat := orderdomain.Category{Name: name}
		if err := s.order.NewCategory(cat); err != nil {
			s.logger.Warn("category create failed", zap.Error(err))
			http.Redirect(w, r, catalogRedirect(0, "Unable+to+create+category."), http.StatusSeeOther)
			return
		}
	} else {
		id, _ := strconv.Atoi(idStr)
		cat := orderdomain.Category{ID: id, Name: name}
		if err := s.order.UpdateCategory(cat); err != nil {
			s.logger.Warn("category update failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, catalogRedirect(id, "Unable+to+update+category."), http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
}

func (s *Server) postDeleteCategoryPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := s.order.DeleteCategory(id); err != nil {
		s.logger.Warn("category delete failed", zap.Int("id", id), zap.Error(err))
		http.Redirect(w, r, catalogRedirect(id, "Unable+to+delete+category.+It+may+still+have+items."), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/catalog", http.StatusSeeOther)
}

func (s *Server) postCatalogItemPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, catalogRedirect(0, "Invalid+form+data"), http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	catID, _ := strconv.Atoi(r.FormValue("category_id"))
	name := r.FormValue("name")
	priceGlad, _ := strconv.Atoi(r.FormValue("price_gladiators"))
	pricePara, _ := strconv.Atoi(r.FormValue("price_parabool"))
	priceCala, _ := strconv.Atoi(r.FormValue("price_calamari"))

	if idStr == "" {
		item := orderdomain.Item{
			CategoryID:      catID,
			Name:            name,
			PriceGladiators: orderdomain.Price(priceGlad),
			PriceParabool:   orderdomain.Price(pricePara),
			PriceCalamari:   orderdomain.Price(priceCala),
		}
		if err := s.order.NewItem(item); err != nil {
			s.logger.Warn("item create failed", zap.Error(err))
			http.Redirect(w, r, catalogRedirect(catID, "Unable+to+create+item."), http.StatusSeeOther)
			return
		}
	} else {
		id, _ := strconv.Atoi(idStr)
		item := orderdomain.Item{
			ID:              id,
			CategoryID:      catID,
			Name:            name,
			PriceGladiators: orderdomain.Price(priceGlad),
			PriceParabool:   orderdomain.Price(pricePara),
			PriceCalamari:   orderdomain.Price(priceCala),
		}
		if err := s.order.UpdateItem(item); err != nil {
			s.logger.Warn("item update failed", zap.Int("id", id), zap.Error(err))
			http.Redirect(w, r, catalogRedirect(catID, "Unable+to+update+item."), http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, catalogRedirect(catID, ""), http.StatusSeeOther)
}

func (s *Server) postDeleteItemPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	catIDStr := r.URL.Query().Get("cat")
	catID, _ := strconv.Atoi(catIDStr)

	if err := s.order.DeleteItem(id); err != nil {
		s.logger.Warn("item delete failed", zap.Int("id", id), zap.Error(err))
		http.Redirect(w, r, catalogRedirect(catID, "Unable+to+delete+item."), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, catalogRedirect(catID, ""), http.StatusSeeOther)
}
