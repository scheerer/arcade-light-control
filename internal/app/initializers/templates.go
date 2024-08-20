package initializers

import (
	_ "embed"
	"errors"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/scheerer/light-control/web"
)

var TemplateRenderer *TemplateRegistry

// Kind of dislike this current template setup but work for a simple app. Would need more abstraction for larger apps.
func LoadTemplates() {
	templates := make(map[string]*template.Template)

	templates["index.tmpl"] = template.Must(template.ParseFS(web.ViewTemplates, "template/index.tmpl", "template/layout.tmpl"))

	TemplateRenderer = &TemplateRegistry{
		templates: templates,
	}
}

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name+".tmpl"]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "layout.tmpl", data)
}
