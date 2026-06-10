package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	bucket := os.Getenv("GCS_BUCKET")
	if bucket == "" {
		log.Fatal("GCS_BUCKET environment variable is required")
	}
	// Strip this prefix from request paths before mapping to a GCS object name.
	// Example: STRIP_PREFIX=/my-bucket means GET /my-bucket/file.html → object "file.html"
	// Leave empty to serve the full request path as the object name.
	stripPrefix := os.Getenv("STRIP_PREFIX")
	if stripPrefix == "" {
		stripPrefix = "/" + bucket
	}
	// ROOT_REDIRECT: where to redirect bare "/" requests. Defaults to /index.html.
	rootRedirect := os.Getenv("ROOT_REDIRECT")
	if rootRedirect == "" {
		rootRedirect = "/index.html"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer client.Close()

	bkt := client.Bucket(bucket)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, rootRedirect, http.StatusFound)
			return
		}

		objectName := strings.TrimPrefix(r.URL.Path, stripPrefix)
		objectName = strings.TrimPrefix(objectName, "/")
		if objectName == "" {
			objectName = "index.html"
		}
		if strings.HasSuffix(objectName, "/") {
			objectName += "index.html"
		}

		obj := bkt.Object(objectName)
		attrs, err := obj.Attrs(ctx)
		if err == storage.ErrObjectNotExist {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if err != nil {
			log.Printf("Error getting attrs for %s: %v", objectName, err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		reader, err := obj.NewReader(ctx)
		if err != nil {
			log.Printf("Error reading %s: %v", objectName, err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		defer reader.Close()

		contentType := attrs.ContentType
		if contentType == "" || contentType == "application/octet-stream" || contentType == "text/plain" {
			if ct := mime.TypeByExtension(filepath.Ext(objectName)); ct != "" {
				contentType = ct
			}
		}
		w.Header().Set("Content-Type", contentType)
		if attrs.CacheControl != "" {
			w.Header().Set("Cache-Control", attrs.CacheControl)
		}
		if attrs.ContentEncoding != "" {
			w.Header().Set("Content-Encoding", attrs.ContentEncoding)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", attrs.Size))

		io.Copy(w, reader)
	})

	h2s := &http2.Server{}
	handler := h2c.NewHandler(http.DefaultServeMux, h2s)

	log.Printf("Starting GCS proxy: bucket=%s strip_prefix=%q root_redirect=%s port=%s", bucket, stripPrefix, rootRedirect, port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
