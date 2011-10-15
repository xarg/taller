package taller

import (
	"testing"
	"os"
	"path"
	"bytes"
)

const (
	TESTDIR string = "test"
)

// return a template absolute path
func update_environment() {
	current_dir, _ := os.Getwd()
	err := os.Setenv(TALLER_ENV_VARIABLE, path.Join(current_dir, TESTDIR))
	if err != nil {
		panic("Failed to set environment variable" + err.String())
	}
}

//test if we can open a file and read the contents out of it
func TestTemplateFile(t *testing.T) {
	update_environment()
	template := NewTemplateFile("base.html")
	content := template.Content()
	if bytes.Count(content, []byte("<html>")) != 1 {
		t.Error("Cannot read the template")
	}
}

//test if the TemplateBytes works
func TestTemplateBytes(t *testing.T) {
	template := NewTemplateBytes([]byte("<html></html>"))
	content := template.Content()
	if bytes.Count(content, []byte("<html>")) != 1 {
		t.Error("Cannot read the template")
	}
}

//rendering base.html
func TestBaseCompile(t *testing.T) {
	/*
	update_environment()
	template := NewTemplateFile("base.html")
	rendered_content := Compile(template)
	expected_content := ReadTemplateFile("results/base.html")
	if string(rendered_content) != string(expected_content) {
		t.Errorf("Got: \n%s\nExpected: \n%s", rendered_content,
			expected_content)
	}
	*/
}

//test a simple [expand "template"] and [include "include.html"]
func TestExpandInclude(t *testing.T) {
	update_environment()
	template := NewTemplateFile("expand_include.html")
	rendered_content := Render(template)
	expected_content := ReadTemplateFile("results/expand_include.html")
	if string(rendered_content) != string(expected_content) {
		t.Errorf("Got: \n%s\nExpected: \n%s", rendered_content,
			expected_content)
	}
}
