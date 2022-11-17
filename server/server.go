package server

import (
	ascii "ascii-art-web/ascii"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const (
	Host = "localhost"
	Port = "8080"
)

type Status struct {
	BadRequest  int
	NotFound    int
	ServerError int
	OK          int
}

var myError = Status{
	BadRequest:  400,
	NotFound:    404,
	ServerError: 500,
	OK:          200,
}

type Font struct {
	Standard   string
	Shadow     string
	Thinkertoy string
}

var Form struct {
	content string
	font    string
}

type Data struct {
	content    string
	statusCode int
}

func doMissingPage() error {
	return &ascii.RequestError{
		StatusCode: 404,
		Err:        errors.New("page not found"),
	}
}

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/error.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data Data) {
	w.WriteHeader(data.statusCode)
	err := templates.ExecuteTemplate(w, tmpl+".html", data.content)
	if err != nil {
		////////// temp
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Post(w http.ResponseWriter, r *http.Request) error {

	content := r.FormValue("inputText")
	font := r.FormValue("font")

	result, err := ascii.Printer(content, font)
	if err != nil {
		data := Data{
			content:    err.Error(),
			statusCode: 400,
		}
		renderTemplate(w, "error", data)
		return err
	}
	data := Data{
		content:    result,
		statusCode: 200,
	}
	renderTemplate(w, "index", data)
	//templates.Execute(w, result)
	return nil
}

func Get(w http.ResponseWriter, r *http.Request) {
	templates.Execute(w, "")
}
func requestHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		fmt.Println(r.URL.Path)
		if r.URL.Path != "/favicon.ico" && r.URL.Path != "/ascii-art" && r.URL.Path != "/" {
			data := Data{
				content:    doMissingPage().Error(),
				statusCode: 404,
			}
			renderTemplate(w, "error", data)
			return
		}
		Get(w, r)
		return
	case "POST":
		Post(w, r)
		return
	default:
		data := Data{
			content:    "",
			statusCode: 500,
		}
		renderTemplate(w, "error", data)
		// fix it later
		http.Redirect(w, r, "error.html", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}

}

func Start() {
	fs := http.FileServer(http.Dir("templates/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.HandleFunc("/", requestHandler)
	err := http.ListenAndServe(Host+":"+Port, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
		return
	}
}
