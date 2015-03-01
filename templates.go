package main

import (
	"html/template"
	"io/ioutil"
)

var (
	templateFilenames = map[string][]string{
		"allPages": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"all-pages.html",
		},
		"deleted": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"deleted.html",
		},
		"diff": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"_delete.html",
			"_edit.html",
			"diff.html",
		},
		"edit": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"_delete.html",
			"_edit.html",
			"edit.html",
		},
		"history": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"history.html",
		},
		"printable": []string{
			"printable.html",
		},
		"preview": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"_delete.html",
			"_edit.html",
			"preview.html",
		},
		"search": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"search.html",
		},
		"view": []string{
			"_base.html",
			"_head.html",
			"_body.html",
			"_delete.html",
			"view.html",
		},
	}
)

func NewTemplates() (map[string]*template.Template, error) {
	fs := FS(false)

	templates := make(map[string]*template.Template)

	for templateName, filenames := range templateFilenames {
		template := template.New(templateName)
		templates[templateName] = template

		for _, filename := range filenames {
			file, err := fs.Open("/templates/" + filename)
			if err != nil {
				return nil, err
			}
			var c []byte
			c, err = ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}

			_, err = template.Parse(string(c))
			if err != nil {
				return nil, err
			}
		}
	}
	return templates, nil
}
