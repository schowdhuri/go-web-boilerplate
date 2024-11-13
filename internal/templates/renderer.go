package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"viabl.ventures/gossr/internal/assets"
)

var funcMap template.FuncMap

type Renderer struct {
	templates map[string]*template.Template
	assetPipe *assets.AssetPipeline
	isDev     bool
}

func NewRenderer(assetPipe *assets.AssetPipeline, isDev bool) *Renderer {
	r := &Renderer{
		assetPipe: assetPipe,
		isDev:     isDev,
	}

	// Add custom functions to the template
	funcMap = template.FuncMap{
		"asset": assetPipe.GetAssetURL,
		"safeJS": func(s string) template.JS {
			return template.JS(s)
		},
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s)
		},
	}

	// Parse all templates
	r.parseTemplates()
	return r
}

func (r *Renderer) parseTemplates() {
	r.templates = make(map[string]*template.Template)
	tpl_map := map[string][]string{
		"home.html":       {"layouts/base.html", "pages/home.html"},
		"signin.html":     {"layouts/base.html", "pages/signin.html"},
		"admin_home.html": {"layouts/base.html", "pages/admin_home.html"},
	}

	for name, deps := range tpl_map {
		files := make([]string, len(deps))
		for i, dep := range deps {
			files[i] = filepath.Join("internal", "templates", dep)
		}
		tpl, err := template.New(name).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			fmt.Println("error parsing template", err)
		} else {
			r.templates[name] = tpl
		}
	}
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	// In development, reload templates on every request
	if r.isDev {
		r.parseTemplates()
	}
	// Buffer the output to handle errors and enable minification
	var buf bytes.Buffer
	// TODO: find a better way to render the base template
	if err := r.templates[name].ExecuteTemplate(&buf, "base.html", data); err != nil {
		fmt.Println("error executing template", err)
		return err
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(buf.Bytes()))
	return err
}
