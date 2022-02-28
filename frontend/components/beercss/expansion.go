package beercss

import "github.com/vugu/vugu"

type Expansion struct {
	Open bool `vugu:"data"`

	Summary vugu.Builder
	Content vugu.Builder
}
