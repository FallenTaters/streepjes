/*
This package is used solely for migrating from the old streepjes-api bbolt repo.
It is entirely self contained except for when it writes to the new sqlite db.
If the bbolt repo is not used anywhere and there are successful production backups of the new db,
and the code has been tested properly, this entire package and its subdirectories can be removed.
*/
package main

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/domain/catalog"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/domain/members"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/domain/orders"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/domain/users"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared/buckets"
)

func main() {
	defer buckets.Init()() //nolint:errcheck

	getLegacyData()

	fmt.Println(`categories:`, len(legacyData.Catalog.Categories))
	fmt.Println(`products:`, len(legacyData.Catalog.Products))
	fmt.Println(`members:`, len(legacyData.Members))
	fmt.Println(`orders:`, len(legacyData.Orders))
	fmt.Println(`users:`, len(legacyData.Users))

	migrateStructs()

	persist()
}

var legacyData struct {
	Catalog catalog.Catalog
	Members []members.Member
	Orders  []orders.Order
	Users   []users.User
}

func getLegacyData() {
	legacyData.Catalog, _ = catalog.Get()
	legacyData.Members, _ = members.GetAll()
	legacyData.Orders, _ = orders.GetAll()
	legacyData.Users, _ = users.GetAll()
}

var newData struct {
	Categories []orderdomain.Category
	Items      []orderdomain.Item
	Members    []orderdomain.Member
	Users      []authdomain.User
	Orders     []orderdomain.Order
}

func migrateStructs() { //nolint:funlen
	for _, category := range legacyData.Catalog.Categories {
		newData.Categories = append(newData.Categories, orderdomain.Category{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	for _, item := range legacyData.Catalog.Products {
		newData.Items = append(newData.Items, orderdomain.Item{
			ID:              item.ID,
			CategoryID:      item.CategoryID,
			Name:            item.Name,
			PriceGladiators: orderdomain.Price(item.PriceGladiators),
			PriceParabool:   orderdomain.Price(item.PriceParabool),
		})
	}

	for _, member := range legacyData.Members {
		newData.Members = append(newData.Members, orderdomain.Member{
			ID:   member.ID,
			Club: domain.Club(member.Club),
			Name: member.Name,
			Debt: member.Debt,
		})
	}

	for i, user := range legacyData.Users {
		newData.Users = append(newData.Users, authdomain.User{ //nolint:exhaustivestruct
			ID:           i + 1, // EYYYY
			Username:     user.Username,
			PasswordHash: string(user.Password),
			Club:         domain.Club(user.Club),
			Name:         user.Name,
			Role:         authdomain.Role(user.Role), // just so happens to align, risky
		})
	}

	for _, order := range legacyData.Orders {
		var bartenderID int
		for _, bartender := range newData.Users {
			if bartender.Username == order.Bartender {
				bartenderID = bartender.ID
			}
		}

		newData.Orders = append(newData.Orders, orderdomain.Order{
			ID:          order.ID,
			Club:        domain.Club(order.Club), // just so happens to align, risky
			BartenderID: bartenderID,
			MemberID:    order.MemberID,
			Contents:    order.Contents,
			Price:       orderdomain.Price(order.Price),
			OrderTime:   order.OrderTime,
			Status:      orderdomain.Status(order.Status), // just so happens to align, risky
			StatusTime:  order.StatusTime,
		})
	}
}

func persist() {
	// TODO
}
