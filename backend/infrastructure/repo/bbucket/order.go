package bbucket

import (
	"errors"
	"time"

	"github.com/FallenTaters/bbucket"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"

	"go.etcd.io/bbolt"
)

func NewOrderRepo(db *bbolt.DB) repo.Order {
	return orderRepo{
		bbucket.New(db, orderBucket),
		bbucket.New(db, memberBucket),
	}
}

type orderRepo struct {
	bucket       bbucket.Bucket
	memberBucket bbucket.Bucket
}

func (or orderRepo) Get(id int) (orderdomain.Order, bool) {
	var o orderdomain.Order

	err := or.bucket.Get(bbucket.Itob(id), &o)
	if errors.Is(err, bbucket.ErrObjectNotFound) {
		return orderdomain.Order{}, false
	}
	if err != nil {
		panic(err)
	}

	return o, true
}

func (or orderRepo) Filter(filter repo.OrderFilter) []orderdomain.Order {
	orders := []orderdomain.Order{}

	err := or.bucket.GetAll(&orderdomain.Order{}, func(ptr interface{}) error {
		o := *ptr.(*orderdomain.Order)

		if filter.Bartender != nil && o.Bartender != *filter.Bartender {
			return nil
		}

		orders = append(orders, o)

		return nil
	})
	if err != nil {
		panic(err)
	}

	return orders
}

func (or orderRepo) Create(order orderdomain.Order) error {
	order.ID = or.bucket.NextSequence()
	order.OrderTime = time.Now().Local()
	order.StatusTime = order.OrderTime

	err := or.bucket.Create(orderKey(order), order)
	if err != nil || order.MemberID == 0 {
		return err
	}

	return or.memberBucket.Update(bbucket.Itob(order.MemberID), &orderdomain.Member{}, func(ptr interface{}) (object interface{}, err error) {
		member := *ptr.(*orderdomain.Member)

		member.Debt += order.Price

		return member, nil
	})
}

func (or orderRepo) DeleteByID(id int) bool {
	var order orderdomain.Order

	err := or.bucket.Update(bbucket.Itob(id), &orderdomain.Order{}, func(ptr interface{}) (object interface{}, err error) {
		order = *ptr.(*orderdomain.Order)
		o := order

		o.Status = orderdomain.OrderStatusCancelled
		o.StatusTime = time.Now()

		return o, nil
	})
	if err != nil {
		return false
	}
	if order.MemberID == 0 || order.Status != orderdomain.OrderStatusOpen {
		return true
	}

	err = or.memberBucket.Update(bbucket.Itob(order.MemberID), &orderdomain.Member{}, func(ptr interface{}) (object interface{}, err error) {
		member := *ptr.(*orderdomain.Member)

		member.Debt -= order.Price

		return member, nil
	})
	if err != nil {
		panic(err)
	}

	return true
}

func orderKey(o orderdomain.Order) []byte {
	return itob(o.ID)
}
