package fileserver

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed static/*
var static embed.FS

func ListenAndServe(port int) {
	files, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	index, err := static.ReadFile("static/index.html")
	if err != nil {
		log.Panicf("error reading static/index.html: %v", err)
	}

	router := http.NewServeMux()
	fileServer := http.FileServer(http.FS(files))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if isFilePath(r.URL.Path) {
			fileServer.ServeHTTP(w, r)
		} else {
			w.Write(index)
		}
	})

	log.Printf("fileserver listening on port: %d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatalf("Error starting server: %s\n", err)
}

func isFilePath(path string) bool {
	return strings.Contains(path, ".")
}
