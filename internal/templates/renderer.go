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
	templates *template.Template
	assetPipe *assets.AssetPipeline
	cache     map[string]*template.Template
	isDev     bool
}

func NewRenderer(assetPipe *assets.AssetPipeline, isDev bool) *Renderer {
	r := &Renderer{
		assetPipe: assetPipe,
		cache:     make(map[string]*template.Template),
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

	// Parse base templates
	r.templates = template.New("").Funcs(funcMap)

	// Parse all templates
	if err := r.parseTemplates(); err != nil {
		panic(err)
	}

	return r
}

func (r *Renderer) parseTemplates() error {
	// Parse all .html files in templates directory
	pattern := filepath.Join("internal", "templates", "**", "*.html")
	templateFiles, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	_, err = r.templates.ParseFiles(templateFiles...)
	return err
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	// In development, reload templates on every request
	if r.isDev {
		r.templates = template.New("").Funcs(funcMap)
		if err := r.parseTemplates(); err != nil {
			fmt.Println("error parsing templates", err)
			return err
		}
	}
	tmpl := r.templates.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("template %s not found", name)
	}

	// Buffer the output to handle errors and enable minification
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Println("error executing template", err)
		return err
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(buf.Bytes()))
	return err
}
