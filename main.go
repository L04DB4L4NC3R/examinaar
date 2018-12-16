package main

import (
	"html/template"
	"log"
	"net/http"

	"./controller"
)

func main() {
	http.Handle("/img/", http.FileServer(http.Dir("public")))
	http.Handle("/js/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("public")))
	http.Handle("/vendor/", http.FileServer(http.Dir("public")))
	controller.Startup(populateTemplates())
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func populateTemplates() *template.Template {
	result := template.New("index.html")

	template.Must(result.ParseGlob("./views/*.html"))
	return result

}
