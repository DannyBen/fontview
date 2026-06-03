package report

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Run(opts Options) error {
	fonts, err := loadFonts(opts.Inputs)
	if err != nil {
		return err
	}

	if opts.WriteHTML {
		page, err := render(fonts, opts.Version)
		if err != nil {
			return err
		}
		if err := os.WriteFile(opts.Output, page, 0644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "fontview: wrote %s with %d font(s)\n", opts.Output, len(fonts))
		return nil
	}

	page, err := render(serverFonts(fonts), opts.Version)
	if err != nil {
		return err
	}
	return serve(opts.Addr, page)
}

func serve(addr string, page []byte) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(page)
	})
	mux.HandleFunc("/font/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/font/")
		path = filepath.Clean(filepath.FromSlash(path))
		if strings.HasPrefix(path, "..") || filepath.IsAbs(path) {
			http.Error(w, "invalid font path", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, path)
	})

	fmt.Fprintf(os.Stderr, "fontview: serving http://%s\n", listener.Addr())
	log.Fatal(http.Serve(listener, mux))
	return nil
}
