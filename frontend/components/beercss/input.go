package beercss

import (
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/beercss"
	"github.com/vugu/vugu"
)

type Input struct {
	Type  string `vugu:"data"`
	Label string `vugu:"data"`
	Value string `vugu:"data"`

	Input        InputHandler `vugu:"data"`
	ShowPassword bool         `vugu:"data"`
}

func (i *Input) HandleChange(event vugu.DOMEvent) {
	v := event.JSEventTarget().Get(`value`).String()

	i.Value = v

	if i.Input != nil {
		go i.Input.InputHandle(InputEvent(v))
	}
}

// Replace with Rendered() once https://github.com/vugu/vugu/issues/224 is resolved, no lock needed
func (i *Input) Compute(vugu.ComputeCtx) {
	go func() {
		defer global.LockOnly()()

		beercss.UI()
	}()
}

func (i *Input) GetType() string {
	if i.Type == `password` && i.ShowPassword {
		return `string`
	}

	return i.Type
}

func (i *Input) Classes() string {
	if i.Type == `password` {
		return `prefix`
	}

	return ``
}
