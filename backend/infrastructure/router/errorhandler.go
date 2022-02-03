package router

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"git.fuyu.moe/Fuyu/router"
)

func panicHandler(c *router.Context, v interface{}) {
	fmt.Fprintf(os.Stderr, "panic: %s\n", v)
	debug.PrintStack()
	fmt.Println()
	_ = c.NoContent(http.StatusInternalServerError)
}
