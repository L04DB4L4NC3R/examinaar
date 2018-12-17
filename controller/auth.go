package controller

import (
	"log"
	"net/http"

	"../model"
)

func (h Host) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		err := r.ParseForm()

		if err != nil {
			log.Println(err)
		}

		f := r.Form

		data := model.HostType{
			Email:    f.Get("email"),
			Password: f.Get("password"),
		}
		_, err = model.CreateHost(data)

		if err != nil {
			w.Write([]byte("Error creating Host"))
			return
		}

		t := h.temp.Lookup("host.html")
		if t != nil {
			err = t.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
