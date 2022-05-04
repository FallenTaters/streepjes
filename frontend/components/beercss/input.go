package beercss

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/beercss"
	"github.com/FallenTaters/streepjes/frontend/jscall/document"
	"github.com/vugu/vugu"
)

type Input struct {
	AttrMap vugu.AttrMap

	ID    string
	Label string `vugu:"data"`

	// oh god this is such a mess just don't touch it
	Value            string `vugu:"data"` // this is the prop, but often changes upstream after input is handled
	OriginalValue    string `vugu:"data"` // this is the previous value of the prop, to see if it changes
	LastUpdate       string `vugu:"data"` // last update to prevent syscall during typing, when value changes to the last update
	BeforeLastUpdate string `vugu:"data"` // for some reason there is usually delay by one

	Input        InputHandler `vugu:"data"`
	ShowPassword bool         `vugu:"data"`
}

func (i *Input) Init() {
	binaryToken := make([]byte, base64.RawURLEncoding.DecodedLen(10)+1)

	_, err := rand.Read(binaryToken)
	if err != nil {
		panic(err)
	}

	i.ID = base64.RawURLEncoding.EncodeToString(binaryToken)[:10]
}

func (i *Input) HandleChange(event vugu.DOMEvent) {
	v := event.JSEventTarget().Get(`value`).String()

	i.BeforeLastUpdate = i.LastUpdate
	i.LastUpdate = v

	if i.Input != nil {
		go i.Input.InputHandle(InputEvent(v))
	}
}

// Replace with Rendered() once https://github.com/vugu/vugu/issues/224 is resolved, no lock needed
func (i *Input) Compute(vugu.ComputeCtx) {
	go func() {
		defer global.LockOnly()()

		beercss.UI()

		if i.Value != i.OriginalValue && i.Value != i.LastUpdate && i.Value != i.BeforeLastUpdate {
			fmt.Println(i.Value, i.OriginalValue, i.LastUpdate)
			elem, ok := document.GetElementById(i.ID)
			if !ok {
				return
			}

			elem.Set(`value`, i.Value)
			i.OriginalValue = i.Value
		}
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

	if _, ok := attrs[`id`]; !ok {
		attrs[`id`] = i.ID
	}

	return attrs
}

func (i *Input) Classes() string {
	if i.AttrMap[`type`] == `password` {
		return `suffix`
	}

	return ``
}
