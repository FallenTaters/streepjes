package streepjes

import (
	"fmt"
	"time"

	"github.com/FallenTaters/streepjes/migratelegacydb/domain/members"
	"github.com/FallenTaters/streepjes/migratelegacydb/domain/orders"
)

func DeleteMember(id int) error {
	if orders.MemberHasUnpaidOrders(id) {
		return members.ErrUnpaidOrders
	}

	return members.ForceDelete(id)
}

func RecalculateMemberDebt() error {
	allMembers, err := members.GetAll()
	if err != nil {
		return err
	}

	month := time.Now().Month()

	for _, m := range allMembers {
		allOrders, err := orders.GetByMemberID(m.ID)
		if err != nil {
			return err
		}

		debt := 0
		for _, o := range allOrders {
			if o.OrderTime.Month() != month {
				continue
			}

			if o.Status == orders.OrderStatusCancelled {
				continue
			}

			debt += o.Price
		}

		m.Debt = debt
		err = members.Put(m)
		if err != nil {
			return err
		}

		fmt.Printf("new debt for %s: %.2f\n", m.Name, float64(m.Debt)/100)
	}

	return nil
}
