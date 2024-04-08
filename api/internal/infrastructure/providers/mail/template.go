package mail

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

// Utilisez `go:embed` pour embarquer les fichiers de template.
//
//go:embed template/*.html
var htmls embed.FS

//go:embed template/*.txt
var txts embed.FS

// templates stocke les templates HTML et texte compilés.
var templates = make(map[string]*Template)

func init() {
	// Charger les templates HTML.
	htmlFiles, err := fs.ReadDir(htmls, "template")
	if err != nil {
		logger.Error(err)
	}
	for _, file := range htmlFiles {
		if file.IsDir() {
			continue
		}
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		tmpl, err := template.ParseFS(htmls, "template/"+file.Name())
		if err != nil {
			logger.Error(err)
			continue
		}
		if existing, exists := templates[name]; exists {
			existing.Html = tmpl
		} else {
			templates[name] = &Template{Html: tmpl}
		}
	}

	// Charger les templates texte.
	textFiles, err := fs.ReadDir(txts, "template")
	if err != nil {
		logger.Error(err)
	}
	for _, file := range textFiles {
		if file.IsDir() {
			continue
		}
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		tmpl, err := template.New(file.Name()).ParseFS(txts, "template/"+file.Name())
		if err != nil {
			logger.Error(err)
			continue
		}
		if existing, exists := templates[name]; exists {
			existing.Text = tmpl
		} else {
			templates[name] = &Template{Text: tmpl}
		}
	}
}

// Template représente un template HTML et texte.
type Template struct {
	Text *template.Template
	Html *template.Template
}

// Inject insère des données dans les templates HTML et texte.
func (t *Template) Inject(data any) ([]byte, []byte, error) {
	var html bytes.Buffer
	var text bytes.Buffer

	if t.Html != nil {
		if err := t.Html.Execute(&html, data); err != nil {
			return nil, nil, err
		}
	}

	if t.Text != nil {
		if err := t.Text.Execute(&text, data); err != nil {
			return nil, nil, err
		}
	}

	return text.Bytes(), html.Bytes(), nil
}

// NewTemplate retourne une nouvelle instance de Template basée sur le nom.
func NewTemplate(name string) *Template {
	if tmpl, exists := templates[name]; exists {
		return tmpl
	}

	return &Template{}
}
