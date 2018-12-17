package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"../model"
)

type Host struct {
	temp *template.Template
}

func (h Host) RegisterRoutes() {
	http.HandleFunc("/host", h.servepage)
	http.HandleFunc("/host/signup", h.Signup)
	http.HandleFunc("/host/login", h.Login)
	http.HandleFunc("/host/session/delete", h.removeSession)
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
		err := r.ParseForm()

		if err != nil {
			log.Println(err)
		}

		f := r.Form

		data := model.HostType{
			Port1:   f.Get("port1"),
			Port2:   f.Get("port2"),
			Image1:  f.Get("image1"),
			Image2:  f.Get("image2"),
			Email:   f.Get("email"),
			Hosting: true,
		}

		_, err = model.CreateSessions(data)
		switch {
		case err == fmt.Errorf("Session already in place"):
			w.Write([]byte("Session already in place"))
			return
		case err != nil:
			log.Println(err)
			return

		}

		go func() {

			// run docker container
			set := exec.Command("container_setup", data.Image1, "1")
			if err := set.Run(); err != nil {
				log.Fatalln(err)
			} else {

				// exec into that container and share
				cmd := exec.Command("session", data.Port1, data.Image1, "1")
				cmd.Stdout = os.Stdout

				if err = cmd.Run(); err != nil {
					log.Fatalln(err)
				}

			}

		}()

		go func() {

			set := exec.Command("container_setup", data.Image2, "2")
			if err := set.Run(); err != nil {
				log.Fatalln(err)
			} else {

				cmd := exec.Command("session", data.Port2, data.Image2, "2")
				cmd.Stdout = os.Stdout
				if err = cmd.Run(); err != nil {
					log.Fatalln(err)
				}

			}

		}()

		time.Sleep(3 * time.Second)

		t := h.temp.Lookup("agora.html")
		if t != nil {
			err := t.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}
}

func (h Host) removeSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var e model.HostType

		err := json.NewDecoder(r.Body).Decode(&e)

		if err != nil {
			log.Println(err)
		}

		go func() {
			cmd := exec.Command("docker", "container", "rm", "-f", e.Image1+"1")
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				log.Println(err)
			}
		}()

		go func() {
			cmd := exec.Command("docker", "container", "rm", "-f", e.Image2+"2")
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				log.Println(err)
			}
		}()

		_, err = model.DeleteSessions(e.Email)

		if err != nil {
			log.Println(err)
		}

		w.Write([]byte("Deleted"))
	}
}
