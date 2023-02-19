package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/caleberi/snippet-bin/pkg/forms"
	"github.com/caleberi/snippet-bin/pkg/models"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippetDb.Lastest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: snippets})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires").
		MaxLength("title", 100).
		PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	title := form.Get("title")
	content := form.Get("content")
	expires := form.Get("expires")

	id, err := app.snippetDb.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLogger.Printf("Created snippet with id %d", id)
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippetDb.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{Snippet: snippet})
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {

	ts, ok := app.templateCache[name]
	if ok {
		buf := bytes.NewBuffer([]byte{})

		err := ts.Execute(buf, app.addDefaultData(data, r))
		if err != nil {
			app.serverError(w, err)
			return
		}
		buf.WriteTo(w)
		return
	} else {
		app.serverError(w, fmt.Errorf("failed to execute template %s", name))
		return
	}

}

func (app *application) addDefaultData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}

	data.CurrentYear = time.Now().Year()
	return data
}
