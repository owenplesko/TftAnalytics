package fileserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Fileserver struct {
	Directory     string
	IndexFilepath string
}

func (fs Fileserver) ListenAndServe(port int) {
	_, err := os.Stat(fs.Directory)
	if err != nil {
		log.Fatalf("Directory '%s' not found.\n", fs.Directory)
	}

	router := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(fs.Directory))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if isFilePath(r.URL.Path) {
			fileServer.ServeHTTP(w, r)
		} else {
			fs.serveIndex(w, r)
		}
	})

	log.Printf("fileserver listening on port: %d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatalf("Error starting server: %s\n", err)
}

func isFilePath(path string) bool {
	return strings.Contains(path, ".")
}

func (fs Fileserver) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fs.IndexFilepath)
}
