package beercss

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/beercss"
	"github.com/vugu/vugu"
)

type Select struct {
	AttrMap vugu.AttrMap

	ID    string
	Label string `vugu:"data"`

	Options []Option

	Select SelectHandler `vugu:"data"`
}

type Option struct {
	Label string
	Value any
}

func (s *Select) Init() {
	binaryToken := make([]byte, base64.RawURLEncoding.DecodedLen(10)+1)

	_, err := rand.Read(binaryToken)
	if err != nil {
		panic(err)
	}

	s.ID = base64.RawURLEncoding.EncodeToString(binaryToken)[:10]
}

func (s *Select) HandleChange(event vugu.DOMEvent) {
	v := event.JSEventTarget().Get(`value`).String()
	i, err := strconv.Atoi(v)
	if err != nil {
		return
	}

	if s.Select != nil && i >= 0 && i < len(s.Options) {
		go s.Select.SelectHandle(SelectEvent(s.Options[i].Value))
	}
}

// Replace with Rendered() once https://github.com/vugu/vugu/issues/224 is resolved, no lock needed
func (s *Select) Compute(vugu.ComputeCtx) {
	go func() {
		defer global.LockOnly()()

		beercss.UI()
	}()
}

func (s *Select) Attrs() vugu.AttrMap {
	attrs := make(vugu.AttrMap)
	for k, v := range s.AttrMap {
		attrs[k] = v
	}

	if _, ok := attrs[`id`]; !ok {
		attrs[`id`] = s.ID
	}

	return attrs
}
