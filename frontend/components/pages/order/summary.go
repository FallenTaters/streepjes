package order

import (
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Summary struct {
	Members []orderdomain.Member `vugu:"data"`

	ShowMemberModal bool   `vugu:"data"`
	MemberSearch    string `vugu:"data"`
	Loading         bool   `vugu:"data"`
	Error           bool   `vugu:"data"`

	ShowMemberPaymentModal bool `vugu:"data"`
	LoadingMember          bool `vugu:"data"`
	ErrorMember            bool `vugu:"data"`
}

func (s *Summary) total() string {
	return store.Order.CalculateTotal().String()
}

func (s *Summary) GetMembers() []orderdomain.Member {
	var members []orderdomain.Member

	for _, member := range s.Members {
		if store.Order.Club == member.Club && (s.MemberSearch == `` || contains(member.Name, s.MemberSearch)) {
			members = append(members, member)
		}
	}

	return members
}

func contains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func (s *Summary) Init() {
	s.Loading = true
	s.Error = false

	go func() {
		defer func() {
			defer global.LockAndRender()()
			s.Loading = false
		}()

		members, err := cache.Members()
		if err != nil {
			defer global.LockOnly()()
			s.Error = true
			return
		}

		defer global.LockOnly()()
		s.Members = members
	}()
}

// TODO order members by last order
// TODO autoselect member if pressing enter while typing (top-most ?)
func (s *Summary) ChooseMember() {
	s.ShowMemberModal = true

	// wait for input to render before focussing it
	go func() {
		defer global.LockOnly()()
		jscall.Focus(`memberSearchInput`)
	}()
}

func (s *Summary) DisableButtons() bool {
	return len(store.Order.Lines) == 0
}

func (s *Summary) selectMember(member orderdomain.Member) {
	store.Order.Member = member
	s.ShowMemberModal = false
	s.ShowMemberPaymentModal = true
}
