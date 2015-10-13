package main

import (
	"encoding/base64"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sausheong/newspaper/paper"
	"github.com/sausheong/newspaper/today"
  "github.com/sausheong/newspaper/thesun"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

var Papers map[string]paper.Paper

func init() {
  Papers = make(map[string]paper.Paper)
	go func() {
		for {
			today := today.CheckAndLoad()
      Papers[today.Name] = today
      
			thesun := thesun.CheckAndLoad()
      Papers[thesun.Name] = thesun      
      
			time.Sleep(15 * time.Minute)
		}
	}()
}

func main() {	
	r := httprouter.New()
	r.ServeFiles("/public/*filepath", http.Dir("public/"))
  r.GET("/", index)
	r.GET("/paper/:paper", newspaper)
	r.GET("/paper/:paper/page/:page", page)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t, _ := template.ParseFiles("html/index.html")
	t.Execute(w, nil)
}

func newspaper(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("paper")
	t, _ := template.ParseFiles("html/page.html")
	t.Execute(w, name)
}

func page(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("paper")
	page, _ := strconv.Atoi(ps.ByName("page"))
	pg := Papers[name].Pages[page]
	format := `{"page": "%s", "num": %d}`
	jsonData := fmt.Sprintf(format, base64.StdEncoding.EncodeToString(pg), page)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}
