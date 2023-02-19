package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/caleberi/snippet-bin/pkg/forms"
	"github.com/caleberi/snippet-bin/pkg/models"
)

type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var funcs = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template, 0)

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		filename := filepath.Base(page)

		ts, err := template.New(filename).Funcs(funcs).ParseFiles(page)

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))

		if err != nil {
			return nil, err
		}

		cache[filename] = ts
	}

	return cache, nil
}
