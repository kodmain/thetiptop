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

const (
	templatesPath = "mails"
	TXT           = ".txt"
	HTML          = ".html"
)

// templates stocke les templates HTML et texte compilés.
var templates = make(map[string]*Template)

func init() {
	loadTemplates()
}

// loadTemplates charge les templates depuis le système de fichiers.
//
// Cette fonction lit les fichiers dans un répertoire donné et les parse comme des templates.
// Elle gère deux types de templates : texte et HTML.
//
// Parameters:
// - aucun
//
// Returns:
// - aucun
func loadTemplates() {
	files, err := fs.ReadDir(assets.Mails, templatesPath)
	if logger.Error(err) {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		processTemplateFile(file.Name())
	}
}

// processTemplateFile traite un fichier de template unique.
//
// La fonction parse le fichier, vérifie son extension, et met à jour la map des templates en conséquence.
//
// Parameters:
// - name: string Le nom du fichier à traiter
//
// Returns:
// - aucun
func processTemplateFile(name string) {
	tmpl, err := template.ParseFS(assets.Mails, fmt.Sprintf("%s/%s", templatesPath, name))
	if logger.Error(err) {
		return
	}

	ext := filepath.Ext(name)
	name = name[:len(name)-len(ext)]
	updateTemplates(name, tmpl, ext)
}

// updateTemplates met à jour la map de templates basé sur le type de fichier.
//
// La fonction décide si le template est texte ou HTML et met à jour ou ajoute le template à la map.
//
// Parameters:
// - name: string Le nom du template sans extension
// - tmpl: *template.Template Le template parsé
// - ext: string L'extension du fichier
//
// Returns:
// - aucun
func updateTemplates(name string, tmpl *template.Template, ext string) {
	if ext == TXT {
		if existing, exists := templates[name]; exists {
			existing.Text = tmpl
		} else {
			templates[name] = &Template{Text: tmpl}
		}
	} else if ext == HTML {
		if existing, exists := templates[name]; exists {
			existing.Html = tmpl
		} else {
			templates[name] = &Template{Html: tmpl}
		}
	}
}

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
