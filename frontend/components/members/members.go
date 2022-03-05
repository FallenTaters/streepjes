package members

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Members struct {
	Members []orderdomain.Member `vugu:"data"`
	OnClick func(orderdomain.Member)
}
