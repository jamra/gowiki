// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"text/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"github.com/knieriem/markdown"
	//"unicode/utf8"
//	"fmt"
	"bytes"
	//"bufio"
)

const data_dir string = "./data/"
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := data_dir + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}


func loadPage(title string) (*Page, error) {
	filename := data_dir + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	p.Body = GetWikiMarkup(p.Body)

	renderTemplate(w, "view", p)
}

func GetWikiMarkup(data []byte) []byte {
	output := bytes.NewBuffer(nil)
	reader := bytes.NewBuffer(data)

	parser.Markdown(reader, markdown.ToHTML(output))
	return output.Bytes()
}


func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func directoryHandler(w http.ResponseWriter, r *http.Request) {
	d, err := loadDirectory()

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "directory", d)
}

type Directory struct{
	Name string
}

func loadDirectory() (d []Directory, err error){
	dirs, dirErr := ioutil.ReadDir(data_dir)

	if dirErr != nil{
		err = dirErr
	}

	for _, dir := range(dirs){
		name := dir.Name()
		if !dir.IsDir() && len(name) > 4 {
			d = append(d, Directory{Name: name[:len(name)-4]})
		}
	}

	return
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	viewHandler(w, r, "FrontPage")
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("./tmpl/directory.html", "./tmpl/edit.html", "./tmpl/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	//fmt.Println("Render Template: ", tmpl)
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const lenPath = len("/view/")

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}

var opt markdown.Extensions
var parser *markdown.Parser
func main() {

	var mde = markdown.Extensions{
		Smart:        true,
		Dlists:       true,
		FilterHTML:   true,
		FilterStyles: true,
	}

	parser = markdown.NewParser(&mde)

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})
	http.HandleFunc("/", frontPageHandler)
	http.HandleFunc("/directory/", directoryHandler)

	http.ListenAndServe(":8080", nil)
}
