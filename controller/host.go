package controller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
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

		port1 := f.Get("port1")
		port2 := f.Get("port2")

		if !CheckPort(port1, port2) {
			w.Write([]byte("Invalid value of port"))
			return
		}

		data := model.HostType{
			Port1:   port1,
			Port2:   port2,
			Image1:  f.Get("image1"),
			Image2:  f.Get("image2"),
			Email:   host.(string),
			Hosting: true,
			Channel: f.Get("channel"),
		}

		_, err_msg := model.CreateSessions(data)
		log.Println(err_msg)
		if len(err_msg) > 0 {
			w.Write([]byte(err_msg))
			return
		}

		go func() {

			// run docker container
			set := exec.Command("docker", "container", "run", "-it", "-d", "--name", data.Email+"1", data.Image1)
			if err := set.Run(); err != nil {
				log.Fatalln(err)
			} else {
				// exec into that container and share
				cmd := exec.Command("gotty", "-w", "-p", data.Port1, "docker", "container", "exec", "-it", data.Email+"1", "sh")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err = cmd.Run(); err != nil {
					log.Fatalln(err)
				}

			}

		}()

		go func() {

			set := exec.Command("docker", "container", "run", "-it", "-d", "--name", data.Email+"2", data.Image2)
			if err := set.Run(); err != nil {
				log.Fatalln(err)
			} else {

				cmd := exec.Command("gotty", "-w", "-p", data.Port2, "docker", "container", "exec", "-it", data.Email+"2", "sh")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
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

			// remove docker image
			cmd := exec.Command("docker", "container", "rm", "-f", e.Email+"1")
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				log.Println(err)
			}

			// kill gotty bounded to a particular port
			c1 := exec.Command("fuser", e.Port1+"/tcp")
			c2 := exec.Command("xargs", "kill")

			r, w := io.Pipe()
			c1.Stdout = w
			c2.Stdin = r

			var b2 bytes.Buffer
			c2.Stdout = &b2

			c1.Start()
			c2.Start()
			c1.Wait()
			w.Close()
			c2.Wait()
			io.Copy(os.Stdout, &b2)
		}()

		go func() {
			cmd := exec.Command("docker", "container", "rm", "-f", e.Email+"2")
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				log.Println(err)
			}

			c1 := exec.Command("fuser", e.Port2+"/tcp")
			c2 := exec.Command("xargs", "kill")

			r, w := io.Pipe()
			c1.Stdout = w
			c2.Stdin = r

			var b2 bytes.Buffer
			c2.Stdout = &b2

			c1.Start()
			c2.Start()
			c1.Wait()
			w.Close()
			c2.Wait()
			io.Copy(os.Stdout, &b2)
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
