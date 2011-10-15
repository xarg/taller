// taller package provides template libraries with flow controls as close to
// golang as possible

package taller

// A template should have a Content which can Compile
type Template interface {
	Content() []byte
}

type TemplateBytes struct {
	data []byte
}

type TemplateFile struct {
	filename string
}

type Context map[string]interface{}

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
	return ReadTemplateFile(self.filename)
}

// Parse the template content, compile the resulted golang code
// execute and return the result
func Render(template Template, context Context) []byte {
	compiled_path := Compile(Parse(template.Content()))
	return Execute(compiled_path, context)
}

