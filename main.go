package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, here is hello-go</h1>")
}

func abouthandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This repository for recording learnning "+
		"<a href=\"mailto:albert@example.com\">albert@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>not found page!</h1>")

}
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "article ID: "+id)
}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "visit list")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "title is required"
	} else if len(title) < 3 || len(title) > 40 {
		errors["title"] = "title length is limited 3-40"
	}

	if body == "" {
		errors["body"] = "body is required"
	} else if len(body) < 10 {
		errors["body"] = "body length must > 10"
	}

	if len(errors) == 0 {
		fmt.Fprint(w, "verify pass<br>")
		fmt.Fprintf(w, "title: %v<br>", title)
		fmt.Fprintf(w, "title length: %v<br>", len(title))

		fmt.Fprintf(w, "body: %v<br>", body)
		fmt.Fprintf(w, "body length: %v<br>", len(body))
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.html")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func main() {

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", abouthandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
