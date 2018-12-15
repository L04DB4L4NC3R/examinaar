package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"

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
		data := model.HostType{}
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			log.Println(err)
		}

		enc := json.NewEncoder(w)

		// executing terminal logic
		go func() {
			cmd := exec.Command("session", data.Port1, data.Image1)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("Running %s on port %d", data.Image1, data.Port1)
			}
		}()

		go func() {
			cmd := exec.Command("session", data.Port2, data.Port1)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("Running %s on port %d", data.Port1, data.Port2)
			}
		}()

		enc.Encode(struct {
			Done bool `json:"done"`
		}{true})

	}
}
