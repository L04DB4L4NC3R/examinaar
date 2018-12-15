package controller

import "html/template"

var (
	user User
)

func Startup(t *template.Template) {
	user.temp = t
	user.RegisterRoutes()
}
