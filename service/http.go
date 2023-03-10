package service

import (
	"fmt"
	"io"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func AppIndex(appservice *AppService, w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("app index"))
}

// handles both websocket and http GET/POST requests on the same path
func handleConnectionForHttpOrWebsocket(appservice *AppService, w http.ResponseWriter, r *http.Request) {
	log.Warnf("serve handleConnectionForHttpOrWebsocket %s", r.URL.Path)
	if r.URL.Path != "/test.company.com/testservice/mytestservice" {
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Upgrade") == "websocket" {
		HandleWebsocket(appservice, w, r)
	} else {
		AppIndex(appservice, w, r)
	}
}

func ServeFile(fs http.FileSystem, w http.ResponseWriter, r *http.Request, filename string, contenttype string) {
	log.Warnf("ServeFile %s", filename)

	file, err := fs.Open(filename)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	content, err_file := io.ReadAll(file)
	if err_file != nil {
		fmt.Fprintf(w, "Error while reading file: %v", err_file)
		http.Error(w, err_file.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", contenttype)
	_, err_write := w.Write(content)
	if err_write != nil {
		fmt.Fprintf(w, "Error while writing to response: %v", err_write)
		http.Error(w, err_write.Error(), http.StatusInternalServerError)
		return
	}

}

// func GetHttpRoutes(appservice *AppService, mux *chi.Mux) error {

// 	// Serve static files on app path
// 	// mux.Get("/static/*",
// 	// 	http.StripPrefix("/test.company.com/testservice/mytestservice/static/",
// 	// 		http.FileServer(appservice.Fs),
// 	// 	).ServeHTTP,
// 	// )

// 	// mux.Get("/admin.htm", func(w http.ResponseWriter, r *http.Request) {
// 	// 	ServeFile(appservice.Fs, w, r, "admin.htm", "text/html")
// 	// })
// 	// mux.Get("/user.htm", func(w http.ResponseWriter, r *http.Request) {
// 	// 	log.Warnf("serve /user.htm %s", r.URL.Path)
// 	// 	ServeFile(appservice.Fs, w, r, "user.htm", "text/html")
// 	// })
// 	// mux.Get("/searchapi.htm", func(w http.ResponseWriter, r *http.Request) {
// 	// 	ServeFile(appservice.Fs, w, r, "searchapi.htm", "text/html")
// 	// })

// 	// mux.Get("/user.png", func(w http.ResponseWriter, r *http.Request) {
// 	// 	ServeFile(appservice.Fs, w, r, "user.png", "image/png")
// 	// })
// 	// mux.Get("/admin.png", func(w http.ResponseWriter, r *http.Request) {
// 	// 	ServeFile(appservice.Fs, w, r, "admin.png", "image/png")
// 	// })
// 	// mux.Get("/app.css", func(w http.ResponseWriter, r *http.Request) {
// 	// 	ServeFile(appservice.Fs, w, r, "app.css", "text/css; charset=utf-8")
// 	// })
// 	// mux.Get("/app.js", func(w http.ResponseWriter, r *http.Request) {
// 	// 	ServeFile(appservice.Fs, w, r, "app.js", "text/javascript; charset=utf-8")
// 	// })

// 	return nil
// }
