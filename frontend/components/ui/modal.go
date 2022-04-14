package ui

import "github.com/vugu/vugu"

type Modal struct {
	Show  bool         `vugu:"data"`
	Close CloseHandler `vugu:"data"`

	DefaultSlot vugu.Builder
}

func (m *Modal) CloseModal() {
	m.Close.CloseHandle(CloseEvent{})
}
