package main

import (
	"html/template"
	"log"
	"net/http"

	"./controller"
)

func main() {
	controller.Startup(populateTemplates())
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func populateTemplates() *template.Template {
	result := template.New("index.html")

	template.Must(result.ParseGlob("./views/*.html"))
	return result

}
