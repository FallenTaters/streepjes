package pages

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/components/beercss"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Members struct {
	Members        []orderdomain.Member
	Loading, Error bool

	SelectedMember orderdomain.Member
	NewMember      bool

	Name string

	LoadingForm, ErrorForm bool
}

func (m *Members) Init() {
	m.Loading = true
	m.Error = false
	m.SelectedMember = orderdomain.Member{}
	m.NewMember = false
	m.Name = ``

	go func() {
		members, err := cache.Members.Get()
		defer global.LockAndRender()()
		m.Loading = false
		if err != nil {
			m.Error = true
			return
		}

		sort.Slice(members, func(i, j int) bool {
			name1, name2 := strings.ToLower(members[i].Name), strings.ToLower(members[j].Name)
			return strings.Compare(name1, name2) < 0
		})

		m.Members = members
	}()
}

func (m *Members) ClickMember(member orderdomain.Member) {
	m.LoadingForm = false
	m.ErrorForm = false
	m.NewMember = false

	m.SelectedMember = member
	m.Name = member.Name
}

func (m *Members) ClickNew() {
	m.LoadingForm = false
	m.ErrorForm = false
	m.SelectedMember = orderdomain.Member{}

	m.NewMember = true
	m.Name = ``
}

func (m *Members) ShowForm() bool {
	return m.NewMember || m.Editing()
}

func (m *Members) Editing() bool {
	return m.SelectedMember != (orderdomain.Member{})
}

func (m *Members) FormTitle() string {
	if m.Editing() {
		return `Edit ` + m.SelectedMember.Name
	}

	if m.NewMember {
		return `New Member`
	}

	return ``
}

func (m *Members) ClubOptions() []beercss.Option {
	return []beercss.Option{
		{
			Label: domain.ClubGladiators.String(),
			Value: domain.ClubGladiators,
		},
		{
			Label: domain.ClubParabool.String(),
			Value: domain.ClubParabool,
		},
	}
}

func (m *Members) SubmitForm() {
	if m.Editing() {
		m.doUpdate(func() error {
			return backend.PostEditMember(orderdomain.Member{
				ID:   m.SelectedMember.ID,
				Club: store.Auth.User.Club,
				Name: m.Name,
			})
		})
		return
	}

	m.doUpdate(func() error {
		return backend.PostNewMember(orderdomain.Member{
			Club: store.Auth.User.Club,
			Name: m.Name,
		})
	})
}

func (m *Members) Delete() {
	if !window.Confirm(`Are you sure?`) {
		return
	}

	m.doUpdate(func() error {
		return backend.PostDeleteMember(m.SelectedMember.ID)
	})
}

func (m *Members) doUpdate(update func() error) {
	m.LoadingForm = true
	m.ErrorForm = false
	go func() {
		err := update()
		defer global.LockAndRender()()
		defer func() { m.LoadingForm = false }()
		if err != nil {
			m.ErrorForm = true
			return
		}

		m.Init()
	}()
}
