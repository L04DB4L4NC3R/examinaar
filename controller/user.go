package controller

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	temp *template.Template
}

func (u User) RegisterRoutes() {
	http.HandleFunc("/", u.ServeBase)
}

func (u User) ServeBase(w http.ResponseWriter, r *http.Request) {
	t := u.temp.Lookup("index.html")
	if t != nil {
		err := t.Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
