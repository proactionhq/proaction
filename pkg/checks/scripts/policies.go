package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fs, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	out, err := os.Create("policies.go")
	if err != nil {
		panic(err)
	}

	out.Write([]byte("package checks \n\nconst (\n"))

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".rego") {
			name := strings.Title(strings.TrimSuffix(f.Name(), ".rego"))
			name = strings.ReplaceAll(name, "-", "")
			out.Write([]byte(name + "_Policy = `"))
			f, err := os.Open(f.Name())
			if err != nil {
				panic(err)
			}

			io.Copy(out, f)
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
