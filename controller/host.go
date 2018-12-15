package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"../model"
)

type Host struct {
	temp *template.Template
}

func (h Host) RegisterRoutes() {
	http.HandleFunc("/host", h.servepage)
}

func (h Host) servepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := h.temp.Lookup("host.html")
		if t != nil {
			err := t.Execute(w, nil)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else if r.Method == http.MethodPost {
		var host model.HostType
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&host)
		if err == nil {
			log.Println(err)
		}

		enc := json.NewEncoder(w)

		enc.Encode(struct {
			done bool
		}{true})

		log.Println(host)
	}
}
