package taller

import (
	"bytes"
	"path"
	"io/ioutil"
	"strings"
	"os"
)

const (
	TALLER_ENV_VARIABLE = "TALLER_PATH"
)
//read taller environment variable and split the paths
func GetTallerPaths() []string {
	return strings.Split(os.Getenv(TALLER_ENV_VARIABLE), ":")
}

//return the content of the first found template
func ReadTemplateFile(filename string) []byte {
	ok := false
	for _, path_dir := range GetTallerPaths() {
		absolute_path := path.Join(path_dir, filename)
		if _, err := os.Stat(absolute_path); err != os.ENOENT {
			filename = absolute_path
			ok = true
			break
		}
	}
	if !ok {
		panic("Template not found in TALLER_PATH folders: " + filename)
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Failed to read template file: " + err.String())
	}
	return content
}

//given a list of strings join them in a []byte
func JoinBytes(string_list ...string) []byte {
	var joined [][]byte
	for _, str := range string_list {
		joined = append(joined, []byte(str))
	}
	return bytes.Join(joined, []byte(""))
}

func SplitToLines(content []byte) [][]byte {
	return bytes.Split(content, []byte("\n"))
}
