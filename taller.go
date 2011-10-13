// taller package provides template libraries with flow controls as close to
// the native language as possible

package taller

import (
	"bytes"
//	"fmt"
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

type Context map[string] interface{}

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
	parsed_content := Parse(content)
	//compile and execute with the context
	return parsed_content
}

func Parse(content []byte) []byte {
	output := []byte(`package taller
import (
	"fmt"
	"bytes"
	/*IMPORTS*/
)
/*BLOCKS*/
func render(context Context) *bytes.Buffer {
	output := bytes.NewBufferString("")
	/*MAIN_BLOCK*/
	return output
}`)
	blocks := map[string] *bytes.Buffer {
		"IMPORTS":bytes.NewBufferString(""),
		"BLOCKS":bytes.NewBufferString(""),
		"MAIN_BLOCK":bytes.NewBufferString("")}
	lines := bytes.Split(content, []byte("\n"))
	for lineno, line := range lines {
		//each line in the template will be parsed an translated to the
		//appropriate block
		block, parsed_line := ParseLine(lineno, line)
		blocks[block].Write(parsed_line)
	}
	for block, parsed_lines := range blocks {
		//replaces places such as /*MAIN_BLOCK*/ in the output using
		//the parsed_lines
		bytes.Replace(output, JoinBytes("/*", block, "*/"),
			parsed_lines.Bytes(), -1)
	}
	return output
}
//give some string join them a single []byte seq
func JoinBytes(string_list ... string) []byte {
	var joined [][]byte
	for _, str := range string_list {
		joined = append(joined, []byte(str))
	}
	return bytes.Join(joined, []byte(""))
}

func ParseLine (lineno int, line []byte) (string, []byte) {
	//On the first line should always be [extends "filename"]
	var block_name string
	var parsed_line []byte
	if lineno == 0 {
	}
	block_name = "aaa"
	parsed_line = []byte("aaa")
	return block_name, parsed_line
}
