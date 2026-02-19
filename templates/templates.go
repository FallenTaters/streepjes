package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"sync"
)

//go:embed *.html admin/*.html
var files embed.FS

var (
	once  sync.Once
	pages map[string]*template.Template
)

func buildPages() {
	pages = make(map[string]*template.Template)

	pageFiles := []string{
		"login.html",
		"profile.html",
		"order.html",
		"history.html",
		"leaderboard.html",
		"admin/billing.html",
		"admin/catalog.html",
		"admin/members.html",
		"admin/users.html",
	}

	shared := []string{"base.html", "nav.html"}

	for _, pf := range pageFiles {
		t := template.Must(template.ParseFS(files, append(shared, pf)...))
		pages[pf] = t
	}
}

func Render(w io.Writer, name string, data any) error {
	once.Do(buildPages)
	t, ok := pages[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	return t.ExecuteTemplate(w, "base", data)
}
