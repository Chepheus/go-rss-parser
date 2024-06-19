package services

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type TemplateData map[string]interface{}

type HTMLRenderer struct {
	templatePath string
}

func (r HTMLRenderer) Render(w http.ResponseWriter, templateName string, data *TemplateData) {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.ParseFiles(r.templatePath + templateName + ".html")
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}
}

func NewHTMLRenderer(templatePath string) HTMLRenderer {
	return HTMLRenderer{
		templatePath: templatePath,
	}
}
