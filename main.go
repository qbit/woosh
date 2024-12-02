package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

var (
	//go:embed templates
	templateFS embed.FS

	Templ *template.Template
)

func main() {
	Templ, err := template.New("production").ParseFS(templateFS, "templates/*")
	if err != nil {
		log.Fatalln(err)
	}

	r := http.NewServeMux()
	r.HandleFunc("GET /ui/{template}", func(w http.ResponseWriter, r *http.Request) {
		templ := r.PathValue("template")

		if err := Templ.ExecuteTemplate(w, templ, nil); err != nil {
			switch err := err.(type) {
			case *template.Error:
				_ = Templ.ExecuteTemplate(w, "error", err.ErrorCode)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}

			log.Println(err)

			return
		}
	})

	s := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	log.Fatal(s.ListenAndServe())
}
