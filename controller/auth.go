package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"../model"
)

var Store = sessions.NewCookieStore([]byte("SOME_SECRET"))

func (h Host) Signup(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		t := h.temp.Lookup("signup.html")
		if t != nil {
			err := t.Execute(w, nil)
			if err != nil {
				log.Println(err)
			}
		}
	}

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

		session, err := Store.Get(r, "host")
		if err != nil {
			log.Println(err)
		}
		session.Values["host"] = data.Email
		session.Save(r, w)

		t := h.temp.Lookup("host.html")
		if t != nil {
			err = t.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (h Host) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		t := h.temp.Lookup("login.html")
		if t != nil {
			err := t.Execute(w, nil)
			if err != nil {
				log.Println(err)
			}
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			log.Println(err)
			return
		}

		f := r.Form

		data := model.HostType{
			Email:    f.Get("email"),
			Password: f.Get("password"),
		}

		val, err := model.GetHost(data)

		if err != nil {
			log.Println(err)
			return
		}

		if len(val.Email) < 1 {
			w.Write([]byte("User not found"))
			return

		} else if val.Password == data.Password {

			session, err := Store.Get(r, "host")
			if err != nil {
				log.Println(err)
			}

			session.Values["host"] = val.Email
			session.Save(r, w)
			t := h.temp.Lookup("host.html")
			if t != nil {
				err = t.Execute(w, val)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			w.Write([]byte("Wrong password"))
			return
		}

	}
}
