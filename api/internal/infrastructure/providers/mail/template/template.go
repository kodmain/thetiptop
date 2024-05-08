package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/kodmain/thetiptop/api/assets"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

const templatesPath = "mails"

// templates stocke les templates HTML et texte compilés.
var templates = make(map[string]*Template)

func init() {
	fmt.Println("init template")
	templates, err := fs.ReadDir(assets.Mails, templatesPath)
	logger.Error(err)
	fmt.Println(templates)

	loadTemplates()
	//loadHTMLTemplates()
	//loadTextTemplates()
}

func loadTemplates() {
	files, err := fs.ReadDir(assets.Mails, templatesPath)
	logger.Error(err)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		tmpl, err := template.ParseFS(assets.Mails, fmt.Sprintf("%s/%s", templatesPath, name))
		if logger.Error(err) {
			continue
		}

		ext := filepath.Ext(name)
		name = name[:len(name)-len(ext)]

		if ext == ".txt" {
			if existing, exists := templates[name]; exists {
				existing.Text = tmpl
			} else {
				templates[name] = &Template{Text: tmpl}
			}
		} else {
			if existing, exists := templates[name]; exists {
				existing.Html = tmpl
			} else {
				templates[name] = &Template{Html: tmpl}
			}
		}
	}
}

/*
func loadHTMLTemplates() {
	htmlFiles, err := fs.ReadDir(htmls, templatesPath)
	logger.Error(err)
	for _, file := range htmlFiles {
		if file.IsDir() {
			continue
		}
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		tmpl, err := template.ParseFS(htmls, path.Join(templatesPath, file.Name()))
		if logger.Error(err) {
			continue
		}

		if existing, exists := templates[name]; exists {
			existing.Html = tmpl
		} else {
			templates[name] = &Template{Html: tmpl}
		}
	}
}

func loadTextTemplates() {
	textFiles, err := fs.ReadDir(txts, templatesPath)
	if err != nil {
		logger.Error(err)
	}
	for _, file := range textFiles {
		if file.IsDir() {
			continue
		}
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		tmpl, err := template.New(file.Name()).ParseFS(txts, path.Join(templatesPath, file.Name()))
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
*/

// Template représente un template HTML et texte.
type Template struct {
	Text *template.Template
	Html *template.Template
}

// Inject insère des données dans les templates HTML et texte.
func (t *Template) Inject(data Data) ([]byte, []byte, error) {
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

	return nil
}
