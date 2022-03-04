package router

import (
	"fmt"

	"git.fuyu.moe/Fuyu/router"
	"github.com/FallenTaters/streepjes/backend/application/auth"
)

func publicRoutes(r *router.Router, static Static, authService auth.Service) {
	r.GET(`/`, getIndex(static))
	r.GET(`/version`, getVersion)
	r.GET(`/static/*name`, getStatic(static))
	r.POST(`/login`, postLogin(authService))
}

var (
	buildVersion string
	buildCommit  string
	buildTime    string
)

func version() string {
	return fmt.Sprintf("Version: %s\nDate: %s\nCommit: %s",
		buildVersion, buildTime, buildCommit)
}
