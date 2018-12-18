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
	http.HandleFunc("/host/logout", h.logoutHost)
	http.HandleFunc("/host/session/view", h.viewSession)

}

func (h Host) servepage(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "host")
	if err != nil {
		log.Println(err)
	}

	host := session.Values["host"]
	if host == nil {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodGet {

		data := model.HostType{
			Email: host.(string),
		}

		val, err := model.GetHost(data)

		if err != nil {
			log.Println(err)
		}

		t := h.temp.Lookup("host.html")
		if t != nil {
			err = t.Execute(w, val)
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
			Email:   host.(string),
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
			set := exec.Command("docker", "container", "run", "-it", "-d", "--name", data.Image1+"1", data.Image1)
			if err := set.Run(); err != nil {
				log.Fatalln(err)
			} else {

				// exec into that container and share
				cmd := exec.Command("gotty", "-w", "-p", data.Port1, "docker", "container", "exec", "-it", data.Image1+"1", "sh")
				cmd.Stdout = os.Stdout

				if err = cmd.Run(); err != nil {
					log.Fatalln(err)
				}

			}

		}()

		go func() {

			set := exec.Command("docker", "container", "run", "-it", "-d", "--name", data.Image2+"2", data.Image2)
			if err := set.Run(); err != nil {
				log.Fatalln(err)
			} else {

				cmd := exec.Command("gotty", "-w", "-p", data.Port2, "docker", "container", "exec", "-it", data.Image2+"2", "sh")
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
	session, err := Store.Get(r, "host")
	if err != nil {
		log.Println(err)
	}
	if host := session.Values["host"]; host == nil {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}
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

func (h Host) logoutHost(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "host")
	if err != nil {
		log.Println(err)
	}

	session.Values["host"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Host) viewSession(w http.ResponseWriter, r *http.Request) {

	session, err := Store.Get(r, "host")
	if err != nil {
		log.Println(err)
	}

	if host := session.Values["host"]; host == nil {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	var data model.HostType
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	t := h.temp.Lookup("agora.html")
	if t != nil {
		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
