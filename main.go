package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"./controller"
	"./model"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.Handle("/img/", http.FileServer(http.Dir("public")))
	http.Handle("/js/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("public")))
	http.Handle("/vendor/", http.FileServer(http.Dir("public")))
	controller.Startup(populateTemplates())

	db := connectDB()
	defer db.Close()
	model.SetDB(db)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func populateTemplates() *template.Template {
	result := template.New("index.html")

	template.Must(result.ParseGlob("./views/*.html"))
	return result

}

func connectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Connected to DB")
	model.SetDB(db)
	return db
}
