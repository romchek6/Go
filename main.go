package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

type Article struct {
	Id                    uint16
	Title, Anons, Content string
}

var posts = []Article{}
var postView = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	connStr := "user=postgres password=1606 dbname=Go3 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err)
	}

	posts = []Article{}

	for res.Next() {

		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Content)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)

	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		title := r.FormValue("title")
		anons := r.FormValue("anons")
		content := r.FormValue("content")

		if title == "" || anons == "" || content == "" {

			fmt.Fprintf(w, "Не все данные заполнены")

		} else {

			connStr := "user=postgres password=1606 dbname=Go3 sslmode=disable"
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (title , anons , content) values ('%s' , '%s' , '%s')", title, anons, content))
			if err != nil {
				panic(err)
			}
			defer insert.Close()

			http.Redirect(w, r, "/create/", http.StatusSeeOther)

		}

	} else {

		t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.ExecuteTemplate(w, "create", nil)

	}

}

func showPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/post.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	connStr := "user=postgres password=1606 dbname=Go3 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id ='%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	postView = Article{}

	for res.Next() {

		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Content)
		if err != nil {
			panic(err)
		}

		postView = post

	}

	t.ExecuteTemplate(w, "post", postView)

}

func handleFunc() {

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create/", create)
	router.HandleFunc("/post/{id:[0-9]+}/", showPost).Methods("GET")
	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)

}

func main() {

	handleFunc()

}
