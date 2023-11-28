package ascii_art

import (
	ascii "ascii_art/pkg"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	Banner  string
	Msg     string
	Format  string
	Output  string
	Display bool
}

// Return Error status code if an error occurred
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {

	w.WriteHeader(status) // Write status code

	switch status {
	case http.StatusNotFound:
		fmt.Fprint(w, http.StatusText(status))
	case http.StatusBadRequest:
		fmt.Fprint(w, http.StatusText(status))
	case http.StatusInternalServerError:
		fmt.Fprint(w, http.StatusText(status))
	case http.StatusMethodNotAllowed:
		fmt.Fprint(w, http.StatusText(status))
	}
}

// Handle home and display the main page of the application
func (d *Data) home(w http.ResponseWriter, r *http.Request) {
	var err error

	if len(r.URL.RawQuery) != 0 {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	tmp := template.Must(template.ParseFiles("templates/index.html")) // Create template of main page

	switch r.Method { // Compare request method
	case "GET":
		d.Display = false
		os.Remove("tmp")
		d.Output = ""
		err := tmp.Execute(w, nil) // send template to client
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	case "POST":
		r.ParseForm()                                               // parse form send by client
		d.Msg = strings.ReplaceAll(r.Form.Get("msg"), "\r\n", "\n") // store message from text area
		d.Banner = r.Form.Get("banner")                             // store banner from banner selector
		d.Format = r.Form.Get("format")
		d.Display = true
		d.Output, err = ascii.Process(d.Msg, d.Banner) // Run ascii art program
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, d) // Send new template with data in
	default:
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
}

// Set header and send file to download to the client
func (d *Data) download(w http.ResponseWriter, r *http.Request) {
	var ct string

	if len(d.Output) == 0 {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	ascii.WriteOutputFile("tmp", d.Output)

	if _, err := os.Stat("tmp"); err != nil {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	switch d.Format {
	case "xml":
		ct = "application/xml"
	case "html":
		ct = "text/html"
	case "txt":
		ct = "text/plain"
	case "md":
		ct = "text/markdown"
	default:
		ct = "application/octet-stream"
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art-web-export-file."+d.Format)
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	if !strings.Contains(w.Header().Get("Content-Disposition"), "attachment;") {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "tmp")
}

// Main handler how define each path
func Handlers() {

	d := &Data{Format: "txt"}

	http.HandleFunc("/", d.home)                                                               // Main path
	http.HandleFunc("/download", d.download)                                                   // Download
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // linked css

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
