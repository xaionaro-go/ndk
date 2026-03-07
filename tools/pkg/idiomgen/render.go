package idiomgen

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

// RenderTemplate parses and executes a Go text/template against a MergedSpec.
func RenderTemplate(name, tmplText string, data MergedSpec) (string, error) {
	tmpl, err := template.New(name).Funcs(FuncMap()).Parse(tmplText)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderTemplateFile reads a template file from disk and renders it against a MergedSpec.
func RenderTemplateFile(path string, data MergedSpec) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return RenderTemplate(filepath.Base(path), string(content), data)
}

// RenderPerType renders a template against a PerTypeData (single opaque type + its methods).
func RenderPerType(tmplPath string, data PerTypeData) (string, error) {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(FuncMap()).Parse(string(content))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderAny renders a template file against arbitrary data.
func RenderAny(tmplPath string, data any) (string, error) {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(FuncMap()).Parse(string(content))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderBridgeTemplate renders a bridge template against BridgeData.
func RenderBridgeTemplate(tmplPath string, data BridgeData) (string, error) {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(FuncMap()).Parse(string(content))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
