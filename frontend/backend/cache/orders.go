package cache

import "github.com/FallenTaters/streepjes/domain/orderdomain"

type Cache[T any] struct{} // TODO, see below

var Orders = Cache[[]orderdomain.Order]

// TODO: make a generic solution for this. Perhaps a type Cache[T] or so that takes an addFunc
