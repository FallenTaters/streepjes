package beercss

import "github.com/vugu/vugu"

// Expansion doesn't work with input fields or click listeners inside it because it sucks
type Expansion struct {
	Summary vugu.Builder
	Content vugu.Builder
}
