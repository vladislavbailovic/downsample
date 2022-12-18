package main

import (
	"compress/gzip"
	"compress/zlib"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	//go:embed public_html/assets/downsample.wasm
	downsampleWasm []byte
	//go:embed public_html/assets/wasm_exec.js
	wasmLoader []byte
	//go:embed public_html/index.html
	indexHtml []byte
	//go:embed public_html/sample.jpg
	sampleImage []byte
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

		// path, err := os.Getwd()
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		if strings.Contains(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
			w.Write(downsampleWasm)
		} else if strings.Contains(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
			w.Write(wasmLoader)
		} else if strings.Contains(r.URL.Path, ".jpg") {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(sampleImage)
		} else if strings.Contains(r.URL.Path, "favicon") {
			fmt.Println("ignoring favicon")
		} else {
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexHtml)
		}
	})))
}
