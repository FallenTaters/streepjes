package ui

import "github.com/vugu/vugu"

type Modal struct {
	Show     bool         `vugu:"data"`
	Close    CloseHandler `vugu:"data"`
	InHeader vugu.Builder `vugu:"data"`
	Side     string       `vugu:"data"`

	DefaultSlot vugu.Builder
}

func (m *Modal) CloseModal() {
	m.Close.CloseHandle(CloseEvent{})
}

func (m *Modal) class() string {
	return `modal active ` + m.Side
}
