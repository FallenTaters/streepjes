package router

import (
	"net/http"

	"git.fuyu.moe/Fuyu/router"
)

type Static func(filename string) ([]byte, error)

func New(static Static) *router.Router {
	r := router.New()

	r.ErrorHandler = panicHandler

	r.GET(`/version`, getVersion)
	r.GET(`/`, getIndex(static))
	r.GET(`/static/:name`, getStatic(static), compressionMiddleware)

	return r
}

func getVersion(c *router.Context) error {
	return c.String(http.StatusOK, version())
}

func getIndex(static Static) router.Handle {
	return func(c *router.Context) error {
		index, err := static(`index.html`)
		if err != nil {
			panic(err)
		}

		c.Response.Header().Set(`Content-Type`, `text/html`)
		return c.Bytes(http.StatusOK, index)
	}
}

func getStatic(static Static) router.Handle {
	return func(c *router.Context) error {
		name := c.Param(`name`)
		asset, err := static(name)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		_, err = c.Response.Write(asset)
		return err
	}
}
