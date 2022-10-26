package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// data to passed to the template
type templateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	IsAuthenticated int
	API             string
	CSSVersion      string
}

// Creating a variable to hold Functions that will be passed to the template
var functions = template.FuncMap{}

/*
	Embed the directory's templates using a Directive as a comment. This enables the
	Application to be compiled with all it's associated templates into a one binary
*/
//go:embed templates
var templateFS embed.FS

// Return Template Data.
// Info from the Struct above, or a variable created with a "type" of template data.
// Will add the default data using this function. td is a pointer to the template data, and makes an http request.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = app.config.api
	return td
}

/*
Template Render. This function takes a response writer, a request, the name
of the template to render, default data passed to templates, the data from the
struct, and partials.
*/
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	// Check to see if the variable from application config, is present in the template cache,(Ignore the first return parameter).
	_, templateInMap := app.templateCache[templateToRender]

	// Don't use the template cache when in development mode.
	// Changes to the base file need auto updating. Therefore, if the app's config environment is in production,
	// and the template is mapped, then use the template cache. And if it's not there, then it needs to be built. (parsing)
	if app.config.env == "production" && templateInMap {
		t = app.templateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	// Add default data (and request)
	td = app.addDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	return nil
}

// Parse Template function.
// It takes the parameters' partials as a slice, page, and template to render.
func (app *application) parseTemplate(partials []string, page string, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	// Build partials: Check to see if any partials are present.
	// If so, then range through the slice of strings.
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.tmpl", x)
		}
	}

	// If partials are present, then must call the parse template function/file system: ParseFS.
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", strings.Join(partials, ","), templateToRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)
	}
	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}
