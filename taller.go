// taller package provides template libraries with flow controls as close to the native language as possible

package taller

import (
	//"fmt"
	"io/ioutil"
)

const (
	SECTION_RIGHT_DELIMITER = '['
	SECTION_LEFT_DELIMITER = '}'
	ECHO_DELIMITER = '`'
)

// A template should have a Content which can Render
type Template interface {
	Content() []byte
}

type TemplateBytes struct {
	data []byte
}

type TemplateFile struct {
	filename string
}

func NewTemplateBytes(data []byte) *TemplateBytes {
	template := new(TemplateBytes)
	template.data = data
	return template
}

func NewTemplateFile(filename string) *TemplateFile {
	template_file := new(TemplateFile)
	template_file.filename = filename
	return template_file
}

func (self *TemplateBytes) Content() []byte {
	return self.data
}

func (self *TemplateFile) Content() []byte {
	content, err := ioutil.ReadFile(self.filename)
	if err != nil {
		panic("taller error: " + err.String())
	}
	return content
}

func Render(template Template) []byte {
	content := template.Content()
	return content
}

func Parse(content []byte) {
}
