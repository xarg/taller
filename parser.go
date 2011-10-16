package taller

import (
	"regexp"
	"bytes"
)

const (
	ECHO_DELIMITER = '`'

	EXPAND_PATTERN  = `\[expand ([^\]]+)\]`
	INCLUDE_PATTERN = `\[include ([^\]]+)\]`
	IMPORT_PATTERN  = `\[import ([^\]]+)\]`

	BLOCK_TAG_PATTERN    = `\[block ([^\]]+)\]`
	ENDBLOCK_TAG_PATTERN = `\[endblock\]`

	IF_TAG_PATTERN    = `\[if ([^\]]+)\]`
	ENDIF_TAG_PATTERN = `\[endif\]`

	FOR_TAG_PATTERN    = `\[for ([^\]]+)\]`
	ENDFOR_TAG_PATTERN = `\[endfor\]`

	GO_TEMPLATE = `
package taller

import (
	"fmt"
	"bytes"
	/*IMPORTS*/
)
func render(context Context) *bytes.Buffer {
	output := bytes.NewBufferString("")
	/*CONTENT*/
	return output
}
`
)
//store things such as imports in this register
type Register map[string][]string

// First try to expand/include all forms and fill up the existing data.
// After everything is expanded we proceed to transforming the code into
// golang code
func Parse(content []byte) []byte {
	//use this register to store stuff in the process of parsing
	register := new(Register)
	//execute expand of the template
	content = expandTemplate(content, register, 0)
	//replace all includes
	content = includesTemplate(content, register, 0)
	//finally remove all [block name], [endblock]
	content = stripBlocks(content)

	replace_map := map[string]*bytes.Buffer{
		"IMPORTS": bytes.NewBufferString(""),
		"CONTENT": bytes.NewBufferString("")}
	//add the imports
	for _, import_line := range (*register)["imports"] {
		replace_map["IMPORTS"].Write(
			[]byte("\"" + string(import_line) + "\"\n"))
	}

	lines := SplitToLines(content)
	//each line in the template will be parsed and translated to the
	//appropriate golang code

	for lineno, line := range lines {
		parsed_line := ParseLine(lineno, line)
		replace_map["CONTENT"].Write(parsed_line)
	}
	output := []byte(GO_TEMPLATE)
	//replaces place holders such as /*CONTENT*/ in the output using the parsed_lines
	for placeholder, parsed_lines := range replace_map {
		output = bytes.Replace(output,
			JoinBytes("/*", placeholder, "*/"),
			parsed_lines.Bytes(),
			1)
	}
	return output
}

//find [expand "template.html"] on the first line and expand recursivly
func expandTemplate(content []byte, register *Register, level int) []byte {
	pattern, err := regexp.Compile(EXPAND_PATTERN)
	if err != nil {
		panic("Failed to parse the expand pattern: " + err.String())
	}
	//make sure we register
	registerImports(content, register)

	first_line := SplitToLines(content)[0]
	match := pattern.FindSubmatch(first_line)
	if match == nil { //no expand - return the content
		return content
	}
	// We now need to match the blocks and replace them in the expanded
	// template
	parent_templatefile := string(match[1])
	parent_content := ReadTemplateFile(parent_templatefile)
	first_line_parent := SplitToLines(parent_content)[0]

	match = pattern.FindSubmatch(first_line_parent)
	if match != nil { //we found an expand tag in the parent template
		if level <= 100 {
			panic("Expand loop detected!")
		}
		expandTemplate(content, register, level+1)
	}

	parent_blocks := templateBlocks(parent_content, parent_templatefile)
	current_blocks := templateBlocks(content, "")
	content = parent_content
	for block_name, block_content := range parent_blocks {
		if current_blocks[block_name] != nil {
			content = bytes.Replace(
				content,
				block_content,
				current_blocks[block_name],
				-1)
		}
	}
	return content
}

//given a template content and filename return it's blocks in a map.
//key/value will match block name/block content
func templateBlocks(content []byte, filename string) map[string][]byte {
	var block_names []string

	// Compile patterns
	blocktag_pattern, err := regexp.Compile(BLOCK_TAG_PATTERN)
	if err != nil {
		panic("Failed to parse the block pattern: " + err.String())
	}
	endblocktag_pattern, err := regexp.Compile(ENDBLOCK_TAG_PATTERN)
	if err != nil {
		panic("Failed to parse the endblock pattern: " + err.String())
	}

	blocktag_matches := blocktag_pattern.FindAllSubmatch(content, -1)
	endblocktag_matches := endblocktag_pattern.FindAll(content, -1)
	if len(blocktag_matches) != len(endblocktag_matches) {
		if filename != "" {
			panic("You missed a [block] or an [endblock] in: " + filename)
		} else {
			panic("You missed a [block] or an [endblock] in the current template")
		}
	}

	for _, match := range blocktag_matches {
		block_names = append(block_names, string(match[1]))
	}

	blocks := make(map[string][]byte)
	for _, block_name := range block_names {
		block_start := JoinBytes("[block ", block_name, "]")
		block_end := []byte("[endblock]")
		start_index := bytes.Index(content, block_start)
		end_index := bytes.Index(content, block_end)

		//copy everything thing between [block name] and [endblock]
		blocks[block_name] = content[start_index : end_index+len(block_end)]
		//remove the saved block
		content = content[end_index+len(block_end) : len(content)-1]
	}
	return blocks
}

//find all [include "template.html"] and update the content by adding the
//contents of those files
func includesTemplate(content []byte, register *Register, level int) []byte {
	included := false
	pattern, err := regexp.Compile(INCLUDE_PATTERN)
	if err != nil {
		panic("Failed to parse the include pattern: " + err.String())
	}
	for _, line := range SplitToLines(content) {
		match := pattern.FindSubmatch(line)
		if match == nil { //no expand - return the content
			continue
		}
		included = true
		placeholder := match[0]
		template_content := ReadTemplateFile(string(match[1]))
		content = bytes.Replace(content, placeholder,
			template_content, 1)
	}
	if level >= 100 {
		panic("Detected an include loop!")
	}
	if included { //we included a template we'll try to expand some more
		content = includesTemplate(content, register, level+1)
	}
	return content
}

func registerImports(content []byte, register *Register) {
	import_pattern, err := regexp.Compile(IMPORT_PATTERN)
	if err != nil {
		panic("Failed to parse the import pattern: " + err.String())
	}
	matches := import_pattern.FindAllSubmatch(content, -1)
	r := *register
	if matches != nil {
		for _, match := range matches {
			matched_import := string(match[1])
			r["imports"] = append(r["imports"], matched_import)
		}
	}
}

//Strip all blocks from the parsed template
func stripBlocks(content []byte) []byte {
	blocktag_pattern, _ := regexp.Compile(BLOCK_TAG_PATTERN)
	endblocktag_pattern, _ := regexp.Compile(ENDBLOCK_TAG_PATTERN)
	import_pattern, _ := regexp.Compile(IMPORT_PATTERN)

	content = blocktag_pattern.ReplaceAll(content, []byte(""))
	content = import_pattern.ReplaceAll(content, []byte(""))
	return endblocktag_pattern.ReplaceAll(content, []byte(""))
}

func ParseLine(lineno int, line []byte) []byte {
	//On the first line should always be [extends "filename"]
	parsed_line := []byte("aaa")
	return parsed_line
}
