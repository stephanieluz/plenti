package main

import (
	"io"
	"io/ioutil"
	"os"
)

// Reads all files in "content" folder
// and encodes them as strings literals in generated_defaults.go
func main() {
	fs, _ := ioutil.ReadDir("defaults/content")
	out, _ := os.Create("defaults/generated_defaults.go")
	out.Write([]byte("package defaults"))
	out.Write([]byte("\n\n// Do not edit, this file is automatically generated."))
	out.Write([]byte("\n\n// Defaults: scaffolding used in 'new site' command"))
	out.Write([]byte("\nvar Defaults = map[string][]byte{\n"))
	for _, f := range fs {
		out.Write([]byte("\t\"" + f.Name() + "\": []byte(`"))
		f, _ := os.Open("defaults/content/" + f.Name())
		io.Copy(out, f)
		out.Write([]byte("`),\n"))
	}
	out.Write([]byte("}\n"))
}
