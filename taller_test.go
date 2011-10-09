package taller

import (
	"testing"
	"os"
	"path"
	"bytes"
	"io/ioutil"
)
const (
	TESTDIR string = "test"
)

// return a template absolute path
func template_path(template_name string) string {
	current_dir, _ := os.Getwd()
	return path.Join(current_dir, TESTDIR, template_name)
}

func read_template(template_name string) []byte {
	filepath := template_path(template_name)
	content, _ := ioutil.ReadFile(filepath)
	return content
}

//test if we can open a file and read the contents out of it
func TestTemplateFile(t *testing.T) {
	template := NewTemplateFile(template_path("base.html"))
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
func TestBaseRender(t *testing.T) {
	template := NewTemplateFile(template_path("base.html"))
	rendered_content := Render(template)
	expected_content := read_template("results/base.html")
	if string(rendered_content) != string(expected_content) {
		t.Errorf("Got: \n%s\nExpected: \n%s", rendered_content,
		expected_content)
	}
}

