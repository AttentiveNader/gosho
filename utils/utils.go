package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func ExecuteT(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(("ExcuteTemplate error : " + err.Error())))
	}
}
