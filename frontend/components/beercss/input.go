package beercss

import (
	"fmt"
	"time"

	"github.com/PotatoesFall/vecty-test/frontend/jscall/beercss"
	"github.com/vugu/vugu"
)

type Input struct {
	Label string `vugu:"data"`
}

func (i *Input) Init(vugu.InitCtx) {
	fmt.Println(`init`)
}

func (i *Input) Compute(vugu.ComputeCtx) {
	fmt.Println(`compute`)
}

func (i *Input) Rendered(vugu.RenderedCtx) {
	fmt.Println(`rendered`)
	go func() {
		time.Sleep(100 * time.Millisecond)
		beercss.UI()
	}()
}

func (i *Input) Destroy(vugu.DestroyCtx) {
	fmt.Println(`destroyed`)
}
