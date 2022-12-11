package main

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type compressedWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *compressedWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}
func (w *compressedWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func serveCompressed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := compressedWriter{ResponseWriter: w}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(w)
			defer gz.Close()

			writer.Writer = gz
			w.Header().Set("Content-Encoding", "gzip")
		} else if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
			zl := zlib.NewWriter(w)
			defer zl.Close()

			writer.Writer = zl
			w.Header().Set("Content-Encoding", "deflate")
		} else {
			writer.Writer = w
		}
		next.ServeHTTP(&writer, r)
	})
}

func main() {
	http.ListenAndServe(":6660", serveCompressed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.Contains(r.URL.Path, ".wasm") {
			fd, err := os.Open(filepath.Join(path, "public_html/assets/downsample.wasm"))
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/wasm")
			io.Copy(w, fd)
		} else if strings.Contains(r.URL.Path, ".js") {
			fd, err := os.Open(filepath.Join(path, "public_html/assets/wasm_exec.js"))
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "text/javascript")
			io.Copy(w, fd)
		} else if strings.Contains(r.URL.Path, ".jpg") {
			fd, err := os.Open(filepath.Join(path, "public_html/sample.jpg"))
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "image/jpeg")
			io.Copy(w, fd)
		} else if strings.Contains(r.URL.Path, "favicon") {
			fmt.Println("ignoring favicon")
		} else {
			index, err := os.Open(filepath.Join(path, "public_html/index.html"))
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			io.Copy(w, index)
		}
	})))
}
