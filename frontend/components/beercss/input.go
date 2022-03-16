package beercss

import (
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/beercss"
	"github.com/vugu/vugu"
)

type Input struct {
	AttrMap vugu.AttrMap

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

func (i *Input) IsPassword() bool {
	return i.AttrMap[`type`] == `password`
}

func (i *Input) Attrs() vugu.AttrMap {
	attrs := make(vugu.AttrMap)
	for k, v := range i.AttrMap {
		attrs[k] = v
	}

	if i.IsPassword() && i.ShowPassword {
		attrs[`type`] = `text`
	}

	return attrs
}

func (i *Input) Classes() string {
	if i.AttrMap[`type`] == `password` {
		return `suffix`
	}

	return ``
}
