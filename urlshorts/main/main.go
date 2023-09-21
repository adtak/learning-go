package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/adtak/urlshort"
)

func main() {
	var (
		file_name = flag.String("yaml", "path-url.yaml", "a yaml file to map url from path")
	)
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, defaultMux())

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yml := readYaml(*file_name)
	yamlHandler, err := urlshort.YAMLHandler(yml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readYaml(file_name string) []byte {
	file, err := os.ReadFile(file_name)
	if err != nil {
		panic(err)
	}
	return file
}
