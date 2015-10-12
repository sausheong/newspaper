package main

import (
	// "io/ioutil"
	"fmt"
	"net/http"
	// "log"
	// "encoding/base64"
	"encoding/base64"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"os"
	"strconv"
)

func main() {
	r := httprouter.New()
	r.ServeFiles("/public/*filepath", http.Dir("public/"))

	r.GET("/paper/:paper", paper)
	r.GET("/paper/:paper/page/:page", page)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	server.ListenAndServe()
}

func paper(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// name := ps.ByName("paper")
	t, _ := template.ParseFiles("html/index.html")
	t.Execute(w, nil)
}

func page(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("paper")
	page, _ := strconv.Atoi(ps.ByName("page"))
	pg := papers[name].Pages[page]
	format := `{"page": "%s", "num": %d}`
	jsonData := fmt.Sprintf(format, base64.StdEncoding.EncodeToString(pg), page)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}
