package controller

import "html/template"

var (
	user User
	host Host
)

func Startup(t *template.Template) {
	user.temp = t
	host.temp = t
	user.RegisterRoutes()
	host.RegisterRoutes()
}
