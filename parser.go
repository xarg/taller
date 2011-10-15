package taller

import (
	"regexp"
	"bytes"
)

const (
	ECHO_DELIMITER          = '`'

	EXPAND_PATTERN          = `\[expand ([^\]]+)\]`
	INCLUDE_PATTERN         = `\[include ([^\]]+)\]`
	IMPORT_PATTERN          = `\[import ([^\]]+)\]`

	BLOCK_TAG_PATTERN       = `\[block ([^\]]+)\]`
	ENDBLOCK_TAG_PATTERN    = `\[endblock\]`

	IF_TAG_PATTERN          = `\[if ([^\]]+)\]`
	ENDIF_TAG_PATTERN       = `\[endif\]`

	FOR_TAG_PATTERN         = `\[for ([^\]]+)\]`
	ENDFOR_TAG_PATTERN      = `\[endfor\]`

	GO_TEMPLATE             = `
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

// First try to expand/include all forms and fill up the existing data.
// After everything is exapanded we proceed to transforming the code into
// golang code
func Parse(content []byte) []byte {
	//execute expand of the template
	content = expandTemplate(content)
	//concatenate all includes
	content = includesTemplate(content)
	//get all imports
	imports := getImports(content)
	//finally remove all [block name], [endblock] and [import]
	content = stripBlocks(content)

	replace_map := map[string]*bytes.Buffer{
		"IMPORTS": bytes.NewBufferString(""),
		"CONTENT": bytes.NewBufferString("")}
	//add the imports
	for _, import_line := range imports {
		replace_map["IMPORTS"].Write([]byte("\"" + string(import_line) +"\"\n"))
	}

	lines := bytes.Split(content, []byte("\n"))
	//each line in the template will be parsed and translated to the appropriate golang code

	for lineno, line := range lines {
		parsed_line := ParseLine(lineno, line)
		replace_map["CONTENT"].Write(parsed_line)
	}
	output := []byte(GO_TEMPLATE)
	//replaces placeholders such as /*CONTENT*/ in the output using the parsed_lines
	for placeholder, parsed_lines := range replace_map {
		output = bytes.Replace(output,
				JoinBytes("/*", placeholder, "*/"),
				parsed_lines.Bytes(),
				1)
	}
	return output
}

//find [expand "template.html"] on the first line and expand
func expandTemplate(content []byte) []byte {
	first_line := SplitToLines(content)[0]
	pattern, err := regexp.Compile(EXPAND_PATTERN)
	if err != nil {
		panic("Failed to parse the expand pattern: " + err.String())
	}
	match := pattern.FindSubmatch(first_line)
	if match == nil {//no expand - return the content
		panic(string(first_line))
		return content
	}
	// We now need to match the blocks and replace them in the expanded
	// template
	parent_templatefile := string(match[1])

	parent_content := ReadTemplateFile(parent_templatefile)
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
//key/val will match block name/block content
func templateBlocks(content []byte, filename string) map[string] []byte {
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

	blocks := make(map[string] []byte)
	for _, block_name := range block_names {
		block_start := JoinBytes("[block ", block_name, "]")
		block_end := []byte("[endblock]")
		start_index := bytes.Index(content, block_start)
		end_index := bytes.Index(content, block_end)

		//copy everything thing between [block name] and [endblock]
		blocks[block_name] = content[start_index:end_index + len(block_end)]
		//remove the saved block
		content = content[end_index + len(block_end):len(content)-1]
	}
	return blocks
}

//find all [include "template.html"] and update the content by adding the
//contents of those files
func includesTemplate(content []byte) []byte {
	return includeTemplate(content, 0)
}

// resolve includes and make sure we avoid infinite include loops.
func includeTemplate(content []byte, level int) []byte {
	included := false
	pattern, err := regexp.Compile(INCLUDE_PATTERN)
	if err != nil {
		panic("Failed to parse the include pattern: " + err.String())
	}
	for _, line := range SplitToLines(content) {
		match := pattern.FindSubmatch(line)
		if match == nil {//no expand - return the content
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
	if included {//we included a template we'll try to expand some more
		content = includeTemplate(content, level+1)
	}
	return content
}

func getImports(content []byte) [][]byte {
	var imports [][]byte
	import_pattern, err := regexp.Compile(IMPORT_PATTERN)
	if err != nil {
		panic("Failed to parse the import pattern: " + err.String())
	}
	matches := import_pattern.FindAllSubmatch(content, -1)
	if matches != nil {
		for _, match := range matches {
			imports = append(imports, match[1])
		}
	}
	return imports
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
