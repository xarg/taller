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
func TestBaseRender(t *testing.T) {
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
	rendered_content := Render(template, Context{})
	expected_content := ReadTemplateFile("results/expand_include.html")
	if string(rendered_content) != string(expected_content) {
		t.Errorf("Got: \n%s\nExpected: \n%s", rendered_content,
			expected_content)
	}
}

//utils

//Check if split of paths works the way it's suppose to
func TestGetTallerPaths(t *testing.T) {
	os.Setenv(TALLER_ENV_VARIABLE, "aaa/b:xxx/yy")
	paths := GetTallerPaths()
	if len(paths) != 2 {
		t.Errorf("Got: %d\nExpected: 2\n", len(paths))
	}
	if paths[0] != "aaa/b" {
		t.Errorf("Got: %s\nExpected: %s\n", paths[0], "aaa/b")
	}
}

//testing order resolution when reading the template.
func TestReadTemplateFile(t *testing.T) {
	current_dir, _ := os.Getwd()
	os.Setenv(TALLER_ENV_VARIABLE, path.Join(current_dir, TESTDIR, "results")+":"+path.Join(current_dir, TESTDIR))
	content := ReadTemplateFile("base.html")
	expected := []byte(`<body>
content`)
	if bytes.Count(content, expected) != 1 {
		t.Errorf("Cannot find:\n %s\n in:\n %s\n", expected, content)
	}
}

func TestJoinBytes(t *testing.T) {
	if !bytes.Equal(JoinBytes("a"), []byte("a")) {
		t.Errorf("Got: %s\nExpected: %s\n", JoinBytes("a"), "a")
	}
	if !bytes.Equal(JoinBytes("a", "b"), []byte("ab")) {
		t.Errorf("Got: %s\nExpected: %s\n", JoinBytes("ab"), "ab")
	}
}

func TestSplitToLines(t *testing.T) {
	splited := SplitToLines([]byte("a\nb"))
	if !bytes.Equal(splited[0], []byte("a")) {
		t.Errorf("Got: %s\nExpected: %s\n", splited[0], []byte("a"))
	}
}
