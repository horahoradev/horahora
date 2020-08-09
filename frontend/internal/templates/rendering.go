package templates

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func New() *Template {
	t := &Template{
		templates: template.Must(template.ParseGlob("./internal/templates/html/*.gohtml")),
	}

	return t
}
