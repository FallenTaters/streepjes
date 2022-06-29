package pages

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/global"
)

type Members struct {
	Members        []orderdomain.Member
	Loading, Error bool

	SelectedMember orderdomain.Member
	NewMember      bool
}

func (m *Members) Init() {
	m.Loading = true

	go func() {
		members, err := cache.Members.Get()
		defer global.LockAndRender()()
		m.Loading = false
		if err != nil {
			m.Error = true
			return
		}

		sort.Slice(members, func(i, j int) bool {
			return strings.Compare(members[i].Name, members[j].Name) < 0
		})

		m.Members = members
	}()
}

func (m *Members) ClickMember(member orderdomain.Member) {
	m.NewMember = false
	m.SelectedMember = member
}

func (m *Members) ClickNew() {
	m.SelectedMember = orderdomain.Member{}
	m.NewMember = true
}
