package handlers

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/dionisvl/go-uuid-inspector/internal/parser"
)

//go:embed templates/index.html
var indexHTML string

var tmpl = template.Must(template.New("index").Parse(indexHTML))

// Index handles the main page
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, parser.UUIDInfo{}); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// Inspect handles UUID inspection requests
func Inspect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	input := strings.TrimSpace(r.FormValue("uuid"))
	info := parser.Parse(input)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, info); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
